package xqiniu

import (
	"github.com/qiniu/go-sdk/v7/auth"
	"github.com/qiniu/go-sdk/v7/storage"
)

type (
	QBox struct {
		mac          *auth.Credentials
		uploader     *storage.FormUploader
		bucketMgr    *storage.BucketManager
		operationMgr *storage.OperationManager
	}
)

func NewMac(ak, sk string) *auth.Credentials {
	return auth.New(ak, sk)
}

func NewQBox(ak, sk string) *QBox {
	mac := auth.New(ak, sk)
	cfg := storage.Config{UseHTTPS: true}

	return &QBox{
		mac:          mac,
		uploader:     storage.NewFormUploader(nil),
		bucketMgr:    storage.NewBucketManager(mac, &cfg),
		operationMgr: storage.NewOperationManager(mac, &cfg),
	}
}
