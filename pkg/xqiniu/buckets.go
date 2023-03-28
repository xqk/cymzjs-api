package xqiniu

import "git.zc0901.com/go/god/lib/container/garray"

type (
	// Bucket 标志七牛的存储空间
	Bucket int64

	// BucketName 标志七牛存储空间的名称
	BucketName string

	// BucketDomain 标志七牛存储空间的域名
	BucketDomain string
)

const (
	BucketImolaStatic Bucket = 1
	BucketNestStatic  Bucket = 2
	BucketNestPhoto   Bucket = 3

	BucketNameImolaStatic BucketName = "zhuke-static"
	BucketNameNestStatic  BucketName = "dome-static"
	BucketNameNestPhoto   BucketName = "nestciao-photo"

	BucketDomainImolaStatic BucketDomain = "static.imolacn.com"
	BucketDomainNestStatic  BucketDomain = "static.zc0901.com"
	BucketDomainNestPhoto   BucketDomain = "img.imolacn.com"
)

var (
	bucketNameMaps = map[Bucket]BucketName{
		BucketImolaStatic: BucketNameImolaStatic,
		BucketNestStatic:  BucketNameNestStatic,
	}
	bucketDomainMaps = map[Bucket]BucketDomain{
		BucketImolaStatic: BucketDomainImolaStatic,
		BucketNestStatic:  BucketDomainNestStatic,
	}
	buckets = garray.NewArrayFrom([]interface{}{
		BucketImolaStatic.Int64(),
		BucketNestStatic.Int64(),
	})
)

func (bucket Bucket) Int64() int64 {
	return int64(bucket)
}

func (bucket Bucket) Name() string {
	bucketName, ok := bucketNameMaps[bucket]
	if !ok {
		return ""
	}
	return bucketName.String()
}

func (bucket Bucket) Domain() string {
	bucketDomain, ok := bucketDomainMaps[bucket]
	if !ok {
		return ""
	}
	return bucketDomain.String()
}

func (bucket Bucket) Host(schemes ...string) string {
	scheme := "https"
	if len(schemes) > 0 && schemes[0] != "" {
		scheme = schemes[0]
	}
	return scheme + "://" + bucket.Domain()
}

func (bucketName BucketName) String() string {
	return string(bucketName)
}

func (bucketDomain BucketDomain) String() string {
	return string(bucketDomain)
}
