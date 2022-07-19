package mailer

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"sync"

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
	<p>Hello
	{{ if .Firstname }}
		{{ .Firstname }}
	{{ else }}
		{{ .Username }}
	{{ end }} , welcome to Animus!</p>
	<p>Please activate your account by clicking this link: <a href={{ .URL }}>{{ .URL }}</a></p>
	<p>We are happy to have you on board.</p>
</body>

</html>
`

type EmailParams struct {
	Type         string
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

func New(key, secret, topic, lookup, sender, senderName string, log *log.Logger, dryRun bool) *Mailer {
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

func (m *Mailer) HandleMessages(wg *sync.WaitGroup) {
	defer wg.Done()

	for raw := range m.messages {
		if raw == nil {
			m.logger.WithFields(log.Fields{"func": "HandleMessages"}).Warn("got nil message")
			break
		}

		msg := &queue.MailerMessage{}
		if err := msg.Unmarshal(raw); err != nil {
			m.logger.WithFields(log.Fields{"func": "HandleMessages"}).Error(err)
			continue
		}

		if err := msg.Validate(); err != nil {
			m.logger.WithFields(log.Fields{"func": "HandleMessages"}).Error(err)
			continue
		}

		var hydratedEmail bytes.Buffer
		params := EmailParams{
			Type:     string(msg.Type),
			From:     m.DefaultSenderEmail,
			FromName: m.DefaultSenderName,
			To:       msg.Email,
			ToName:   msg.Email,
		}

		if msg.Firstname != nil {
			params.ToName = *msg.Firstname

		}

		if msg.Lastname != nil {
			params.ToName = fmt.Sprintf("%s %s", params.ToName, *msg.Lastname)
		}

		switch msg.Type {
		case queue.MailerTypeResetPass:
			m.logger.WithFields(log.Fields{"func": "HandleMessages"}).Info("password reset flow not implemented")
			continue

		case queue.MailerTypeRegister:
			tmp, err := template.New("mail").Parse(registrationTemplate)
			if err != nil {
				m.logger.WithFields(log.Fields{"func": "HandleMessages"}).Error(err)
				continue
			}
			if err := tmp.Execute(&hydratedEmail, msg); err != nil {
				m.logger.WithFields(log.Fields{"func": "HandleMessages"}).Error(err)
				continue
			}
			params.Subject = "Activate your Animus account"
		default:
			m.logger.WithFields(log.Fields{"func": "HandleMessages"}).Error(fmt.Errorf("unknown message type: %s", msg.Type))
			continue
		}

		params.HTMLTemplate = hydratedEmail.String()

		wg.Add(1)
		go m.Send(wg, &params)
	}
	m.logger.WithFields(log.Fields{"func": "HandleMessages"}).Debug("terminating")
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

	if m.dryRun {
		m.logger.WithFields(log.Fields{
			"from":     p.From,
			"fromName": p.FromName,
			"to":       p.To,
			"toName":   p.ToName,
			"content":  p.HTMLTemplate,
			"subject":  p.Subject,
			"type":     p.Type}).Info("dry run")
		return nil
	}

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
		m.logger.WithFields(log.Fields{"func": "Send", "to": p.To, "subject": p.Subject, "type": p.Type}).Error(err)
		return err
	}
	m.logger.WithFields(log.Fields{"func": "Send", "to": p.To, "subject": p.Subject, "type": p.Type}).Info("sent email request")
	return nil

}
