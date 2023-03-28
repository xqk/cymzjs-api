package cymzjsmodel

import (
	"database/sql"
	"fmt"
	"sort"
	"strings"

	"git.zc0901.com/go/god/lib/container/garray"
	"git.zc0901.com/go/god/lib/fx"
	"git.zc0901.com/go/god/lib/g"
	"git.zc0901.com/go/god/lib/gconv"
	"git.zc0901.com/go/god/lib/gutil"
	"git.zc0901.com/go/god/lib/logx"
	"git.zc0901.com/go/god/lib/mathx"
	"git.zc0901.com/go/god/lib/mr"
	"git.zc0901.com/go/god/lib/store/cache"
	"git.zc0901.com/go/god/lib/store/sqlx"
	"git.zc0901.com/go/god/lib/stringx"
	"git.zc0901.com/go/god/tools/god/mysql/builder"
)

var (
	apiFieldList             = builder.FieldList(&Api{})
	apiFields                = strings.Join(apiFieldList, ",")
	apiFieldsAutoSet         = strings.Join(stringx.RemoveDBFields(apiFieldList, "id", "created_at", "updated_at", "create_time", "update_time"), ",")
	apiFieldsWithPlaceHolder = strings.Join(stringx.RemoveDBFields(apiFieldList, "id", "created_at", "updated_at", "create_time", "update_time"), "=?,") + "=?"

	cacheCymzjsApiPathPrefix = "cache:cymzjs:api:path:"
	cacheCymzjsApiIdPrefix   = "cache:cymzjs:api:id:"
)

type (
	Api struct {
		Id    int64  `db:"id" json:"id"`
		Path  string `db:"path" json:"path"`
		Value string `db:"value" json:"value"`
	}

	ApiModel struct {
		sqlx.CachedConn
		table string
	}
)

func NewApiModel(conn sqlx.Conn, clusterConf cache.ClusterConf) *ApiModel {
	return &ApiModel{
		CachedConn: sqlx.NewCachedConnWithCluster(conn, clusterConf),
		table:      "cymzjs.api",
	}
}

func (m *ApiModel) Insert(data Api) (sql.Result, error) {
	apiIdKey := fmt.Sprintf("%s%v", cacheCymzjsApiIdPrefix, data.Id)
	apiPathKey := fmt.Sprintf("%s%v", cacheCymzjsApiPathPrefix, data.Path)
	_ = m.DelCache(apiIdKey, apiPathKey)
	query := `insert into ` + m.table + ` (` + apiFieldsAutoSet + `) values (?, ?)`
	return m.ExecNoCache(query, data.Path, data.Value)
}

func (m *ApiModel) TxInsert(tx sqlx.TxSession, data Api) (sql.Result, error) {
	apiIdKey := fmt.Sprintf("%s%v", cacheCymzjsApiIdPrefix, data.Id)
	apiPathKey := fmt.Sprintf("%s%v", cacheCymzjsApiPathPrefix, data.Path)
	_ = m.DelCache(apiIdKey, apiPathKey)
	query := `insert into ` + m.table + ` (` + apiFieldsAutoSet + `) values (?, ?)`
	return tx.Exec(query, data.Path, data.Value)
}

func (m *ApiModel) FindOne(id int64) (*Api, error) {
	apiIdKey := fmt.Sprintf("%s%v", cacheCymzjsApiIdPrefix, id)
	var dest Api
	err := m.Query(&dest, apiIdKey, func(conn sqlx.Conn, v interface{}) error {
		query := `select ` + apiFields + ` from ` + m.table + ` where id = ? limit 1`
		return conn.Query(v, query, id)
	})
	if err == nil {
		return &dest, nil
	} else if err == sqlx.ErrNotFound {
		return nil, ErrNotFound
	} else {
		return nil, err
	}
}

func (m *ApiModel) FindMany(ids []int64, workers ...int) (list []*Api) {
	ids = gconv.Int64s(garray.NewArrayFrom(gconv.Interfaces(ids), true).Unique())

	var nWorkers int
	if len(workers) > 0 {
		nWorkers = workers[0]
	} else {
		nWorkers = mathx.MinInt(10, len(ids))
	}

	channel := mr.Map(func(source chan<- interface{}) {
		for _, id := range ids {
			source <- id
		}
	}, func(item interface{}, writer mr.Writer) {
		id := item.(int64)
		one, err := m.FindOne(id)
		if err == nil {
			writer.Write(one)
		} else {
			logx.Error(err)
		}
	}, mr.WithWorkers(nWorkers))

	for one := range channel {
		list = append(list, one.(*Api))
	}

	sort.Slice(list, func(i, j int) bool {
		return gutil.IndexOf(list[i].Id, ids) < gutil.IndexOf(list[j].Id, ids)
	})

	return
}

func (m *ApiModel) FindOneByPath(path string) (*Api, error) {
	apiPathKey := fmt.Sprintf("%s%v", cacheCymzjsApiPathPrefix, path)
	var dest Api
	err := m.QueryIndex(&dest, apiPathKey, func(primary interface{}) string {
		// 主键的缓存键
		return fmt.Sprintf("%s%v", cacheCymzjsApiIdPrefix, primary)
	}, func(conn sqlx.Conn, v interface{}) (i interface{}, e error) {
		// 无索引建——主键对应缓存，通过索引键查目标行
		query := `select ` + apiFields + ` from ` + m.table + ` where path = ? limit 1`
		if err := conn.Query(&dest, query, path); err != nil {
			return nil, err
		}
		return dest.Id, nil
	}, func(conn sqlx.Conn, v, primary interface{}) error {
		// 如果有索引建——主键对应缓存，则通过主键直接查目标航
		query := `select ` + apiFields + ` from ` + m.table + ` where id = ? limit 1`
		return conn.Query(v, query, primary)
	})
	if err == nil {
		return &dest, nil
	} else if err == sqlx.ErrNotFound {
		return nil, ErrNotFound
	} else {
		return nil, err
	}
}

func (m *ApiModel) FindManyByPaths(oKeys []string, workers ...int) (list []*Api) {
	keys := gconv.Strings(garray.NewArrayFrom(gconv.Interfaces(oKeys), true).Unique())

	var nWorkers int
	if len(workers) > 0 {
		nWorkers = workers[0]
	} else {
		nWorkers = mathx.MinInt(10, len(keys))
	}

	channel := mr.Map(func(source chan<- interface{}) {
		for _, key := range keys {
			source <- key
		}
	}, func(item interface{}, writer mr.Writer) {
		key := item.(string)
		one, err := m.FindOneByPath(key)
		if err == nil {
			writer.Write(one)
		} else {
			logx.Error(err)
		}
	}, mr.WithWorkers(nWorkers))

	for one := range channel {
		list = append(list, one.(*Api))
	}

	sort.Slice(list, func(i, j int) bool {
		return gutil.IndexOf(list[i].Path, keys) < gutil.IndexOf(list[j].Path, keys)
	})

	return
}

func (m *ApiModel) Update(data Api) error {
	apiIdKey := fmt.Sprintf("%s%v", cacheCymzjsApiIdPrefix, data.Id)
	_, err := m.Exec(func(conn sqlx.Conn) (result sql.Result, err error) {
		query := `update ` + m.table + ` set ` + apiFieldsWithPlaceHolder + ` where id = ?`
		return conn.Exec(query, data.Path, data.Value, data.Id)
	}, apiIdKey)
	return err
}

func (m *ApiModel) UpdatePartial(ms ...g.Map) (err error) {
	okNum := 0
	fx.From(func(source chan<- interface{}) {
		for _, data := range ms {
			source <- data
		}
	}).Parallel(func(item interface{}) {
		err = m.updatePartial(item.(g.Map))
		if err != nil {
			return
		}
		okNum++
	})

	if err == nil && okNum != len(ms) {
		err = fmt.Errorf("部分局部更新失败！待更新(%d) != 实际更新(%d)", len(ms), okNum)
	}

	return err
}

func (m *ApiModel) updatePartial(data g.Map) error {
	updateArgs, err := sqlx.ExtractUpdateArgs(apiFieldList, data)
	if err != nil {
		return err
	}

	apiIdKey := fmt.Sprintf("%s%v", cacheCymzjsApiIdPrefix, updateArgs.Id)
	_, err = m.Exec(func(conn sqlx.Conn) (result sql.Result, err error) {
		query := `update ` + m.table + ` set ` + updateArgs.Fields + ` where id = ` + updateArgs.Id
		return conn.Exec(query, updateArgs.Args...)
	}, apiIdKey)
	return err
}

func (m *ApiModel) TxUpdate(tx sqlx.TxSession, data Api) error {
	apiIdKey := fmt.Sprintf("%s%v", cacheCymzjsApiIdPrefix, data.Id)
	_, err := m.Exec(func(conn sqlx.Conn) (result sql.Result, err error) {
		query := `update ` + m.table + ` set ` + apiFieldsWithPlaceHolder + ` where id = ?`
		return tx.Exec(query, data.Path, data.Value, data.Id)
	}, apiIdKey)
	return err
}

func (m *ApiModel) TxUpdatePartial(tx sqlx.TxSession, ms ...g.Map) (err error) {
	okNum := 0
	fx.From(func(source chan<- interface{}) {
		for _, data := range ms {
			source <- data
		}
	}).Parallel(func(item interface{}) {
		err = m.txUpdatePartial(tx, item.(g.Map))
		if err != nil {
			return
		}
		okNum++
	})

	if err == nil && okNum != len(ms) {
		err = fmt.Errorf("部分事务型局部更新失败！待更新(%d) != 实际更新(%d)", len(ms), okNum)
	}
	return err
}

func (m *ApiModel) txUpdatePartial(tx sqlx.TxSession, data g.Map) error {
	updateArgs, err := sqlx.ExtractUpdateArgs(apiFieldList, data)
	if err != nil {
		return err
	}

	apiIdKey := fmt.Sprintf("%s%v", cacheCymzjsApiIdPrefix, updateArgs.Id)
	_, err = m.Exec(func(conn sqlx.Conn) (result sql.Result, err error) {
		query := `update ` + m.table + ` set ` + updateArgs.Fields + ` where id = ` + updateArgs.Id
		return tx.Exec(query, updateArgs.Args...)
	}, apiIdKey)
	return err
}

func (m *ApiModel) Delete(id ...int64) error {
	if len(id) == 0 {
		return nil
	}

	datas := m.FindMany(id)
	keys := make([]string, len(id)*2)
	for i, v := range id {
		data := datas[i]
		keys[i] = fmt.Sprintf("%s%v", cacheCymzjsApiIdPrefix, v)
		keys[i+1] = fmt.Sprintf("%s%v", cacheCymzjsApiPathPrefix, data.Path)
	}

	_, err := m.Exec(func(conn sqlx.Conn) (result sql.Result, err error) {
		query := fmt.Sprintf(`delete from `+m.table+` where id in (%s)`, sqlx.In(len(id)))
		return conn.Exec(query, gconv.Interfaces(id)...)
	}, keys...)
	return err
}

func (m *ApiModel) TxDelete(tx sqlx.TxSession, id ...int64) error {
	if len(id) == 0 {
		return nil
	}

	datas := m.FindMany(id)
	keys := make([]string, len(id)*2)
	for i, v := range id {
		data := datas[i]
		keys[i] = fmt.Sprintf("%s%v", cacheCymzjsApiIdPrefix, v)
		keys[i+1] = fmt.Sprintf("%s%v", cacheCymzjsApiPathPrefix, data.Path)
	}

	_, err := m.Exec(func(conn sqlx.Conn) (result sql.Result, err error) {
		query := fmt.Sprintf(`delete from `+m.table+` where id in (%s)`, sqlx.In(len(id)))
		return tx.Exec(query, gconv.Interfaces(id)...)
	}, keys...)
	return err
}
