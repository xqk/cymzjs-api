package logic

import (
	"context"
	"encoding/json"
	"github.com/xqk/cymzjs-api/api/cymzjs/internal/logicutil"

	"github.com/xqk/cymzjs-api/api/cymzjs/internal/svc"
	"github.com/xqk/cymzjs-api/api/cymzjs/internal/types"

	"git.zc0901.com/go/god/lib/logx"
)

type FindSystemFeteSceneCloudMournLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFindSystemFeteSceneCloudMournLogic(ctx context.Context, svcCtx *svc.ServiceContext) FindSystemFeteSceneCloudMournLogic {
	return FindSystemFeteSceneCloudMournLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FindSystemFeteSceneCloudMournLogic) FindSystemFeteSceneCloudMourn(req types.FindSystemFeteSceneReq) (*types.FindSystemFeteSceneResp, error) {
	var resp *types.FindSystemFeteSceneResp

	grjResp, err := logicutil.DemoApiGrjResp(l.svcCtx, "/gurenju/wechat/system/feteScene/findSystemFeteScene?objectId=8E77EAB2FD714912A43B3C42039B1258&objectType=cloud_mourn")
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
