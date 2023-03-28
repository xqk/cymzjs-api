package logic

import (
	"context"
	"encoding/json"
	"github.com/xqk/cymzjs-api/api/cymzjs/internal/logicutil"

	"github.com/xqk/cymzjs-api/api/cymzjs/internal/svc"
	"github.com/xqk/cymzjs-api/api/cymzjs/internal/types"

	"git.zc0901.com/go/god/lib/logx"
)

type FindSystemFeteScenePlatformLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFindSystemFeteScenePlatformLogic(ctx context.Context, svcCtx *svc.ServiceContext) FindSystemFeteScenePlatformLogic {
	return FindSystemFeteScenePlatformLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FindSystemFeteScenePlatformLogic) FindSystemFeteScenePlatform(req types.FindSystemFeteSceneReq) (*types.FindSystemFeteSceneResp, error) {
	var resp *types.FindSystemFeteSceneResp

	grjResp, err := logicutil.DemoApiGrjResp(l.svcCtx, "/gurenju/wechat/system/feteScene/findSystemFeteScene?objectId=71D03B0337DD457E8CC7908C431A815D&objectType=platform")
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

	return resp, nil
}
