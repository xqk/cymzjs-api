package logic

import (
	"context"
	"fmt"
	cymzjsmodel "github.com/xqk/cymzjs-api/dao/cymzjs"
	"github.com/xqk/cymzjs-api/pkg"
	"github.com/xqk/cymzjs-api/pkg/xerr"

	"github.com/xqk/cymzjs-api/api/cymzjs/internal/svc"
	"github.com/xqk/cymzjs-api/api/cymzjs/internal/types"

	"git.zc0901.com/go/god/lib/logx"
	"github.com/medivhzhan/weapp/v2"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) LoginLogic {
	return LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req types.LoginReq) (*types.LoginResp, error) {
	// 微信登录，code换openId
	openId, _, err := l.code2session(l.svcCtx.Config.AppId, l.svcCtx.Config.AppSecret, req.Code)
	if err != nil {
		return nil, err
	}

	user, err := l.svcCtx.UserModel.FindOneByOpenid(openId)
	if err != nil && err != cymzjsmodel.ErrNotFound {
		return nil, err
	}
	var userId int64
	if err != nil && err == cymzjsmodel.ErrNotFound {
		insertResult, err := l.svcCtx.UserModel.Insert(cymzjsmodel.User{
			Openid: openId,
		})
		if err != nil {
			return nil, err
		}

		// 删除缓存
		_, err = l.svcCtx.Store.Del(fmt.Sprintf("cache:cymzjs:user:openid:%s", openId))
		if err != nil {
			return nil, err
		}

		userId, err = insertResult.LastInsertId()
		if err != nil {
			return nil, err
		}
	} else {
		userId = user.Id
	}

	// 生成jwt身份令牌
	token, err := pkg.JwtFromUid(l.svcCtx.Config.Auth.AccessSecret, l.svcCtx.Config.Auth.AccessExpire, userId)
	if err != nil {
		return nil, xerr.NewCodeErrFmt(xerr.ErrCodeCommon, "生成身份令牌失败，error:%v", err)
	}

	return &types.LoginResp{
		UserId: userId,
		Token:  fmt.Sprintf("Bearer %s", token),
	}, nil
}

// code2session 微信登录，code换openId
func (l *LoginLogic) code2session(appId, appSecret, code string) (string, string, error) {
	resp, err := weapp.Login(appId, appSecret, code)
	if err != nil {
		return "", "", xerr.NewCodeErrFmt(xerr.ErrCodeMPLoginFailed, "error:%v", err)
	}
	if err := resp.GetResponseError(); err != nil {
		return "", "", xerr.NewCodeErrFmt(xerr.ErrCodeMPLoginFailed, "error:%v", err)
	}
	return resp.OpenID, resp.UnionID, nil
}
