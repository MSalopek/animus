package ipfs

import "time"

// IPFSClient handles file manipulation on IPFS by interacting with
// an IPFS Node's HTTP-RPC API
type IPFSClient struct {
	NodeApiURL  string        // API URL of the IPFS Node used for file uploads
	CIDOnly     bool          // return only CID on upload completion
	PinInterval time.Duration // pin unpinned files every PinInterval
	AutoPin     bool          // pin files automatically
	AutoPurge   bool          // automatically delete all locally cached files/dirs after storing

	// TODO: make this use queues or Redis for communication
	// make requests retryable etc.
	UploadPaths  chan string // read from this chan and upload to IPFS
	UploadFolder string      // folder where the files to upload are located
}
