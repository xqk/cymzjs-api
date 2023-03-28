package logic

import (
	"context"
	"github.com/xqk/cymzjs-api/pkg/xqiniu"

	"github.com/xqk/cymzjs-api/api/cymzjs/internal/svc"
	"github.com/xqk/cymzjs-api/api/cymzjs/internal/types"

	"git.zc0901.com/go/god/lib/logx"
)

type QiniuTokenLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewQiniuTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) QiniuTokenLogic {
	return QiniuTokenLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *QiniuTokenLogic) QiniuToken(req types.QiniuTokenReq) (string, error) {
	var uploadToken string
	var ak string
	var sk string
	var bucket xqiniu.Bucket
	switch req.Bucket {
	case xqiniu.BucketImolaStatic.Int64():
		ak = l.svcCtx.Config.QiniuZhuKe.AK
		sk = l.svcCtx.Config.QiniuZhuKe.SK
		bucket = xqiniu.BucketImolaStatic

	}

	if ak != "" && sk != "" && bucket != 0 {
		uploadToken = xqiniu.UploadToken(ak, sk, bucket, 7200)
	}

	return uploadToken, nil
}
