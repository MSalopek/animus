package pinner

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	shell "github.com/ipfs/go-ipfs-api"
	"github.com/minio/minio-go/v7"

	"github.com/msalopek/animus/engine/repo"
	"github.com/msalopek/animus/model"
	"github.com/msalopek/animus/queue"
	"github.com/msalopek/animus/storage"

	log "github.com/sirupsen/logrus"
)

const (
	WebhookFail    = "failed"
	WebhookSuccess = "success"
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

	Messages chan []byte
	Webhooks chan []byte

	subscriber queue.Subscriber
	webhookSub queue.Subscriber
	rateLimit  chan struct{}

	publisher queue.Publisher

	logger *log.Logger
	conf   *Config
}

func New(conf *Config, repo *repo.Repo, logger *log.Logger) *Pinner {
	p := &Pinner{
		autoPin:  true,
		Messages: make(chan []byte),
		Webhooks: make(chan []byte),

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

	p.publisher = queue.MustNewPublisher(conf.PublishTopic, conf.NsqdURL)
	p.subscriber = queue.MustNewSubscriber(
		conf.SubscribeTopic,
		"pinner",
		conf.NsqLookupdURL,
		p.Messages)
	p.webhookSub = queue.MustNewSubscriber(
		"webhooks",
		"pinner",
		conf.NsqLookupdURL,
		p.Webhooks)

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
	p.webhookSub.Unsubscribe()
	p.publisher.Stop()
	close(p.Messages)
	close(p.Webhooks)
	close(p.rateLimit)
}

func (p *Pinner) SetAutopin(val bool) {
	p.autoPin = val
}

// Sends scheduled webhooks.
// Upon termination all in flight requests will be completed.
func (p *Pinner) HandleWebhooks(wg *sync.WaitGroup) {
	defer wg.Done()

	for m := range p.Webhooks {
		if m == nil {
			p.logger.WithFields(log.Fields{"action": "HandleWebhooks"}).Warn("got nil message")
			break
		}
		p.logger.WithFields(log.Fields{"action": "HandleWebhooks"}).Debug("new message")

		req := &queue.WebhookMessage{}
		if err := req.Unmarshal(m); err != nil {
			p.logger.WithFields(log.Fields{"action": "HandlePinRequests"}).Error(err)
			continue
		}

		wg.Add(1)
		go p.SendWebhook(wg, req)
	}
}

// Read and process PinRequests from Messages chan.
// All in-flight requests will be processed before terminating.
func (p *Pinner) HandlePinRequests(wg *sync.WaitGroup) {
	defer wg.Done()
	for m := range p.Messages {
		if m == nil {
			p.logger.WithFields(log.Fields{"action": "HandlePinRequests"}).Warn("got nil message")
			break
		}

		// primitive throttling
		p.rateLimit <- struct{}{}
		p.logger.WithFields(log.Fields{"action": "HandlePinRequests"}).Debug("new message")

		req := &queue.PinRequest{}
		if err := req.Unmarshal(m); err != nil {
			p.logger.WithFields(log.Fields{"action": "HandlePinRequests"}).Error(err)
			continue
		}
		wg.Add(1)
		go p.handle(wg, req)
	}
	p.logger.WithFields(log.Fields{"action": "HandlePinRequests"}).Debug("terminating")
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

// adds file specified by req.StorageID thereby pinning it to the node
// webhooks are sent for both successful and unsuccessful addition
func (p *Pinner) handleAdd(req *queue.PinRequest) {
	storage, err := p.repo.GetUserUploadByID(req.UserID, req.StorageID)
	if err != nil {
		p.logger.WithFields(
			log.Fields{
				"action":  "handleAdd",
				"message": "failed to fetch record",
				"request": fmt.Sprintf("%#v", req),
			}).Error(err)
		return
	}

	hash, err := p.Add(req)
	if err != nil {
		p.logger.WithFields(
			log.Fields{
				"action":  "handleAdd",
				"message": "failed to add",
				"request": fmt.Sprintf("%#v", req),
			}).Error(err)

		p.publishWebhook(storage, WebhookFail)
		return
	}

	p.logger.WithFields(
		log.Fields{
			"action":  "handleAdd",
			"cid":     hash,
			"request": fmt.Sprintf("%#v", req),
		}).Debug("pinning successful")

	stage := model.UploadStageIPFS
	storage.Cid = &hash
	storage.UploadStage = &stage
	storage.UpdatedAt = time.Now()
	storage.Pinned = true
	// fail silently and don't send webhooks
	// the upload was successful but we didn't store it in DB
	if res := p.repo.Model(storage).
		Select("upload_stage", "cid", "updated_at", "pinned").
		Updates(*storage); res.Error != nil {

		p.logger.WithFields(
			log.Fields{
				"message": "failed to save after pinning",
				"request": fmt.Sprintf("%#v", req),
			}).Error(err)

		return
	}

	p.publishWebhook(storage, WebhookSuccess)

	p.logger.WithFields(
		log.Fields{
			"cid":     hash,
			"request": fmt.Sprintf("%#v", req),
		}).Info("success")
}

// unpins the file specified by req.Cid
// webhooks are not sent when unpinning
func (p *Pinner) handleUnpin(req *queue.PinRequest) {
	err := p.Unpin(req)
	if err != nil {
		p.logger.WithFields(
			log.Fields{
				"action":  "handleUnpin",
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
				"action":  "handleUnpin",
				"message": "failed to save",
				"request": fmt.Sprintf("%#v", req),
			}).Error(err)
		return
	}

	p.logger.WithFields(
		log.Fields{
			"action": "handleUnpin",
			"cid":    req.CID,
		}).Info("success")
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

// Send HTTP request to provided URL.
func (p *Pinner) SendWebhook(wg *sync.WaitGroup, m *queue.WebhookMessage) error {
	defer wg.Done()

	req := WebhookRequest{
		Object: m,
		Status: m.Status,
	}

	// TODO: maybe cache user webhook URLs if DB is under heavly load.
	u, err := p.repo.GetUserById(int(m.UserID))
	if err != nil {
		p.logger.WithFields(log.Fields{
			"action":    "SendWebhook",
			"userID":    m.UserID,
			"storageID": m.StorageID,
			"message":   "user not found",
		}).Error(err)
		return err
	}

	p.logger.WithFields(log.Fields{
		"action":    "SendWebhook",
		"userID":    m.UserID,
		"storageID": m.StorageID,
	}).Debug("new message")

	if u.WebhooksActive == nil || !*u.WebhooksActive {
		return nil
	}

	if u.WebhooksURL == nil || *u.WebhooksURL == "" {
		p.logger.WithFields(log.Fields{
			"action":    "SendWebhook",
			"message":   fmt.Sprintf("user WebhookURL not configured: %d", m.UserID),
			"storageID": m.StorageID,
		}).Error(err)
		return fmt.Errorf("user WebhookURL not configured")
	}

	if req.Status == WebhookFail {
		req.RetryUrl = fmt.Sprintf("%s/auth/storage/pin/id/%d", p.conf.RetryUrl, m.StorageID)
	}

	body, err := json.Marshal(req)
	if err != nil {
		return err
	}

	buf := bytes.NewBuffer(body)

	_, err = http.Post(*u.WebhooksURL, "application/json", buf)
	if err != nil {
		p.logger.WithFields(log.Fields{
			"action":    "SendWebhook",
			"message":   "could not send webhook post request",
			"storageID": m.StorageID,
		}).Error(err)
		return err
	}
	p.logger.WithFields(log.Fields{
		"action":    "SendWebhook",
		"storageID": m.StorageID,
	}).Debug("success")
	return nil
}

func (p *Pinner) publishWebhook(s *model.Storage, status string) {
	w := queue.WebhookMessage{
		StorageID:   s.ID,
		Cid:         s.Cid,
		Name:        s.Name,
		Metadata:    s.Metadata,
		UploadStage: s.UploadStage,
		Pinned:      s.Pinned,
		CreatedAt:   s.CreatedAt,
		UpdatedAt:   s.UpdatedAt,
		DeletedAt:   s.DeletedAt,

		Status: status,

		// needed downstream to fetch user WebhookURL
		UserID: s.UserID,
	}

	body, err := w.Marshal()
	if err != nil {
		p.logger.WithFields(log.Fields{
			"action":    "publishWebhook",
			"message":   "could not marshall webhook message",
			"storageID": s.ID,
		}).Error(err)
	}

	err = p.publisher.Publish(p.conf.PublishTopic, body)
	if err != nil {
		p.logger.WithFields(log.Fields{
			"action":    "publishWebhook",
			"message":   "could not publish webhook message",
			"storageID": s.ID,
		}).Error(err)
	}
	p.logger.WithFields(log.Fields{
		"action":    "publishWebhook",
		"storageID": s.ID,
	}).Info("success")
}

type WebhookRequest struct {
	Object   *queue.WebhookMessage `json:"object"`
	Status   string                `json:"status"`
	RetryUrl string                `json:"retry_url"`
}
