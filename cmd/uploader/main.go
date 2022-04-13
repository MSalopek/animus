package main

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	shell "github.com/ipfs/go-ipfs-api"
)

func main() {
	router := gin.Default()
	router.POST("/upload", upload)

	router.Run(":8083")
}

func upload(c *gin.Context) {
	meta := c.PostForm("meta")

	// Source
	file, err := c.FormFile("file")
	if err != nil {
		c.String(http.StatusBadRequest, "get form err: %s", err.Error())
		return
	}

	// filename := filepath.Base(file.Filename)
	src, err := file.Open()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	defer src.Close()
	hash, err := UploadToIpfs(src)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	// if err := c.SaveUploadedFile(file, filename); err != nil {
	// 	c.String(http.StatusBadRequest, "upload file err: %s", err.Error())
	// 	return
	// }

	c.String(http.StatusOK, "File %s uploaded successfully with fields meta=%s and hash=%s.", file.Filename, meta, hash)
}

func UploadToIpfs(r io.Reader) (string, error) {
	sh := shell.NewShell("localhost:5001")
	fileHash, err := sh.Add(r)
	if err != nil {
		return "", err
	}
	return fileHash, nil
}
