package logic

import (
	"context"
	"encoding/json"
	"github.com/xqk/cymzjs-api/api/cymzjs/internal/logicutil"

	"github.com/xqk/cymzjs-api/api/cymzjs/internal/svc"
	"github.com/xqk/cymzjs-api/api/cymzjs/internal/types"

	"git.zc0901.com/go/god/lib/logx"
)

type FindPlatformEnterpriseDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFindPlatformEnterpriseDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) FindPlatformEnterpriseDetailLogic {
	return FindPlatformEnterpriseDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FindPlatformEnterpriseDetailLogic) FindPlatformEnterpriseDetail(req types.FindPlatformEnterpriseDetailReq) (*types.FindPlatformEnterpriseDetailResp, error) {
	var resp *types.FindPlatformEnterpriseDetailResp

	grjResp, err := logicutil.DemoApiGrjResp(l.svcCtx, "/gurenju/wechat/platform/enterprise/findPlatformEnterpriseList?platformId=71D03B0337DD457E8CC7908C431A815D")
	if err != nil {
		return nil, err
	}

	dataBytes, err := json.Marshal(grjResp.Data)
	if err != nil {
		return nil, err
	}
	var items []*types.FindPlatformEnterpriseDetailResp
	err = json.Unmarshal(dataBytes, &items)
	if err != nil {
		return nil, err
	}

	if len(items) > 0 {
		return items[0], nil
	}

	return resp, nil
}
