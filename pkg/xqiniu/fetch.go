package xqiniu

import "github.com/qiniu/go-sdk/v7/storage"

func (box *QBox) Fetch(url string, bucket Bucket, keys ...string) (storage.FetchRet, error) {
	if len(keys) > 0 && keys[0] != "" {
		return box.bucketMgr.Fetch(url, bucket.Name(), keys[0])
	} else {
		return box.bucketMgr.FetchWithoutKey(url, bucket.Name())
	}
}
