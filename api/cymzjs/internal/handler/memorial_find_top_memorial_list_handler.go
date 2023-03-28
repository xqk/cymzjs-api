package handler

import (
	"net/http"

	"git.zc0901.com/go/god/api/httpx"
	"github.com/xqk/cymzjs-api/api/cymzjs/internal/logic"
	"github.com/xqk/cymzjs-api/api/cymzjs/internal/svc"
)

func MemorialFindTopMemorialListHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		l := logic.NewMemorialFindTopMemorialListLogic(r.Context(), ctx)
		resp, err := l.MemorialFindTopMemorialList()
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
