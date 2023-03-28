package logicutil

import (
	"git.zc0901.com/go/god/lib/gconv"
	"git.zc0901.com/go/god/lib/store/kv"
)

const (
	keyVisitCount = "cymzjs:visit:count"
)

// CacheGetVisitCount 获取访问次数
func CacheGetVisitCount(store kv.Store) (count int64, err error) {
	countStr, err := store.Get(keyVisitCount)
	if err != nil {
		return count, err
	}

	if countStr != "" {
		count = gconv.Int64(countStr)
	}

	return count, nil
}

// CacheSetVisitCount 获取访问次数
func CacheSetVisitCount(store kv.Store) (count int64, err error) {
	count, err = store.IncrBy(keyVisitCount, 1)
	if err != nil {
		return count, err
	}
	return count, nil
}
