package logic

import (
	"context"
	"git.zc0901.com/go/god/lib/gconv"
	"github.com/xqk/cymzjs-api/pkg"

	"github.com/xqk/cymzjs-api/api/cymzjs/internal/svc"
	"github.com/xqk/cymzjs-api/api/cymzjs/internal/types"

	"git.zc0901.com/go/god/lib/logx"
)

type FindMemorialFeteHomePageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFindMemorialFeteHomePageLogic(ctx context.Context, svcCtx *svc.ServiceContext) FindMemorialFeteHomePageLogic {
	return FindMemorialFeteHomePageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FindMemorialFeteHomePageLogic) FindMemorialFeteHomePage(req types.FindMemorialFeteHomePageReq, authorization string) (*types.FindMemorialFeteHomePageResp, error) {
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
		return &types.FindMemorialFeteHomePageResp{}, nil
	}

	one, err := l.svcCtx.MemorialModel.FindOne(ids[0])
	if err != nil {
		return nil, err
	}

	res := &types.FindMemorialFeteHomePageResp{
		Hudie: []string{
			"resource/wltang/hudie/left.gif",
			"resource/wltang/hudie/right.gif",
		},
		Pinyao: []string{
			"image/2021/08/13/5230673852044D1DA76259074BAA5E64.png",
			"image/2021/08/13/6B1D0C2E120E4F0CBDC907ADE336654F.png",
			"image/2021/08/13/2B227C4BE8704EA5AFC78BF87251B485.png",
			"image/2021/08/13/C7C25C2EB2724A9A9AE3B6BA6E144ACA.png",
		},
		Peiji: []string{
			"resource/wltang/peiji/left.gif",
			"resource/wltang/peiji/right.gif",
		},
		Xianheshengui: []string{
			"image/2021/08/13/0DE72E756AE34C60823B02B782EB7D8E.png",
			"image/2021/08/13/1AC581B08EA34C8BB535D8D94555FF4F.png",
		},
		Saomu: []string{},
		Xianhuabsdz: []string{
			"resource/wltang/xianhuabsdz/nanshiguibai.gif",
			"resource/wltang/xianhuabsdz/nvshiguibai.gif",
		},
		Lazhu: []string{
			"image/2021/08/13/A00AA1A2F7B34592A23BC57FE959D029.gif",
			"image/2021/08/13/F72D9EFE3BA84B1C9B04E9B6002957BA.gif",
		},
		Chajiu: []string{
			"image/2021/08/13/BE77BB1F2C73450BA51631DBA65E3DFC.png",
		},
		Jibai: []string{},
		Bianpao: []string{
			"resource/wltang/bianpao/left.gif",
			"resource/wltang/bianpao/right.gif",
		},
		Guangyun: []string{
			"resource/wltang/guangyun/guangyun_01.gif",
		},
		Xianglu: []string{
			"image/2021/08/13/ACA666997A3D44E0AD6EE86E0952F813.gif",
		},
		Saomubsdz: []string{
			"resource/wltang/saomubsdz/left.gif",
			"resource/wltang/saomubsdz/right.gif",
		},
		BianpaoMusic: []string{
			"resource/wltang/bianpao/backgroundMusic.mp3",
		},
		Xianhua:             []string{},
		UserFeteList:        []string{},
		BackgroundMusicName: "天上人间",
		BackgroundMusicUrl:  "music/2020/08/09/A5BDCECA9D9D421E9D9E5F8A67A69864.mp3",
		MemorialInfo: &types.MemorialInfo{
			Id:                     userId,
			MemorialNo:             "",
			HeadImg:                one.HeadImg,
			DefaultHeadImg:         "image/2020/05/13/B7A381B0370547FA88154D14095245EA.png",
			Rahmen:                 "resource/default_rahmen.gif",
			Homepage:               "resource/default_homepage.jpg",
			BackgroundMusic:        "music/2019/03/27/16448345CC924ECFB82DEFF2B31983F9.mp3",
			JntBackgroundMusic:     "",
			WltBackgroundMusic:     "",
			QftBackgroundMusic:     "",
			BackgroundMusicName:    "清明雨上",
			JntBackgroundMusicName: "",
			WltBackgroundMusicName: "",
			QftBackgroundMusicName: "",
			Name:                   one.NameOne,
			Relation:               "",
			Birthdate:              "",
			Deathdate:              "",
			Remark:                 "",
			Blog:                   "",
			JntBackgroundImage:     "resource/jntang/background.jpg",
			WltBackgroundImage:     "resource/wltang/background.jpg",
			QftBackgroundImage:     "resource/cifutang/background.jpg",
		},
	}

	return res, nil
}
