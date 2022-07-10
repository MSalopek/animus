package pinner

import (
	"context"
	"errors"
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

// internal request throttling
// there will be at maximum defaultMaxConcurrentRequests
// goroutines at any time during application operation
const defaultMaxConcurrentRequests = 5

type Pinner struct {
	sm *storage.Manager // s3 compatible storage

	sh      *shell.Shell // ipfs api shell
	autoPin bool         // pin files automatically

	repo *repo.Repo

	Messages   chan []byte
	subscriber queue.Subscriber
	rateLimit  chan struct{}

	producer *nsq.Producer

	logger *log.Logger
	conf   *Config
}

func New(conf *Config, repo *repo.Repo, logger *log.Logger) *Pinner {
	p := &Pinner{
		autoPin:  true,
		Messages: make(chan []byte),

		repo: repo,
		sm:   storage.MustNewManager(conf.Storage),
		sh:   shell.NewShell(conf.NodeApiURL),

		rateLimit: make(chan struct{}, defaultMaxConcurrentRequests),
		logger:    logger,
		conf:      conf,
	}
	if conf.MaxConcurrentRequests < 1 {
		p.rateLimit = make(chan struct{}, defaultMaxConcurrentRequests)
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
	close(p.rateLimit)
}

func (p *Pinner) SetAutopin(val bool) {
	p.autoPin = val
}

// Read and process PinRequests from Messages chan.
// All in-flight requests will be processed before terminating.
func (p *Pinner) HandlePinRequests(wg *sync.WaitGroup) {
	defer wg.Done()
	for m := range p.Messages {
		if m == nil {
			p.logger.WithFields(log.Fields{"func": "HandlePinRequests"}).Warn("got nil message")
			break
		}

		// primitive throttling
		p.rateLimit <- struct{}{}
		p.logger.WithFields(log.Fields{"func": "HandlePinRequests"}).Debug("NEW MESSAGE")

		req := &queue.PinRequest{}
		if err := req.Unmarshal(m); err != nil {
			p.logger.WithFields(log.Fields{"func": "HandlePinRequests"}).Error(err)
			continue
		}
		wg.Add(1)
		go p.handle(wg, req)
	}
	p.logger.WithFields(log.Fields{"func": "HandlePinRequests"}).Debug("terminating")
}

// handle will Add or Unpin an object with regards to req.Unpin field.
func (p *Pinner) handle(wg *sync.WaitGroup, req *queue.PinRequest) {
	defer wg.Done()
	defer func() {
		<-p.rateLimit // primitive throttling
	}()

	if req.Unpin {
		p.handleUnpin(req)
	} else {
		p.handleAdd(req)
	}
}

func (p *Pinner) handleAdd(req *queue.PinRequest) {
	hash, err := p.Add(req)
	if err != nil {
		p.logger.WithFields(
			log.Fields{
				"message": "failed to add",
				"request": fmt.Sprintf("%#v", req),
			}).Error(err)
		return
	}
	p.logger.WithFields(
		log.Fields{
			"cid":     hash,
			"request": fmt.Sprintf("%#v", req),
		}).Debug("added pin request")

	stage := model.UploadStageIPFS
	s := &model.Storage{
		ID:          int64(req.StorageID),
		Cid:         &hash,
		UploadStage: &stage,
		UpdatedAt:   time.Now(),
		Pinned:      true,
	}
	// check errors
	if res := p.repo.Model(s).Updates(*s); res.Error != nil {
		p.logger.WithFields(
			log.Fields{
				"message": "failed to save",
				"request": fmt.Sprintf("%#v", req),
			}).Error(err)
		return
	}
	p.logger.WithFields(
		log.Fields{
			"cid":     hash,
			"request": fmt.Sprintf("%#v", req),
		}).Info("processed pin request")
}

func (p *Pinner) handleUnpin(req *queue.PinRequest) {
	err := p.Unpin(req)
	if err != nil {
		p.logger.WithFields(
			log.Fields{
				"message": "failed to unpin",
				"request": fmt.Sprintf("%#v", req),
			}).Error(err)
		return
	}

	stage := model.UploadStageStorage
	s := &model.Storage{
		ID:          int64(req.StorageID),
		UploadStage: &stage,
		UpdatedAt:   time.Now(),
		Pinned:      false,
	}

	if res := p.repo.Model(s).
		Select("upload_stage", "updated_at", "pinned").
		Updates(*s); res.Error != nil {
		p.logger.WithFields(
			log.Fields{
				"message": "failed to save",
				"request": fmt.Sprintf("%#v", req),
			}).Error(err)
		return
	}
	p.logger.WithFields(
		log.Fields{
			"cid": req.CID,
		}).Info("processed unpin request")
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
		p.logger.WithFields(log.Fields{
			"temp":    tmp,
			"request": fmt.Sprintf("%#v", pr),
		}).Debug("created temp dir")

		hash, err := p.sh.AddDir(tmp)
		if err != nil {
			return hash, err
		}

		p.logger.WithFields(log.Fields{
			"temp": tmp,
		}).Debug("removing temp dir")
		err = os.RemoveAll(tmp)
		if err != nil {
			p.logger.WithFields(log.Fields{
				"temp":    tmp,
				"message": "error removing temp dir",
			}).Error(err)
		}
		return hash, nil
	}
	obj, err := p.sm.GetObject(ctx, p.conf.Bucket, pr.Key, minio.GetObjectOptions{})
	if err != nil {
		p.logger.WithFields(log.Fields{
			"request": fmt.Sprintf("%#v", pr),
			"message": "error streaming object",
		}).Error(err)
		return "", err
	}
	return p.sh.Add(obj)
}

// Remove unpins file or directory identified by CID
// so it can be garbage collected in the next GC cycle.
func (p *Pinner) Unpin(pr *queue.PinRequest) error {
	if pr.CID == "" {
		err := errors.New("missing CID")
		p.logger.WithFields(log.Fields{
			"request": fmt.Sprintf("%#v", pr),
			"message": "cannot unpin objects without CID",
		}).Error(err)
		return err
	}

	p.logger.WithFields(log.Fields{
		"request": fmt.Sprintf("%#v", pr),
	}).Debug("unpinning object")

	return p.sh.Unpin(pr.CID)
}
