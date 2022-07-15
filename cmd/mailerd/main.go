package main

import (
	"bytes"
	"context"
	"errors"
	"flag"
	"fmt"
	"html/template"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/mailjet/mailjet-apiv3-go"
	"github.com/msalopek/animus/queue"

	log "github.com/sirupsen/logrus"
)

const registrationTemplate = `<!DOCTYPE html>
<html lang="en">

<head>
    <meta charset="UTF-8">
</head>

<body>
	<h3>Dear {{ .Username }}, welcome to Animus!</h3>
	<br />
	<p>Please activate your account by clicking this link: <a href={{ .URL }}>{{ .URL }}</a></p>!
</body>

</html>
`

var dryRun = flag.Bool("dry-run", false, "output emails to stdout (does not send API requests)")
var key = flag.String("api-key", "", "mailer REST API key")
var secret = flag.String("api-secret", "", "mailer REST API secret")
var topic = flag.String("topic", "", "nsq email topic")
var nsqLookup = flag.String("nsq-lookup", "", "nsq lookup for nsq connections")
var debug = flag.Bool("debug", false, "show debug logs")

func main() {
	flag.Parse()

	if *key == "" {
		log.Fatal("api-key is required")
	}
	if *secret == "" {
		log.Fatal("api-secret is required")
	}
	if *topic == "" {
		log.Fatal("topic is required")
	}
	if *nsqLookup == "" {
		log.Fatal("nsq-lookup is required")
	}

	var wg sync.WaitGroup
	ctx := exitSignal()

	logger := log.New()
	logger.Out = os.Stdout
	logger.SetFormatter(&log.JSONFormatter{TimestampFormat: "2006-01-02 15:04:05"})

	if *debug {
		logger.SetFormatter(&log.TextFormatter{TimestampFormat: "2006-01-02 15:04:05", FullTimestamp: true})
		logger.SetLevel(log.DebugLevel)
	}

	m := NewMailer(*key, *secret, *topic, *nsqLookup, "matija@animus.store", "Matija from Animus", logger, *dryRun)
	go m.HandleRegistrationEmails(&wg)

	<-ctx.Done()
	logger.Info("mailer stopping")
	m.Stop()
	wg.Wait()
	logger.Info("mailer stopped")
}

// exitSignal cancels the context on SIGINT, SIGTERM
// use the returned context as the root context
func exitSignal() context.Context {
	ctx, stop := context.WithCancel(context.Background())
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		<-c
		stop()
	}()
	return ctx
}

type EmailParams struct {
	Subject      string
	From         string
	FromName     string
	To           string
	ToName       string
	HTMLTemplate string
	CustomID     string
}

func (ep *EmailParams) Valid() error {
	if ep.Subject == "" {
		return fmt.Errorf("must specify EmailParams.Subject")
	}
	if ep.To == "" {
		return fmt.Errorf("must specify EmailParams.To")
	}
	if ep.ToName == "" {
		return fmt.Errorf("must specify EmailParams.ToName")
	}
	if ep.HTMLTemplate == "" {
		return fmt.Errorf("must specify EmailParams.HTMLTemplate")
	}

	return nil
}

type Mailer struct {
	Client             *mailjet.Client
	DefaultSenderEmail string
	DefaultSenderName  string

	// nsq subscriber
	subscriber queue.Subscriber
	messages   chan []byte

	logger *log.Logger

	// log email templates, don't send requests
	dryRun bool
}

func NewMailer(key, secret, topic, lookup, sender, senderName string, log *log.Logger, dryRun bool) *Mailer {
	if key == "" {
		log.Fatal("missing key")
	}
	if secret == "" {
		log.Fatal("missing secret")
	}
	if sender == "" {
		log.Fatal("missing sender")
	}
	if senderName == "" {
		log.Fatal("missing senderName")
	}

	m := &Mailer{
		Client:             mailjet.NewMailjetClient(key, secret),
		DefaultSenderEmail: sender,
		DefaultSenderName:  senderName,
		messages:           make(chan []byte),
		logger:             log,
		dryRun:             dryRun,
	}
	m.subscriber = queue.MustNewSubscriber(
		topic,
		"mailer",
		lookup,
		m.messages,
	)
	return m
}

func (m *Mailer) Stop() {
	m.subscriber.Unsubscribe()
	close(m.messages)
}

func (m *Mailer) HandleRegistrationEmails(wg *sync.WaitGroup) {
	defer wg.Done()
	for msg := range m.messages {
		if msg == nil {
			m.logger.WithFields(log.Fields{"func": "HandleRegistrationEmails"}).Warn("got nil message")
			break
		}

		req := &queue.RegisterEmail{}
		if err := req.Unmarshal(msg); err != nil {
			m.logger.WithFields(log.Fields{"func": "HandleRegistrationEmails"}).Error(err)
			continue
		}

		tmp, err := template.New("mail").Parse(registrationTemplate)
		if err != nil {
			m.logger.WithFields(log.Fields{"func": "HandleRegistrationEmails"}).Error(err)
			continue
		}

		var buf bytes.Buffer
		if err := tmp.Execute(&buf, req); err != nil {
			m.logger.WithFields(log.Fields{"func": "HandleRegistrationEmails"}).Error(err)
		}

		params := &EmailParams{
			From:         m.DefaultSenderEmail,
			FromName:     m.DefaultSenderName,
			To:           req.Email,
			ToName:       req.Email,
			HTMLTemplate: buf.String(),
			Subject:      "Activate your Animus account",
		}

		wg.Add(1)
		go m.Send(wg, params)
	}
	m.logger.WithFields(log.Fields{"func": "HandleRegistrationEmails"}).Debug("terminating")
}

func (m *Mailer) Send(wg *sync.WaitGroup, p *EmailParams) error {
	defer wg.Done()

	if p == nil {
		err := errors.New("email params not specified")
		m.logger.WithFields(log.Fields{"func": "Send", "to": p.To, "subject": p.Subject}).Error(err)
		return err
	}

	sender := p.From
	if sender == "" {
		sender = m.DefaultSenderEmail
	}

	senderName := p.FromName
	if senderName == "" {
		senderName = m.DefaultSenderName
	}

	if err := p.Valid(); err != nil {
		m.logger.WithFields(log.Fields{"func": "Send", "to": p.To, "subject": p.Subject}).Error(err)
		return err
	}

	if !m.dryRun {
		msgs := []mailjet.InfoMessagesV31{
			{
				From: &mailjet.RecipientV31{
					Email: sender,
					Name:  senderName,
				},
				To: &mailjet.RecipientsV31{
					mailjet.RecipientV31{
						Email: p.To,
						Name:  p.ToName,
					},
				},
				Subject: p.Subject,
				// TextPart: "My first Mailjet email",
				HTMLPart: p.HTMLTemplate,
				CustomID: p.CustomID,
			},
		}

		send := mailjet.MessagesV31{Info: msgs}
		if _, err := m.Client.SendMailV31(&send); err != nil {
			m.logger.WithFields(log.Fields{"func": "Send", "to": p.To, "subject": p.Subject}).Error(err)
			return err
		}

		return nil
	}

	m.logger.WithFields(log.Fields{
		"from":     p.From,
		"fromName": p.FromName,
		"to":       p.To,
		"toName":   p.ToName,
		"content":  p.HTMLTemplate,
		"subject":  p.Subject}).Info("dry run")
	return nil
}
