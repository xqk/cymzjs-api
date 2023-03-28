package logic

import (
	"context"
	"encoding/json"
	"github.com/xqk/cymzjs-api/api/cymzjs/internal/logicutil"
	"github.com/xqk/cymzjs-api/api/cymzjs/internal/svc"
	"github.com/xqk/cymzjs-api/api/cymzjs/internal/types"

	"git.zc0901.com/go/god/lib/logx"
)

type BannerFindBannerListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewBannerFindBannerListLogic(ctx context.Context, svcCtx *svc.ServiceContext) BannerFindBannerListLogic {
	return BannerFindBannerListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BannerFindBannerListLogic) BannerFindBannerList(req types.BannerFindBannerListReq) ([]*types.BannerFindBannerListItem, error) {
	items := make([]*types.BannerFindBannerListItem, 0)

	grjResp, err := logicutil.DemoApiGrjResp(l.svcCtx, "/gurenju/wechat/banner/findBannerList?type=0")
	if err != nil {
		return nil, err
	}

	dataBytes, err := json.Marshal(grjResp.Data)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(dataBytes, &items)
	if err != nil {
		return nil, err
	}

	return items, nil
}
