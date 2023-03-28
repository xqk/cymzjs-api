package logic

import (
	"context"
	"git.zc0901.com/go/god/lib/gconv"
	"github.com/xqk/cymzjs-api/api/cymzjs/internal/logicutil"
	"github.com/xqk/cymzjs-api/pkg"

	"github.com/xqk/cymzjs-api/api/cymzjs/internal/svc"
	"github.com/xqk/cymzjs-api/api/cymzjs/internal/types"

	"git.zc0901.com/go/god/lib/logx"
)

const (
	sqlUserFirstMemorialId = "SELECT id FROM memorial WHERE user_id=?"
)

type GetMemorialFormLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetMemorialFormLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetMemorialFormLogic {
	return GetMemorialFormLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetMemorialFormLogic) GetMemorialForm(req types.GetMemorialFormReq, authorization string) (*types.MemorialForm, error) {
	secret := l.svcCtx.Config.Auth.AccessSecret
	claims, err := pkg.JwtDecode(authorization, secret)
	var userId int64
	if err == nil {
		userId = gconv.Int64(claims["userId"])
	}

	var ids []int64
	err = l.svcCtx.MemorialModel.QueryNoCache(&ids, sqlUserFirstMemorialId, userId)
	if err != nil {
		return nil, err
	}

	if len(ids) == 0 {
		return &types.MemorialForm{}, nil
	}

	one, err := l.svcCtx.MemorialModel.FindOne(ids[0])
	if err != nil {
		return nil, err
	}

	visitCount, err := logicutil.CacheGetVisitCount(l.svcCtx.Store)
	if err != nil {
		visitCount = 0
	}
	return &types.MemorialForm{
		HeadImg:     one.HeadImg,
		NameOne:     one.NameOne,
		RelationOne: one.RelationOne,
		NameTwo:     one.NameTwo,
		RelationTwo: one.RelationTwo,
		VisitCount:  visitCount,
	}, nil
}
