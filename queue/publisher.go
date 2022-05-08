package queue

import (
	log "github.com/sirupsen/logrus"

	"github.com/nsqio/go-nsq"
)

type Publisher interface {
	Publish(topic string, body []byte) error
	Stop()
}

type NsqPublisher struct {
	Topic string
	*nsq.Producer
}

func NewPublisher(topic, nsqd string) (*NsqPublisher, error) {
	conf := nsq.NewConfig()
	n := &NsqPublisher{Topic: topic}
	p, err := nsq.NewProducer(nsqd, conf)
	if err != nil {
		return nil, err
	}
	n.Producer = p
	return n, nil
}

func MustNewPublisher(topic, nsqd string) *NsqPublisher {
	p, err := NewPublisher(topic, nsqd)
	if err != nil {
		log.Fatal(err)
	}
	return p
}
