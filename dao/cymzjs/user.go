package cymzjsmodel

import (
	"database/sql"
	"fmt"
	"sort"
	"strings"
	"time"

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
	userFieldList             = builder.FieldList(&User{})
	userFields                = strings.Join(userFieldList, ",")
	userFieldsAutoSet         = strings.Join(stringx.RemoveDBFields(userFieldList, "id", "created_at", "updated_at", "create_time", "update_time"), ",")
	userFieldsWithPlaceHolder = strings.Join(stringx.RemoveDBFields(userFieldList, "id", "created_at", "updated_at", "create_time", "update_time"), "=?,") + "=?"

	cacheCymzjsUserIdPrefix     = "cache:cymzjs:user:id:"
	cacheCymzjsUserOpenidPrefix = "cache:cymzjs:user:openid:"
)

type (
	User struct {
		Id         int64     `db:"id" json:"id"`
		Openid     string    `db:"openid" json:"openid"`
		CreateTime time.Time `db:"create_time" json:"createTime"` // 创建时间
		UpdateTime time.Time `db:"update_time" json:"updateTime"` // 更新时间
	}

	UserModel struct {
		sqlx.CachedConn
		table string
	}
)

func NewUserModel(conn sqlx.Conn, clusterConf cache.ClusterConf) *UserModel {
	return &UserModel{
		CachedConn: sqlx.NewCachedConnWithCluster(conn, clusterConf),
		table:      "cymzjs.user",
	}
}

func (m *UserModel) Insert(data User) (sql.Result, error) {
	userOpenidKey := fmt.Sprintf("%s%v", cacheCymzjsUserOpenidPrefix, data.Openid)
	userIdKey := fmt.Sprintf("%s%v", cacheCymzjsUserIdPrefix, data.Id)
	_ = m.DelCache(userOpenidKey, userIdKey)
	query := `insert into ` + m.table + ` (` + userFieldsAutoSet + `) values (?)`
	return m.ExecNoCache(query, data.Openid)
}

func (m *UserModel) TxInsert(tx sqlx.TxSession, data User) (sql.Result, error) {
	userIdKey := fmt.Sprintf("%s%v", cacheCymzjsUserIdPrefix, data.Id)
	userOpenidKey := fmt.Sprintf("%s%v", cacheCymzjsUserOpenidPrefix, data.Openid)
	_ = m.DelCache(userIdKey, userOpenidKey)
	query := `insert into ` + m.table + ` (` + userFieldsAutoSet + `) values (?)`
	return tx.Exec(query, data.Openid)
}

func (m *UserModel) FindOne(id int64) (*User, error) {
	userIdKey := fmt.Sprintf("%s%v", cacheCymzjsUserIdPrefix, id)
	var dest User
	err := m.Query(&dest, userIdKey, func(conn sqlx.Conn, v interface{}) error {
		query := `select ` + userFields + ` from ` + m.table + ` where id = ? limit 1`
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

func (m *UserModel) FindMany(ids []int64, workers ...int) (list []*User) {
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
		list = append(list, one.(*User))
	}

	sort.Slice(list, func(i, j int) bool {
		return gutil.IndexOf(list[i].Id, ids) < gutil.IndexOf(list[j].Id, ids)
	})

	return
}

func (m *UserModel) FindOneByOpenid(openid string) (*User, error) {
	userOpenidKey := fmt.Sprintf("%s%v", cacheCymzjsUserOpenidPrefix, openid)
	var dest User
	err := m.QueryIndex(&dest, userOpenidKey, func(primary interface{}) string {
		// 主键的缓存键
		return fmt.Sprintf("%s%v", cacheCymzjsUserIdPrefix, primary)
	}, func(conn sqlx.Conn, v interface{}) (i interface{}, e error) {
		// 无索引建——主键对应缓存，通过索引键查目标行
		query := `select ` + userFields + ` from ` + m.table + ` where openid = ? limit 1`
		if err := conn.Query(&dest, query, openid); err != nil {
			return nil, err
		}
		return dest.Id, nil
	}, func(conn sqlx.Conn, v, primary interface{}) error {
		// 如果有索引建——主键对应缓存，则通过主键直接查目标航
		query := `select ` + userFields + ` from ` + m.table + ` where id = ? limit 1`
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

func (m *UserModel) FindManyByOpenids(oKeys []string, workers ...int) (list []*User) {
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
		one, err := m.FindOneByOpenid(key)
		if err == nil {
			writer.Write(one)
		} else {
			logx.Error(err)
		}
	}, mr.WithWorkers(nWorkers))

	for one := range channel {
		list = append(list, one.(*User))
	}

	sort.Slice(list, func(i, j int) bool {
		return gutil.IndexOf(list[i].Openid, keys) < gutil.IndexOf(list[j].Openid, keys)
	})

	return
}

func (m *UserModel) Update(data User) error {
	userIdKey := fmt.Sprintf("%s%v", cacheCymzjsUserIdPrefix, data.Id)
	_, err := m.Exec(func(conn sqlx.Conn) (result sql.Result, err error) {
		query := `update ` + m.table + ` set ` + userFieldsWithPlaceHolder + ` where id = ?`
		return conn.Exec(query, data.Openid, data.Id)
	}, userIdKey)
	return err
}

func (m *UserModel) UpdatePartial(ms ...g.Map) (err error) {
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

func (m *UserModel) updatePartial(data g.Map) error {
	updateArgs, err := sqlx.ExtractUpdateArgs(userFieldList, data)
	if err != nil {
		return err
	}

	userIdKey := fmt.Sprintf("%s%v", cacheCymzjsUserIdPrefix, updateArgs.Id)
	_, err = m.Exec(func(conn sqlx.Conn) (result sql.Result, err error) {
		query := `update ` + m.table + ` set ` + updateArgs.Fields + ` where id = ` + updateArgs.Id
		return conn.Exec(query, updateArgs.Args...)
	}, userIdKey)
	return err
}

func (m *UserModel) TxUpdate(tx sqlx.TxSession, data User) error {
	userIdKey := fmt.Sprintf("%s%v", cacheCymzjsUserIdPrefix, data.Id)
	_, err := m.Exec(func(conn sqlx.Conn) (result sql.Result, err error) {
		query := `update ` + m.table + ` set ` + userFieldsWithPlaceHolder + ` where id = ?`
		return tx.Exec(query, data.Openid, data.Id)
	}, userIdKey)
	return err
}

func (m *UserModel) TxUpdatePartial(tx sqlx.TxSession, ms ...g.Map) (err error) {
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

func (m *UserModel) txUpdatePartial(tx sqlx.TxSession, data g.Map) error {
	updateArgs, err := sqlx.ExtractUpdateArgs(userFieldList, data)
	if err != nil {
		return err
	}

	userIdKey := fmt.Sprintf("%s%v", cacheCymzjsUserIdPrefix, updateArgs.Id)
	_, err = m.Exec(func(conn sqlx.Conn) (result sql.Result, err error) {
		query := `update ` + m.table + ` set ` + updateArgs.Fields + ` where id = ` + updateArgs.Id
		return tx.Exec(query, updateArgs.Args...)
	}, userIdKey)
	return err
}

func (m *UserModel) Delete(id ...int64) error {
	if len(id) == 0 {
		return nil
	}

	datas := m.FindMany(id)
	keys := make([]string, len(id)*2)
	for i, v := range id {
		data := datas[i]
		keys[i] = fmt.Sprintf("%s%v", cacheCymzjsUserIdPrefix, v)
		keys[i+1] = fmt.Sprintf("%s%v", cacheCymzjsUserOpenidPrefix, data.Openid)
	}

	_, err := m.Exec(func(conn sqlx.Conn) (result sql.Result, err error) {
		query := fmt.Sprintf(`delete from `+m.table+` where id in (%s)`, sqlx.In(len(id)))
		return conn.Exec(query, gconv.Interfaces(id)...)
	}, keys...)
	return err
}

func (m *UserModel) TxDelete(tx sqlx.TxSession, id ...int64) error {
	if len(id) == 0 {
		return nil
	}

	datas := m.FindMany(id)
	keys := make([]string, len(id)*2)
	for i, v := range id {
		data := datas[i]
		keys[i] = fmt.Sprintf("%s%v", cacheCymzjsUserIdPrefix, v)
		keys[i+1] = fmt.Sprintf("%s%v", cacheCymzjsUserOpenidPrefix, data.Openid)
	}

	_, err := m.Exec(func(conn sqlx.Conn) (result sql.Result, err error) {
		query := fmt.Sprintf(`delete from `+m.table+` where id in (%s)`, sqlx.In(len(id)))
		return tx.Exec(query, gconv.Interfaces(id)...)
	}, keys...)
	return err
}
