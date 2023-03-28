package handler

import (
	"net/http"

	"git.zc0901.com/go/god/api/httpx"
	"github.com/xqk/cymzjs-api/api/cymzjs/internal/logic"
	"github.com/xqk/cymzjs-api/api/cymzjs/internal/svc"
)

func FindPlatformMemorialConfigHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		l := logic.NewFindPlatformMemorialConfigLogic(r.Context(), ctx)
		resp, err := l.FindPlatformMemorialConfig()
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
