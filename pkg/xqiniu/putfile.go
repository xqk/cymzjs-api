package xqiniu

import (
	"context"
	"crypto/sha1"
	"fmt"
	"github.com/qiniu/go-sdk/v7/storage"
	"io/ioutil"
)

type (
	UploadedRet struct {
		Key           string `json:"key"`
		Size          int64  `json:"size"`
		Type          string `json:"type"`
		Hash          string `json:"hash"`
		Width         int64  `json:"width"`
		Height        int64  `json:"height"`
		Orientation   string `json:"orientation"`
		Color         string `json:"color"`
		VideoRotate   string `json:"videoRotate"`
		VideoWidth    int64  `json:"videoWidth"`
		VideoHeight   int64  `json:"videoHeight"`
		VideoDuration string `json:"videoDuration"`
		Url           string `json:"url"`
		OriginalUrl   string `json:"originalUrl"`
	}
)

func PutFile(ak, sk, p string, bucket Bucket, keyPrefix ...string) (*UploadedRet, error) {
	fileB, err := ioutil.ReadFile(p)
	if err != nil {
		return nil, err
	}

	sha1B := sha1.Sum(fileB)
	key := fmt.Sprintf("%x", sha1B)
	if len(keyPrefix) > 0 && keyPrefix[0] != "" {
		key = fmt.Sprintf("%s/%s", keyPrefix[0], key)
	}

	uploader := storage.NewFormUploader(nil)
	token := UploadToken(ak, sk, bucket)
	r := new(UploadedRet)
	err = uploader.PutFile(context.Background(), &r, token, key, p, nil)

	return r, err
}

func (box *QBox) PutFile(p string, bucket Bucket, keyPrefix ...string) (*UploadedRet, error) {
	fileB, err := ioutil.ReadFile(p)
	if err != nil {
		return nil, err
	}

	sha1B := sha1.Sum(fileB)
	key := fmt.Sprintf("%x", sha1B)
	if len(keyPrefix) > 0 && keyPrefix[0] != "" {
		key = fmt.Sprintf("%s/%s", keyPrefix[0], key)
	}

	token := box.UploadToken(bucket)
	r := new(UploadedRet)
	err = box.uploader.PutFile(context.Background(), &r, token, key, p, nil)

	return r, err
}
