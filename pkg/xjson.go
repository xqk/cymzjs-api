package pkg

import (
	"fmt"
	jsoniter "github.com/json-iterator/go"
)

// CopyTo 将source结构体内容转换为dest结构体
func CopyTo(source, dest interface{}, debug ...bool) (err error) {
	if source == nil {
		return
	}

	if len(debug) > 0 && debug[0] {
		bytes, _ := jsoniter.MarshalIndent(source, "", "  ")
		fmt.Println("来源", string(bytes))
	}

	bytes, err := jsoniter.Marshal(source)
	if err != nil {
		return
	}

	err = jsoniter.Unmarshal(bytes, &dest)
	if err != nil {
		return
	}

	return
}

// JsonTo 将字符串转换为dest结构体
func JsonTo(json string, dest interface{}) (err error) {
	if json == "" {
		return nil
	}

	err = jsoniter.Unmarshal([]byte(json), dest)
	if err != nil {
		return
	}
	return nil
}

// ToStr 将结构体转换为字符串
func ToStr(dest interface{}) string {
	bytes, err := jsoniter.Marshal(dest)
	if err != nil {
		return ""
	}
	return string(bytes)
}
