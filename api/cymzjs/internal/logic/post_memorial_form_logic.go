package logic

import (
	"context"
	"git.zc0901.com/go/god/lib/g"
	"git.zc0901.com/go/god/lib/gconv"
	cymzjsmodel "github.com/xqk/cymzjs-api/dao/cymzjs"
	"github.com/xqk/cymzjs-api/pkg"

	"github.com/xqk/cymzjs-api/api/cymzjs/internal/svc"
	"github.com/xqk/cymzjs-api/api/cymzjs/internal/types"

	"git.zc0901.com/go/god/lib/logx"
)

type PostMemorialFormLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPostMemorialFormLogic(ctx context.Context, svcCtx *svc.ServiceContext) PostMemorialFormLogic {
	return PostMemorialFormLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PostMemorialFormLogic) PostMemorialForm(req types.MemorialForm, authorization string) (*types.PostMemorialFormResp, error) {
	secret := l.svcCtx.Config.Auth.AccessSecret
	claims, err := pkg.JwtDecode(authorization, secret)
	var userId int64
	if err == nil {
		userId = gconv.Int64(claims["userId"])
	}
	if userId == 0 {

	}

	var ids []int64
	err = l.svcCtx.MemorialModel.QueryNoCache(&ids, sqlUserFirstMemorialId, userId)
	if err != nil {
		return nil, err
	}

	var memorialId int64
	if len(ids) == 0 {
		insertResult, err := l.svcCtx.MemorialModel.Insert(cymzjsmodel.Memorial{
			UserId:      userId,
			HeadImg:     req.HeadImg,
			NameOne:     req.NameOne,
			RelationOne: req.RelationOne,
			NameTwo:     req.NameTwo,
			RelationTwo: req.RelationTwo,
		})
		if err != nil {
			return nil, err
		}

		memorialId, err = insertResult.LastInsertId()
		if err != nil {
			return nil, err
		}
	} else {
		err := l.svcCtx.MemorialModel.UpdatePartial(g.Map{
			"id":           ids[0],
			"head_img":     req.HeadImg,
			"name_one":     req.NameOne,
			"relation_one": req.RelationOne,
			"name_two":     req.NameTwo,
			"relation_two": req.RelationTwo,
		})
		if err != nil {
			return nil, err
		}

		memorialId = ids[0]
	}

	return &types.PostMemorialFormResp{
		Id: memorialId,
	}, nil
}
