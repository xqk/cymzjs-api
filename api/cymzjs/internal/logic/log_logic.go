package logic

import (
	"context"
	"github.com/xqk/cymzjs-api/api/cymzjs/internal/logicutil"

	"github.com/xqk/cymzjs-api/api/cymzjs/internal/svc"
	"github.com/xqk/cymzjs-api/api/cymzjs/internal/types"

	"git.zc0901.com/go/god/lib/logx"
)

type LogLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLogLogic(ctx context.Context, svcCtx *svc.ServiceContext) LogLogic {
	return LogLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LogLogic) Log(req types.LogReq) error {
	switch req.Id {
	case logicutil.LogPos11001.Int64():
		_, err := logicutil.CacheSetVisitCount(l.svcCtx.Store)
		if err != nil {
			logx.Error(err)
			return nil
		}
	}

	return nil
}
