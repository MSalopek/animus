package pinner

import (
	"context"
	"log"
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
)

type Pinner struct {
	sm *storage.Manager // s3 compatible storage

	sh      *shell.Shell // ipfs api shell
	autoPin bool         // pin files automatically

	repo *repo.Repo

	consumer *nsq.Consumer
	producer *nsq.Producer

	conf *Config
	reqC chan *queue.PinRequest
}

func New(conf Config) *Pinner {
	p := &Pinner{
		conf: &conf,
		sm:   storage.MustNewManager(conf.Storage),
		sh:   shell.NewShell(conf.NodeApiURL),
		reqC: make(chan *queue.PinRequest, conf.QueueSize),
	}
	if conf.QueueSize < 1 {
		// this will process 1 request at a time
		p.reqC = make(chan *queue.PinRequest)
	}

	if conf.LocalShell {
		p.sh = shell.NewLocalShell()
	}

	config := nsq.NewConfig()

	consumer, err := nsq.NewConsumer(conf.SubscribeTopic, "pinner", config)
	consumer.AddHandler(nsq.HandlerFunc(p.HandleMessages))

	// Use nsqlookupd to discover nsqd instances.
	// See also ConnectToNSQD, ConnectToNSQDs, ConnectToNSQLookupds.
	err = consumer.ConnectToNSQLookupd(conf.NsqLookupdURL)
	if err != nil {
		log.Fatal(err)
	}
	if err != nil {
		log.Fatal(err)
	}

	return p
}

func (p *Pinner) SetAutopin(val bool) {
	p.autoPin = val
}

func (p *Pinner) HandleMessages(m *nsq.Message) error {
	pr := &queue.PinRequest{}
	if err := pr.Unmarshal(m.Body); err != nil {
		return err
	}
	p.reqC <- pr
	return nil

}

func (p *Pinner) handlePinRequests(wg *sync.WaitGroup, reqs <-chan *queue.PinRequest) {
	for req := range reqs {
		wg.Add(1)
		go p.handleAdd(req)
	}
}

func (p *Pinner) handleAdd(req *queue.PinRequest) {
	hash, err := p.Add(req)
	if err != nil {
		// log error
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
		hash, err := p.sh.AddDir(tmp)
		if err != nil {
			return hash, err
		}
		// check err
		os.RemoveAll(tmp)
		return hash, nil
	}
	obj, err := p.sm.GetObject(ctx, p.conf.Bucket, pr.Key, minio.GetObjectOptions{})
	return p.sh.Add(obj)
}
