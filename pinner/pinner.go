package pinner

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	shell "github.com/ipfs/go-ipfs-api"
	"github.com/minio/minio-go/v7"
	"github.com/nsqio/go-nsq"

	"github.com/msalopek/animus/engine/repo"
	"github.com/msalopek/animus/model"
	"github.com/msalopek/animus/queue"
	"github.com/msalopek/animus/storage"

	log "github.com/sirupsen/logrus"
)

const defaultMaxConcurrentRequests = 10

type Pinner struct {
	sm *storage.Manager // s3 compatible storage

	sh      *shell.Shell // ipfs api shell
	autoPin bool         // pin files automatically

	repo *repo.Repo

	Messages              chan []byte
	subscriber            queue.Subscriber
	maxConcurrentRequests int

	producer *nsq.Producer

	logger *log.Logger
	conf   *Config
}

func New(conf *Config, repo *repo.Repo, logger *log.Logger) *Pinner {
	p := &Pinner{
		conf:                  conf,
		logger:                logger,
		sm:                    storage.MustNewManager(conf.Storage),
		sh:                    shell.NewShell(conf.NodeApiURL),
		maxConcurrentRequests: conf.MaxConcurrentRequests,

		Messages: make(chan []byte),
	}
	if conf.MaxConcurrentRequests < 1 {
		p.maxConcurrentRequests = defaultMaxConcurrentRequests
	}
	p.subscriber = queue.MustNewSubscriber(
		conf.SubscribeTopic,
		"pinner",
		conf.NsqLookupdURL,
		p.Messages)

	if conf.LocalShell {
		p.sh = shell.NewLocalShell()
	}

	p.logger.WithFields(log.Fields{
		"topic":          conf.SubscribeTopic,
		"max-concurrent": conf.MaxConcurrentRequests,
	}).Info("started pinner")
	return p
}

func (p *Pinner) Stop() {
	p.subscriber.Unsubscribe()
	close(p.Messages)
}

func (p *Pinner) SetAutopin(val bool) {
	p.autoPin = val
}

// read messages from Messages chan and process pin requests.
// All in-flight requests will be processed before terminating.
func (p *Pinner) handlePinRequests(wg *sync.WaitGroup) {
	for m := range p.Messages {
		if m == nil {
			log.WithFields(log.Fields{"func": "handlePinRequests"}).Error("got nil message")
			break
		}

		req := &queue.PinRequest{}
		if err := req.Unmarshal(m); err != nil {
			log.WithFields(log.Fields{"func": "handlePinRequests"}).Error(err)
			continue
		}
		wg.Add(1)
		go p.handleAdd(wg, req)
	}
	log.WithFields(log.Fields{"func": "handlePinRequests"}).Error("terminating")
}

func (p *Pinner) handleAdd(wg *sync.WaitGroup, req *queue.PinRequest) {
	defer wg.Done()
	hash, err := p.Add(req)
	if err != nil {
		log.WithFields(
			log.Fields{
				"request": fmt.Sprintf("%#v", req),
			}).Error(err)
		return
	}
	stage := model.UploadStageIPFS
	s := &model.Storage{
		StorageKey:  &hash,
		UploadStage: &stage,
		UpdatedAt:   time.Now(),
	}
	// check errors
	p.repo.Save(s)
	log.WithFields(
		log.Fields{
			"hash":    hash,
			"request": fmt.Sprintf("%#v", req),
		}).Info("uploaded pin request")
}

// Add uploads and pins file or directory under pr.Key to IPFS using StorageMessage.
// If key corresponds to a directory its contents will be downloaded to
// a tmp directory and uploaded to IPFS. The directory is removed after upload.

// If key corresponds to a file the file will be streamed to IPFS.
func (p *Pinner) Add(pr *queue.PinRequest) (string, error) {
	ctx := context.Background()
	var hash, tmp string
	var err error
	if pr.Dir {
		tmp, err = p.sm.DownloadDir(ctx, p.conf.Bucket, pr.Key)
		if err != nil {
			return hash, err
		}
		log.WithFields(log.Fields{
			"temp":    tmp,
			"request": fmt.Sprintf("%#v", pr),
		}).Debug("created temp dir")

		hash, err := p.sh.AddDir(tmp)
		if err != nil {
			return hash, err
		}

		log.WithFields(log.Fields{
			"temp": tmp,
		}).Debug("removing temp dir")
		err = os.RemoveAll(tmp)
		if err != nil {
			log.WithFields(log.Fields{
				"temp":    tmp,
				"message": "error removing temp dir",
			}).Error(err)
		}
		return hash, nil
	}
	obj, err := p.sm.GetObject(ctx, p.conf.Bucket, pr.Key, minio.GetObjectOptions{})
	if err != nil {
		log.WithFields(log.Fields{
			"request": fmt.Sprintf("%#v", pr),
			"message": "error streaming object",
		}).Error(err)
		return "", err
	}
	return p.sh.Add(obj)
}
