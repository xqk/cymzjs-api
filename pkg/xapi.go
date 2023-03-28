package pkg

import (
	"git.zc0901.com/go/god/lib/logx"
	"github.com/xqk/cymzjs-api/pkg/xerr"
	"net/http"
)

type Response struct {
	Code    int64       `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

func ApiErrHandler(err error) (int, interface{}) {
	if e, ok := err.(*xerr.CodeErr); ok {
		return http.StatusOK, Response{Code: e.GetErrCode(), Message: e.GetErrMsg()}
	}
	return http.StatusOK, Response{Code: -1, Message: err.Error()}
}

func ApiOKHandler(data interface{}) interface{} {
	return Response{Code: 0, Data: data}
}

// ApiErrorLog 打印名为 name 的 API 产生的错误。
func ApiErrorLog(name, msg string, v interface{}) {
	if v == nil {
		logx.Errorf("[%v]: %s", name, msg)
	} else {
		logx.Errorf("[%v]: %s\n%v", name, msg, v)
	}
}
