package logic

import (
	"context"
	"encoding/json"
	"github.com/xqk/cymzjs-api/api/cymzjs/internal/logicutil"

	"github.com/xqk/cymzjs-api/api/cymzjs/internal/svc"
	"github.com/xqk/cymzjs-api/api/cymzjs/internal/types"

	"git.zc0901.com/go/god/lib/logx"
)

type FindPlatformDetailsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFindPlatformDetailsLogic(ctx context.Context, svcCtx *svc.ServiceContext) FindPlatformDetailsLogic {
	return FindPlatformDetailsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FindPlatformDetailsLogic) FindPlatformDetails(req types.FindPlatformDetailsReq) (*types.FindPlatformDetailsResp, error) {
	var resp *types.FindPlatformDetailsResp

	grjResp, err := logicutil.DemoApiGrjResp(l.svcCtx, "/gurenju/wechat/platform/platform/findPlatformDetails?idOrNo=71D03B0337DD457E8CC7908C431A815D")
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
