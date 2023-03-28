// Code generated by god. DO NOT EDIT.
package handler

import (
	"net/http"

	"github.com/xqk/cymzjs-api/api/cymzjs/internal/svc"

	"git.zc0901.com/go/god/api"
)

func RegisterHandlers(engine *api.Server, serverCtx *svc.ServiceContext) {
	engine.AddRoutes(
		[]api.Route{
			{
				Method:  http.MethodGet,
				Path:    "/api/qiniu/token",
				Handler: QiniuTokenHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/api/login",
				Handler: LoginHandler(serverCtx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/api/log",
				Handler: LogHandler(serverCtx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/api/banner/findBannerList",
				Handler: BannerFindBannerListHandler(serverCtx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/api/memorial/findTopMemorialList",
				Handler: MemorialFindTopMemorialListHandler(serverCtx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/api/system/feteScene/findSystemFeteScenePlatform",
				Handler: FindSystemFeteScenePlatformHandler(serverCtx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/api/system/feteScene/findSystemFeteSceneCloudMourn",
				Handler: FindSystemFeteSceneCloudMournHandler(serverCtx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/api/platform/platform/findPlatformMemorialConfig",
				Handler: FindPlatformMemorialConfigHandler(serverCtx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/api/platform/platform/findPlatformMemorialConfigDetail",
				Handler: FindPlatformMemorialConfigDetailHandler(serverCtx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/api/platform/platform/findPlatformDetails",
				Handler: FindPlatformDetailsHandler(serverCtx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/api/platform/enterprise/findPlatformEnterpriseDetail",
				Handler: FindPlatformEnterpriseDetailHandler(serverCtx),
			},
		},
	)

	engine.AddRoutes(
		[]api.Route{
			{
				Method:  http.MethodGet,
				Path:    "/api/memorial/get/form",
				Handler: GetMemorialFormHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/api/memorial/post/form",
				Handler: PostMemorialFormHandler(serverCtx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/api/memorial/findMemorialFeteHomePage",
				Handler: FindMemorialFeteHomePageHandler(serverCtx),
			},
		},
		api.WithJwt(serverCtx.Config.Auth.AccessSecret),
	)
}