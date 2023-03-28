package handler

import (
	"net/http"

	"github.com/xqk/cymzjs-api/api/cymzjs/internal/logic"
	"github.com/xqk/cymzjs-api/api/cymzjs/internal/svc"
	"github.com/xqk/cymzjs-api/api/cymzjs/internal/types"

	"git.zc0901.com/go/god/api/httpx"
)

func FindPlatformEnterpriseDetailHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.FindPlatformEnterpriseDetailReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := logic.NewFindPlatformEnterpriseDetailLogic(r.Context(), ctx)
		resp, err := l.FindPlatformEnterpriseDetail(req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
