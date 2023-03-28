package pkg

import (
	"fmt"
	"time"
)

func Duration(s, e time.Time) (time.Duration, error) {
	if e.Unix() < s.Unix() {
		return 0, nil
	}
	return e.Sub(s), nil
}

// Sec2MS 秒数转为"分:秒"格式
func Sec2MS(second int64, minuteTpl string, secondTpl string) (t string) {
	if second == 0 {
		return ""
	}

	if minuteTpl == "" {
		minuteTpl = "%d:"
	}
	if secondTpl == "" {
		secondTpl = "%d"
	}

	var tmp int64 = second
	if m := tmp / 60; m > 0 {
		t += fmt.Sprintf(minuteTpl, m)
		tmp = tmp % 60
	}
	if tmp > 0 {
		t += fmt.Sprintf(secondTpl, tmp)
	}
	return t
}
