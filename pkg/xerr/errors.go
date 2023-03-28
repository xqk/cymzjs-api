package xerr

import (
	"fmt"

	"git.zc0901.com/go/god/lib/logx"
	"github.com/pkg/errors"
	"google.golang.org/grpc/status"
)

// CodeErr 通用的错误信息结构体
type CodeErr struct {
	errCode int64
	errMsg  string
}

const (
	ErrCodeCommon        int64 = 100001
	ErrCodeDB            int64 = 100002
	ErrCodeRequestParams int64 = 100003
	ErrCodeAccessDenied  int64 = 100004
	ErrCodeData          int64 = 100005
	ErrCodeNoLogin       int64 = 100006
	ErrCodeUserNotFound  int64 = 100007
	ErrCodeMPLoginFailed int64 = 100010
	ErrCodeMPInterFailed int64 = 100011
	ErrCodeNotFound      int64 = 100012
)

var (
	messages = map[int64]string{
		ErrCodeCommon:        "服务器开小差啦，请稍后再试",
		ErrCodeDB:            "数据库繁忙，请稍后再试",
		ErrCodeRequestParams: "请求参数有误，请稍后再试",
		ErrCodeAccessDenied:  "没有操作权限，拒绝访问",
		ErrCodeData:          "数据出现异常，请稍后再试",
		ErrCodeNoLogin:       "当前还未登录，请登录后再试",
		ErrCodeUserNotFound:  "用户不存在",
		ErrCodeMPLoginFailed: "微信小程序登录失败，请稍后再试",
		ErrCodeMPInterFailed: "请求微信数据失败，请稍后再试",
		ErrCodeNotFound:      "数据走丢啦，请稍后再试",
	}
)

func (err *CodeErr) Error() string {
	if err.errMsg == "" {
		return ErrMsgByCode(err.errCode)
	}
	return err.errMsg
}

func (err *CodeErr) GetErrCode() int64 {
	return err.errCode
}

func (err *CodeErr) GetErrMsg() string {
	return err.errMsg
}

func NewCodeErr(errCode int64, errMsg string) error {
	return &CodeErr{errCode: errCode, errMsg: errMsg}
}

func NewCodeErrFromCode(errCode int64) error {
	return &CodeErr{errCode: errCode, errMsg: ErrMsgByCode(errCode)}
}

func NewCodeErrFmt(errCode int64, errMsg string, args ...interface{}) error {
	errMsg = ErrMsgByCode(errCode) + "；" + errMsg
	return &CodeErr{errCode: errCode, errMsg: fmt.Sprintf(errMsg, args...)}
}

func NewCodeErrFromRpcErr(err error) error {
	if err == nil {
		return nil
	}

	errCode := ErrCodeCommon
	errMsg := ErrMsgByCode(errCode)

	causeErr := errors.Cause(err)
	if e, ok := causeErr.(*CodeErr); ok {
		errCode = e.GetErrCode()
		errMsg = e.GetErrMsg()
	} else {
		if st, ok := status.FromError(causeErr); ok {
			errCode = int64(st.Code())
			errMsg = st.Message()
			if VerifyErrCode(errCode) {
				errMsg = ErrMsgByCode(errCode)
			}
		}
	}

	return NewCodeErr(errCode, errMsg)
}

func NewCodeErrPrintf(errCode int64, errMsg string, args ...interface{}) error {
	err := &CodeErr{errCode: errCode, errMsg: ErrMsgByCode(errCode)}
	logx.Errorf(errMsg, args...)
	return err
}

// ErrMsgByCode 根据错误码返回
func ErrMsgByCode(errCode int64) string {
	if errMsg, ok := messages[errCode]; ok {
		return errMsg
	}
	return messages[ErrCodeCommon]
}

// VerifyErrCode 验证errCode是否已定义
func VerifyErrCode(errCode int64) bool {
	if _, ok := messages[errCode]; ok {
		return true
	}
	return false
}
