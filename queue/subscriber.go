package queue

import (
	log "github.com/sirupsen/logrus"

	"github.com/nsqio/go-nsq"
)

type Subscriber interface {
	Subscribe(chan<- []byte) nsq.HandlerFunc
	Unsubscribe()
}

type NsqConsumer struct {
	*nsq.Consumer
}

func NewSubscriber(topic, channel, lookupd string, msgs chan<- []byte) (*NsqConsumer, error) {
	conf := nsq.NewConfig()
	n := &NsqConsumer{}
	c, err := nsq.NewConsumer(topic, channel, conf)
	if err != nil {
		return nil, err
	}
	c.AddHandler(nsq.HandlerFunc(n.Subscribe(msgs)))

	err = c.ConnectToNSQLookupd(lookupd)
	if err != nil {
		return nil, err
	}
	n.Consumer = c
	return n, nil
}

func MustNewSubscriber(topic, channel, lookupd string, msgs chan<- []byte) *NsqConsumer {
	c, err := NewSubscriber(topic, channel, lookupd, msgs)
	if err != nil {
		log.Fatal(err)
	}
	return c
}

// Consume pushes messages into C channel.
func (n *NsqConsumer) Subscribe(C chan<- []byte) nsq.HandlerFunc {
	return func(m *nsq.Message) error {
		C <- m.Body
		return nil
	}
}

func (n *NsqConsumer) Unsubscribe() {
	n.Consumer.Stop()
	// wait for cleanup to complete
	<-n.StopChan
}
