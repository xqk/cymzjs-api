package logic

import (
	"context"
	"github.com/xqk/cymzjs-api/api/cymzjs/internal/types"

	"git.zc0901.com/go/god/lib/logx"
	"github.com/xqk/cymzjs-api/api/cymzjs/internal/svc"
)

type MemorialFindTopMemorialListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMemorialFindTopMemorialListLogic(ctx context.Context, svcCtx *svc.ServiceContext) MemorialFindTopMemorialListLogic {
	return MemorialFindTopMemorialListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MemorialFindTopMemorialListLogic) MemorialFindTopMemorialList() ([]*types.MemorialFindTopMemorialListItem, error) {
	items := make([]*types.MemorialFindTopMemorialListItem, 0)

	//grjResp, err := logicutil.DemoApiGrjResp(l.svcCtx, "/gurenju/wechat/memorial/findTopMemorialListV3")
	//if err != nil {
	//	return nil, err
	//}
	//
	//dataBytes, err := json.Marshal(grjResp.Data)
	//if err != nil {
	//	return nil, err
	//}
	//
	//var groupItems []*struct {
	//	Name     string      `json:"name"`
	//	Type     string      `json:"type"`
	//	DataList interface{} `json:"dataList"`
	//}
	//err = json.Unmarshal(dataBytes, &groupItems)
	//if err != nil {
	//	return nil, err
	//}
	//
	//for _, v := range groupItems {
	//	if v.Type == "6" {
	//		dataListBytes, err := json.Marshal(v.DataList)
	//		if err != nil {
	//			return nil, err
	//		}
	//		err = json.Unmarshal(dataListBytes, &items)
	//		if err != nil {
	//			return nil, err
	//		}
	//	}
	//}

	return items, nil
}
