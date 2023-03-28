package logic

import (
	"context"
	"encoding/json"
	"github.com/xqk/cymzjs-api/api/cymzjs/internal/logicutil"
	"github.com/xqk/cymzjs-api/api/cymzjs/internal/types"

	"git.zc0901.com/go/god/lib/logx"
	"github.com/xqk/cymzjs-api/api/cymzjs/internal/svc"
)

type FindPlatformMemorialConfigLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFindPlatformMemorialConfigLogic(ctx context.Context, svcCtx *svc.ServiceContext) FindPlatformMemorialConfigLogic {
	return FindPlatformMemorialConfigLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FindPlatformMemorialConfigLogic) FindPlatformMemorialConfig() (resp []*types.FindPlatformMemorialConfigItem, err error) {
	grjResp, err := logicutil.DemoApiGrjResp(l.svcCtx, "/gurenju/wechat/platform/platform/findPlatformMemorialConfig?platformId=71D03B0337DD457E8CC7908C431A815D")
	if err != nil {
		return nil, err
	}

	dataBytes, err := json.Marshal(grjResp.Data)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(dataBytes, &resp)
	if err != nil {
		return nil, err
	}

	return
}
