package xqiniu

import (
	"github.com/qiniu/go-sdk/v7/storage"
	"strings"
)

const (
	ReturnBody = `
		{
			"key": $(key),
			"size": $(fsize),
			"type": $(mimeType),
			"hash": $(etag),
			"width": $(imageInfo.width),
			"height": $(imageInfo.height),
			"orientation": $(imageInfo.Orientation.val),
			"color": $(exif.ColorSpace.val),
			"videoDuration": $(avinfo.video.duration),
			"videoWidth": $(avinfo.video.width),
			"videoHeight": $(avinfo.video.height),
			"videoRotate": $(avinfo.video.tags.rotate),
			"url": "{host}/$(key)?imageView2/2/w/600",
			"originalUrl": "{host}/$(key)"
		}
	`
)

func UploadToken(ak, sk string, bucket Bucket, expires ...uint64) string {
	expire := uint64(8600)
	if len(expires) > 0 && expires[0] > 0 {
		expire = expires[0]
	}

	policy := storage.PutPolicy{
		Scope:      bucket.Name(),
		Expires:    expire,
		ReturnBody: strings.ReplaceAll(ReturnBody, "{host}", bucket.Host("https")),
	}
	mac := NewMac(ak, sk)

	return policy.UploadToken(mac)
}

func (box *QBox) UploadToken(bucket Bucket, expires ...uint64) string {
	if box.mac == nil {
		return ""
	}

	expire := uint64(8600)
	if len(expires) > 0 && expires[0] > 0 {
		expire = expires[0]
	}

	policy := storage.PutPolicy{
		Scope:      bucket.Name(),
		Expires:    expire,
		ReturnBody: strings.ReplaceAll(ReturnBody, "{host}", bucket.Host("https")),
	}

	return policy.UploadToken(box.mac)
}
