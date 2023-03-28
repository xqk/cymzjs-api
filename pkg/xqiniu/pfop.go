package xqiniu

import "strings"

// Pfop 持久化数据处理
func (box *QBox) Pfop(bucket Bucket, key string, fops []string, pipeline, notifyUrl string, force bool) (string, error) {
	return box.operationMgr.Pfop(bucket.Name(), key, strings.Join(fops, ";"), pipeline, notifyUrl, force)
}
