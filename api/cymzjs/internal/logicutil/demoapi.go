package logicutil

import (
	"encoding/json"
	"git.zc0901.com/go/god/lib/gconv"
	"github.com/xqk/cymzjs-api/api/cymzjs/internal/svc"
	"github.com/xqk/cymzjs-api/api/cymzjs/internal/types"
)

func DemoApiGrjResp(svcCtx *svc.ServiceContext, path string) (*types.GrjResp, error) {
	// 模拟数据
	demo, err := svcCtx.ApiModel.FindOneByPath(path)
	if err != nil {
		return nil, err
	}

	var grjResp *types.GrjResp
	err = json.Unmarshal(gconv.Bytes(demo.Value), &grjResp)
	if err != nil {
		return nil, err
	}

	return grjResp, nil
}
