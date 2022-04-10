package service

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/msalopek/animus/api"
	"github.com/msalopek/animus/ipfs"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type HttpAPI interface {
	Login(c *gin.Context)
	Register(c *gin.Context)
	AddFile(c *gin.Context)
	PinFile(c *gin.Context)
	DeleteFile(c *gin.Context)

	AddGate(c *gin.Context)
	DeleteGate(c *gin.Context)
	UpdateGate(c *gin.Context)

	Ping(c *gin.Context)
}

// IPFSProxy service accepts HTTP file upload requests, stores the
// received files to a local cache and uploads them to IPFS.
// The service can keep track of file changes (WIP, and very primitive)
// which allows the user to manipulate files.
//
// The service has 2 working parts:
//
// 1. HTTP-API
//    - allows the service to restrict user access (users must register etc.)
//    - accepts file uploads and stores them to disk
//    - makes file manipulation easier
//    - when a file is cached locally, it is queued for IPFS upload
//
// 2. IPFS Client
//    - uses IPFS Node's HTTP-RPC to send file operation requests
//    - checks uploaded files and uploads them to IPFS
//    - can be used to manipulate the file/dir on IPFS
type IPFSProxy struct {
	db         *gorm.DB
	httpAPI    *api.HttpAPI
	ipfsClient *ipfs.IPFSClient

	logger *logrus.Logger
}

func New(db *gorm.DB, *api.HttpAPI, ipfsNodeAPI string) IPFSProxy {
	if ipfsNodeAPI == "" {
		panic("IPFSNodeAPI must be provided")
	}

	logger := log.New()
	logger.Out = os.Stdout
	logger.SetFormatter(&log.TextFormatter{})

	s := IPFSProxy{
		httpAPI: httpAPi,
		db: db,
		logger:  logger,
	}
	return s
}
