package svc

import (
	"git.zc0901.com/go/god/lib/store/kv"
	"git.zc0901.com/go/god/lib/store/sqlx"
	"github.com/xqk/cymzjs-api/api/cymzjs/internal/config"
	cymzjsmodel "github.com/xqk/cymzjs-api/dao/cymzjs"
)

type ServiceContext struct {
	Config config.Config
	Store  kv.Store

	UserModel     *cymzjsmodel.UserModel
	MemorialModel *cymzjsmodel.MemorialModel
	ApiModel      *cymzjsmodel.ApiModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMySQL(c.Mysql.Cymzjs)
	return &ServiceContext{
		Config: c,
		Store:  kv.NewStore(c.Cache),

		UserModel:     cymzjsmodel.NewUserModel(conn, c.Cache),
		MemorialModel: cymzjsmodel.NewMemorialModel(conn, c.Cache),
		ApiModel:      cymzjsmodel.NewApiModel(conn, c.Cache),
	}
}
