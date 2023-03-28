package handler

import (
	"net/http"

	"github.com/xqk/cymzjs-api/api/cymzjs/internal/logic"
	"github.com/xqk/cymzjs-api/api/cymzjs/internal/svc"
	"github.com/xqk/cymzjs-api/api/cymzjs/internal/types"

	"git.zc0901.com/go/god/api/httpx"
)

func FindMemorialFeteHomePageHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.FindMemorialFeteHomePageReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		authorization := r.Header.Get("Authorization")
		l := logic.NewFindMemorialFeteHomePageLogic(r.Context(), ctx)
		resp, err := l.FindMemorialFeteHomePage(req, authorization)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
