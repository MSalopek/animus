package queue

import (
	"log"
	"sync"

	"github.com/nsqio/go-nsq"
)

func BodyLogger(m *nsq.Message) error {
	log.Println("message:", string(m.Body))
	return nil
}

func ChanBodyLogger(C <-chan []byte) {
	log.Println("ChanBodyLogger STARTED")
	for m := range C {
		log.Println("message:", string(m))
	}
}

func WgChanBodyLogger(wg *sync.WaitGroup, C <-chan []byte) {
	defer wg.Done()
	log.Println("WgChanBodyLogger STARTED")
	for m := range C {
		log.Println("message:", string(m))
	}
	log.Println("WgChanBodyLogger STOPPED")
}
