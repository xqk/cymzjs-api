package logicutil

type (
	LogPos int64
)

const (
	LogPos11001 LogPos = 11001 // 任意页面-进入页面
)

func (t LogPos) Int64() int64 {
	return int64(t)
}

func (t LogPos) Int() int {
	return int(t)
}
