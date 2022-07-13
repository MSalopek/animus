package storage

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"os"
	"path"
	"strings"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/minio/minio-go/v7/pkg/lifecycle"
)

// TODO: write some integration tests for this
type Opts struct {
	Region       string
	LifetimeDays int
}

type Manager struct {
	*minio.Client
	conf  Config
	rules *lifecycle.Configuration // defines bucket lifecycle
}

// Basic upload metadata.
type Upload struct {
	Bucket       string    `json:"bucket,omitempty"`
	Key          string    `json:"key,omitempty"`
	Location     string    `json:"bucket,omitempty"`
	LastModified time.Time `json:"last_modified,omitempty"`
	Size         int64     `json:"size,omitempty"`
}

type Uploads struct {
	Uploads []*Upload `json:"uploads"`
}

// FileStreamer can stream from reader and upload to target bucket.
type FileStreamer interface {
	StreamFile(ctx context.Context, bucket, fname string, reader io.Reader, size int64, opts Opts) (*Upload, error)
}

// UploadFile is a convenience method for uploading file using FileStreamer interface using default storage options.
func UploadFile(fs FileStreamer, file *multipart.FileHeader, bucket, objName string) (*Upload, error) {
	src, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	// TODO: add deadline
	ctx := context.Background()
	return fs.StreamFile(ctx, bucket, objName, src, file.Size, Opts{})
}

// Initializes new manager with provided configuration.
func MustNewManager(conf Config) *Manager {
	c, err := minio.New(conf.Endpoint,
		&minio.Options{
			Creds: credentials.NewStaticV4(
				conf.AccessKeyID,
				conf.AccessKeySecret,
				conf.AccessToken,
			),
		},
	)
	if err != nil {
		log.Fatalf("could not start minIO client: %s", err)
	}

	lcr := lifecycle.NewConfiguration()
	lcr.Rules = []lifecycle.Rule{
		{
			ID:     "expire-bucket",
			Status: "Enabled",
			Expiration: lifecycle.Expiration{
				Days: lifecycle.ExpirationDays(conf.Lifetime),
			},
		},
	}

	if conf.TmpDir != "" {
		if _, err := os.Stat(conf.TmpDir); os.IsNotExist(err) {
			err := os.Mkdir(conf.TmpDir, 0700)
			if err != nil {
				log.Fatalf("could not create temp directory %s - err: %s", conf.TmpDir, err)
			}
		}
	}

	return &Manager{
		Client: c,
		conf:   conf,
		rules:  lcr,
	}
}

// Create bucket with defined lifetime rules and ignore errors raised if bucket already exists.
func (m *Manager) SafeCreateExpirableBucket(ctx context.Context, bucket, region string) error {
	err := m.MakeBucket(ctx, bucket, minio.MakeBucketOptions{Region: region})
	if err != nil {
		exists, errBucketExists := m.BucketExists(ctx, bucket)
		if errBucketExists == nil && exists {
			return nil
		}
		return err
	}

	if err := m.SetBucketLifecycle(ctx, bucket, m.rules); err != nil {
		return err
	}

	return nil
}

// Create bucket and ignore errors raised if bucket already exists.
func (m *Manager) SafeCreateBucket(ctx context.Context, bucket, region string) error {
	err := m.MakeBucket(ctx, bucket, minio.MakeBucketOptions{Region: region})
	if err != nil {
		exists, errBucketExists := m.BucketExists(ctx, bucket)
		if errBucketExists == nil && exists {
			return nil
		}
		return err
	}

	return nil
}

// Recursively upload directory to bucket. Create bucket if it does not exist.
// In case of errors on upload the content will be partially uploaded.
func (m *Manager) FSUploadDir(ctx context.Context, path, bucket string, opts Opts) error {
	loc := m.conf.Region
	if opts.Region != "" {
		loc = opts.Region
	}

	finfo, err := ReadDir(path)
	if err != nil {
		return err
	}

	if err := m.SafeCreateBucket(ctx, bucket, loc); err != nil {
		return err
	}

	for _, f := range finfo {
		_, err := m.FPutObject(ctx, bucket, f.Name, f.Path, minio.PutObjectOptions{})
		if err != nil {
			return err
		}
	}
	return nil
}

// Recursively upload directory to expireable bucket. Create bucket if it does not exist.
// In case of errors on upload the content will be partially uploaded.
func (m *Manager) FSUploadExpireableDir(ctx context.Context, path, bucket string, opts Opts) error {
	loc := m.conf.Region
	if opts.Region != "" {
		loc = opts.Region
	}

	finfo, err := ReadDir(path)
	if err != nil {
		return err
	}

	if err := m.SafeCreateBucket(ctx, bucket, loc); err != nil {
		return err
	}

	for _, f := range finfo {
		_, err := m.FPutObject(ctx, bucket, f.Name, f.Path, minio.PutObjectOptions{})
		if err != nil {
			return err
		}
	}
	return nil
}

// func (c *Client) PutObject(ctx context.Context, bucketName, objectName string, reader io.Reader, objectSize int64,
// 	opts PutObjectOptions,

// Upload single file from stream to bucket.
// On error file upload is aborted.
func (m *Manager) StreamFile(ctx context.Context, bucket, fname string, reader io.Reader, size int64, opts Opts) (*Upload, error) {
	loc := m.conf.Region
	if opts.Region != "" {
		loc = opts.Region
	}

	if err := m.SafeCreateBucket(ctx, bucket, loc); err != nil {
		return nil, err
	}

	info, err := m.PutObject(ctx, bucket, fname, reader, size, minio.PutObjectOptions{})
	if err != nil {
		return nil, err
	}

	return &Upload{
		Bucket:       info.Bucket,
		Key:          info.Key,
		Location:     info.Location,
		Size:         info.Size,
		LastModified: info.LastModified,
	}, nil
}

// Upload single file from stream to expirable bucket.
func (m *Manager) StreamExpirableFile(ctx context.Context, bucket, fname string, reader io.Reader, size int64, opts Opts) (*Upload, error) {
	loc := m.conf.Region
	if opts.Region != "" {
		loc = opts.Region
	}

	if err := m.SafeCreateExpirableBucket(ctx, bucket, loc); err != nil {
		return nil, err
	}

	info, err := m.PutObject(ctx, bucket, fname, reader, size, minio.PutObjectOptions{})
	if err != nil {
		return nil, err
	}

	return &Upload{
		Bucket:       info.Bucket,
		Key:          info.Key,
		Location:     info.Location,
		Size:         info.Size,
		LastModified: info.LastModified,
	}, nil
}

// Upload single file from filesystem to bucket.
func (m *Manager) FSUploadFile(ctx context.Context, bucket, path string, opts Opts) error {
	loc := m.conf.Region
	if opts.Region != "" {
		loc = opts.Region
	}

	f, err := os.Open(path)
	if err != nil {
		return err
	}
	f.Close()

	if err := m.SafeCreateBucket(ctx, bucket, loc); err != nil {
		return err
	}

	fname := fileName(path)
	if _, err := m.FPutObject(
		ctx, bucket, fname, path, minio.PutObjectOptions{}); err != nil {
		return err
	}

	return nil
}

// Upload single file from filesystem to expirable bucket.
func (m *Manager) FSUploadExpirableFile(ctx context.Context, bucket, path string, opts Opts) error {
	loc := m.conf.Region
	if opts.Region != "" {
		loc = opts.Region
	}

	f, err := os.Open(path)
	if err != nil {
		return err
	}
	f.Close()

	if err := m.SafeCreateExpirableBucket(ctx, bucket, loc); err != nil {
		return err
	}

	fname := fileName(path)
	if _, err := m.FPutObject(
		ctx, bucket, fname, path, minio.PutObjectOptions{}); err != nil {
		return err
	}

	return nil
}

// Recursively download bucket contents and return path to downloaded files.
func (m *Manager) downloadBucket(ctx context.Context, bucket string) (string, error) {
	dir, err := ioutil.TempDir(m.conf.TmpDir, m.conf.TmpDirPrefix)
	if err != nil {
		return "", err
	}

	for c := range m.ListObjects(ctx, bucket, minio.ListObjectsOptions{Recursive: true}) {
		err := m.DownloadFileToDir(ctx, bucket, dir, c.Key)
		if err != nil {
			return "", err
		}
	}
	return dir, nil
}

// Recursively download folder contents and return path to downloaded files.
// prefix represents a bucket prefix from which the recursive directory traversal will start.
func (m *Manager) DownloadDir(ctx context.Context, bucket, prefix string) (string, error) {
	if bucket == "" {
		return "", errors.New("bucket not provided")
	}
	if prefix == "" {
		return "", errors.New("path prefix not provided")
	}

	tmpDir, err := ioutil.TempDir(m.conf.TmpDir, m.conf.TmpDirPrefix)
	if err != nil {
		return "", err
	}

	for c := range m.ListObjects(ctx, bucket, minio.ListObjectsOptions{Recursive: true, Prefix: prefix}) {
		// ex. key: 1/some-filename.ext
		sKey := strings.SplitN(c.Key, "/", 2)
		if len(sKey) != 2 {
			return "", fmt.Errorf("invalid key: %s", c.Key)
		}
		fname := sKey[1]
		err := m.DownloadFileToDirWithName(ctx, bucket, tmpDir, c.Key, fname)
		if err != nil {
			return "", err
		}
	}

	return tmpDir, nil
}

// Recursively download directory contents, tar.gz on destination and return path to downloaded file.
// prefix represents a bucket prefix from which the recursive directory traversal will start.
func (m *Manager) DownloadDirTarGz(ctx context.Context, bucket, prefix string) (string, error) {
	if bucket == "" {
		return "", errors.New("bucket not provided")
	}
	if prefix == "" {
		return "", errors.New("folder path not provided")
	}

	tmpDir, err := ioutil.TempDir(m.conf.TmpDir, m.conf.TmpDirPrefix)
	if err != nil {
		return "", err
	}

	archive := path.Join(tmpDir, fmt.Sprintf("%s.tar.gz", prefix))
	fileW, err := os.Create(archive)
	if err != nil {
		return "", err
	}
	defer fileW.Close()

	zw := gzip.NewWriter(fileW)
	tw := tar.NewWriter(zw)
	for c := range m.ListObjects(ctx, bucket, minio.ListObjectsOptions{Recursive: true, Prefix: prefix}) {
		object, err := m.GetObject(ctx, bucket, c.Key, minio.GetObjectOptions{})
		if err != nil {
			return "", nil
		}
		fHead := &tar.Header{
			Name: c.Key,
			Mode: 0600,
			Size: c.Size,
		}
		if err := tw.WriteHeader(fHead); err != nil {
			return "", err
		}

		buf := make([]byte, 4096)
		for {
			n, rErr := object.Read(buf)
			if rErr != nil && rErr != io.EOF {
				return "", rErr
			}
			_, wErr := tw.Write(buf[:n])
			if wErr != nil {
				return "", wErr
			}

			if rErr == io.EOF {
				break
			}
		}
	}

	if err := tw.Close(); err != nil {
		return "", err
	}

	if err := zw.Close(); err != nil {
		return "", err
	}

	return archive, nil
}

// Download a single file from bucket and store it to a temp folder returning the path to tempfile.
func (m *Manager) DownloadFileToTmp(ctx context.Context, bucket, fname string) (string, error) {
	// just create dir, someone else needs to remove.
	dir, err := ioutil.TempDir(m.conf.TmpDir, m.conf.TmpDirPrefix)
	if err != nil {
		return "", err
	}

	fpath := path.Join(dir, fname)

	if err = m.FGetObject(ctx, bucket, fname, fpath, minio.GetObjectOptions{}); err != nil {
		return "", err
	}

	return fpath, nil
}

// Download a single file from bucket and store it to the dir returning path to the file.
func (m *Manager) DownloadFileToDir(ctx context.Context, bucket, dir, fname string) error {
	if _, err := os.Stat(dir); err != nil {
		return nil
	}

	fpath := path.Join(dir, fname)
	if err := m.FGetObject(ctx, bucket, fname, fpath, minio.GetObjectOptions{}); err != nil {
		return err
	}

	return nil
}

// Download a single file from bucket and store it to the dir returning path to the file.
// The file name will be set to fname. This is useful for striping folder prefixes.
func (m *Manager) DownloadFileToDirWithName(ctx context.Context, bucket, dir, key, fname string) error {
	if _, err := os.Stat(dir); err != nil {
		return nil
	}

	fpath := path.Join(dir, fname)
	if err := m.FGetObject(ctx, bucket, key, fpath, minio.GetObjectOptions{}); err != nil {
		return err
	}

	return nil
}

// RemoveDirObjects deletes storage objects specified by rmList.
func (m *Manager) RemoveDirObjects(ctx context.Context, bucket, dir string) error {
	objectsCh := make(chan minio.ObjectInfo)
	var err error

	// Send object names that are needed to be removed to objectsCh
	go func() {
		defer close(objectsCh)
		// List all objects from a bucket-name with a matching prefix.
		for object := range m.ListObjects(ctx, bucket, minio.ListObjectsOptions{Prefix: dir, Recursive: true}) {
			if object.Err != nil {
				// TODO: handle this!
				break
			}
			objectsCh <- object
		}
	}()

	// NOTE: don't know what this does
	opts := minio.RemoveObjectsOptions{
		GovernanceBypass: true,
	}

	// TODO: this only catches the last error
	for rErr := range m.RemoveObjects(ctx, bucket, objectsCh, opts) {
		err = rErr.Err
	}
	return err
}

// Basic file information.
type FileInfo struct {
	Name string
	Path string
}

// ReadDir recursively traverses dir defined by path.
// An error is returned if the path does not exist.
// If an error is encountered during traversal the dir will be
// partially read (up to the error).
func ReadDir(path string) ([]*FileInfo, error) {
	var res []*FileInfo
	return res, readDirRecursive(path, &res)
}

func readDirRecursive(path string, result *[]*FileInfo) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	base := strings.TrimLeft(path, "./")
	names, _ := file.Readdirnames(0)
	for _, name := range names {
		fp := fmt.Sprintf("%s/%s", path, name)
		f, err := os.Open(fp)
		if err != nil {
			return err
		}
		defer f.Close()
		fInfo, err := f.Stat()
		if err != nil {
			return err
		}
		if !fInfo.IsDir() {
			*result = append(*result, &FileInfo{
				Name: fmt.Sprintf("%s/%s", base, name),
				Path: fp,
			})
		} else {
			readDirRecursive(fp, result)
		}
	}
	return nil
}

func fileName(path string) string {
	sp := strings.Split(path, "/")
	return sp[len(sp)-1]
}
