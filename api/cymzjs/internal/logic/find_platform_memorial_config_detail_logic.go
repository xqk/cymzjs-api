package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"git.zc0901.com/go/god/lib/gconv"
	"github.com/xqk/cymzjs-api/api/cymzjs/internal/logicutil"
	"github.com/xqk/cymzjs-api/api/cymzjs/internal/svc"
	"github.com/xqk/cymzjs-api/api/cymzjs/internal/types"
	"github.com/xqk/cymzjs-api/pkg"

	"git.zc0901.com/go/god/lib/logx"
)

type FindPlatformMemorialConfigDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFindPlatformMemorialConfigDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) FindPlatformMemorialConfigDetailLogic {
	return FindPlatformMemorialConfigDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FindPlatformMemorialConfigDetailLogic) FindPlatformMemorialConfigDetail(req types.FindPlatformMemorialConfigDetailReq, authorization string) (*types.FindPlatformMemorialConfigItem, error) {
	if req.Id == "undefined" {
		req.Id = ""
	}
	resp := &types.FindPlatformMemorialConfigItem{}
	if req.Id == "" {
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
			return resp, nil
		}

		one, err := l.svcCtx.MemorialModel.FindOne(ids[0])
		if err != nil {
			return nil, err
		}
		resp = &types.FindPlatformMemorialConfigItem{
			Id:              fmt.Sprintf("%d", one.Id),
			Name:            one.NameOne,
			OptId:           "",
			OptName:         "",
			IsAble:          0,
			HeadImage:       one.HeadImg,
			CreateTime:      0,
			PlatformId:      "",
			UpdateTime:      0,
			BackgroundImage: "image/2022/03/23/916C1E6E362E42BCA4847E1CE759EC70.jpg",
			NameOne:         one.NameOne,
			RelationOne:     one.RelationOne,
			NameTwo:         one.NameTwo,
			RelationTwo:     one.RelationTwo,
		}
	} else {
		var items []*types.FindPlatformMemorialConfigItem
		grjResp, err := logicutil.DemoApiGrjResp(l.svcCtx, "/gurenju/wechat/platform/platform/findPlatformMemorialConfig?platformId=71D03B0337DD457E8CC7908C431A815D")
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

		for _, v := range items {
			if v.Id == req.Id {
				resp = v
			}
		}
	}

	return resp, nil
}
