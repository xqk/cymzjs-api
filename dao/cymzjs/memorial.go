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
	memorialFieldList             = builder.FieldList(&Memorial{})
	memorialFields                = strings.Join(memorialFieldList, ",")
	memorialFieldsAutoSet         = strings.Join(stringx.RemoveDBFields(memorialFieldList, "id", "created_at", "updated_at", "create_time", "update_time"), ",")
	memorialFieldsWithPlaceHolder = strings.Join(stringx.RemoveDBFields(memorialFieldList, "id", "created_at", "updated_at", "create_time", "update_time"), "=?,") + "=?"

	cacheCymzjsMemorialIdPrefix = "cache:cymzjs:memorial:id:"
)

type (
	Memorial struct {
		Id          int64     `db:"id" json:"id"`
		UserId      int64     `db:"user_id" json:"userId"`
		HeadImg     string    `db:"head_img" json:"headImg"`
		NameOne     string    `db:"name_one" json:"nameOne"`
		RelationOne string    `db:"relation_one" json:"relationOne"`
		NameTwo     string    `db:"name_two" json:"nameTwo"`
		RelationTwo string    `db:"relation_two" json:"relationTwo"`
		CreateTime  time.Time `db:"create_time" json:"createTime"` // 创建时间
		UpdateTime  time.Time `db:"update_time" json:"updateTime"` // 更新时间
	}

	MemorialModel struct {
		sqlx.CachedConn
		table string
	}
)

func NewMemorialModel(conn sqlx.Conn, clusterConf cache.ClusterConf) *MemorialModel {
	return &MemorialModel{
		CachedConn: sqlx.NewCachedConnWithCluster(conn, clusterConf),
		table:      "cymzjs.memorial",
	}
}

func (m *MemorialModel) Insert(data Memorial) (sql.Result, error) {
	memorialIdKey := fmt.Sprintf("%s%v", cacheCymzjsMemorialIdPrefix, data.Id)
	_ = m.DelCache(memorialIdKey)
	query := `insert into ` + m.table + ` (` + memorialFieldsAutoSet + `) values (?, ?, ?, ?, ?, ?)`
	return m.ExecNoCache(query, data.UserId, data.HeadImg, data.NameOne, data.RelationOne, data.NameTwo, data.RelationTwo)
}

func (m *MemorialModel) TxInsert(tx sqlx.TxSession, data Memorial) (sql.Result, error) {
	memorialIdKey := fmt.Sprintf("%s%v", cacheCymzjsMemorialIdPrefix, data.Id)
	_ = m.DelCache(memorialIdKey)
	query := `insert into ` + m.table + ` (` + memorialFieldsAutoSet + `) values (?, ?, ?, ?, ?, ?)`
	return tx.Exec(query, data.UserId, data.HeadImg, data.NameOne, data.RelationOne, data.NameTwo, data.RelationTwo)
}

func (m *MemorialModel) FindOne(id int64) (*Memorial, error) {
	memorialIdKey := fmt.Sprintf("%s%v", cacheCymzjsMemorialIdPrefix, id)
	var dest Memorial
	err := m.Query(&dest, memorialIdKey, func(conn sqlx.Conn, v interface{}) error {
		query := `select ` + memorialFields + ` from ` + m.table + ` where id = ? limit 1`
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

func (m *MemorialModel) FindMany(ids []int64, workers ...int) (list []*Memorial) {
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
		list = append(list, one.(*Memorial))
	}

	sort.Slice(list, func(i, j int) bool {
		return gutil.IndexOf(list[i].Id, ids) < gutil.IndexOf(list[j].Id, ids)
	})

	return
}

func (m *MemorialModel) Update(data Memorial) error {
	memorialIdKey := fmt.Sprintf("%s%v", cacheCymzjsMemorialIdPrefix, data.Id)
	_, err := m.Exec(func(conn sqlx.Conn) (result sql.Result, err error) {
		query := `update ` + m.table + ` set ` + memorialFieldsWithPlaceHolder + ` where id = ?`
		return conn.Exec(query, data.UserId, data.HeadImg, data.NameOne, data.RelationOne, data.NameTwo, data.RelationTwo, data.Id)
	}, memorialIdKey)
	return err
}

func (m *MemorialModel) UpdatePartial(ms ...g.Map) (err error) {
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

func (m *MemorialModel) updatePartial(data g.Map) error {
	updateArgs, err := sqlx.ExtractUpdateArgs(memorialFieldList, data)
	if err != nil {
		return err
	}

	memorialIdKey := fmt.Sprintf("%s%v", cacheCymzjsMemorialIdPrefix, updateArgs.Id)
	_, err = m.Exec(func(conn sqlx.Conn) (result sql.Result, err error) {
		query := `update ` + m.table + ` set ` + updateArgs.Fields + ` where id = ` + updateArgs.Id
		return conn.Exec(query, updateArgs.Args...)
	}, memorialIdKey)
	return err
}

func (m *MemorialModel) TxUpdate(tx sqlx.TxSession, data Memorial) error {
	memorialIdKey := fmt.Sprintf("%s%v", cacheCymzjsMemorialIdPrefix, data.Id)
	_, err := m.Exec(func(conn sqlx.Conn) (result sql.Result, err error) {
		query := `update ` + m.table + ` set ` + memorialFieldsWithPlaceHolder + ` where id = ?`
		return tx.Exec(query, data.UserId, data.HeadImg, data.NameOne, data.RelationOne, data.NameTwo, data.RelationTwo, data.Id)
	}, memorialIdKey)
	return err
}

func (m *MemorialModel) TxUpdatePartial(tx sqlx.TxSession, ms ...g.Map) (err error) {
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

func (m *MemorialModel) txUpdatePartial(tx sqlx.TxSession, data g.Map) error {
	updateArgs, err := sqlx.ExtractUpdateArgs(memorialFieldList, data)
	if err != nil {
		return err
	}

	memorialIdKey := fmt.Sprintf("%s%v", cacheCymzjsMemorialIdPrefix, updateArgs.Id)
	_, err = m.Exec(func(conn sqlx.Conn) (result sql.Result, err error) {
		query := `update ` + m.table + ` set ` + updateArgs.Fields + ` where id = ` + updateArgs.Id
		return tx.Exec(query, updateArgs.Args...)
	}, memorialIdKey)
	return err
}

func (m *MemorialModel) Delete(id ...int64) error {
	if len(id) == 0 {
		return nil
	}

	keys := make([]string, len(id))
	for i, v := range id {
		keys[i] = fmt.Sprintf("%s%v", cacheCymzjsMemorialIdPrefix, v)
	}

	_, err := m.Exec(func(conn sqlx.Conn) (result sql.Result, err error) {
		query := fmt.Sprintf(`delete from `+m.table+` where id in (%s)`, sqlx.In(len(id)))
		return conn.Exec(query, gconv.Interfaces(id)...)
	}, keys...)
	return err
}

func (m *MemorialModel) TxDelete(tx sqlx.TxSession, id ...int64) error {
	if len(id) == 0 {
		return nil
	}

	keys := make([]string, len(id))
	for i, v := range id {
		keys[i] = fmt.Sprintf("%s%v", cacheCymzjsMemorialIdPrefix, v)
	}

	_, err := m.Exec(func(conn sqlx.Conn) (result sql.Result, err error) {
		query := fmt.Sprintf(`delete from `+m.table+` where id in (%s)`, sqlx.In(len(id)))
		return tx.Exec(query, gconv.Interfaces(id)...)
	}, keys...)
	return err
}
