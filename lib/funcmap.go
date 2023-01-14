package lib

import (
	"bytes"
	"fmt"
	"strconv"

	"github.com/astaxie/beego"
)

//AddFuncMaps .
func AddFuncMaps() {
	_ = beego.AddFuncMap("field_error_message", func(v map[string]map[string]string, key string) map[string]string {
		if val, ok := v[key]; ok {
			return val
		}
		return make(map[string]string)
	})
	_ = beego.AddFuncMap("field_error_exist", func(v map[string]map[string]string, key string) bool {
		if _, ok := v[key]; ok {
			return true
		}
		return false
	})
	_ = beego.AddFuncMap("printkb", func(i interface{}) string {
		switch v := i.(type) {
		case uint64:
			return num2str(int64(i.(uint64)/1024), '\u00A0')
		case int64:
			return num2str(i.(int64)/1024, '\u00A0')
		default:
			beego.Error("Unknown type:", v)
		}
		return "Mapping error"
	})
	_ = beego.AddFuncMap("printmb", func(i interface{}) string {
		switch v := i.(type) {
		case uint64:
			return num2str(int64(i.(uint64)/1024/1024), '\u00A0')
		case int64:
			return num2str(i.(int64)/1024/1024, '\u00A0')
		default:
			beego.Error("Unknown type:", v)
		}
		return "Mapping error"
	})
	_ = beego.AddFuncMap("printmbold", func(i uint64) string {
		return num2str(int64(i/1024/1024), ' ')
	})
	_ = beego.AddFuncMap("printgb", func(i uint64) string {
		return num2str(int64(i/1024/1024/1024), ' ')
	})
	_ = beego.AddFuncMap("percent", func(x, y interface{}) string {
		beego.Notice("Percent", x, y)
		zValue := "0"
		switch v := x.(type) {
		case string:
			beego.Error("Not implemented")
		case int32:
			if x.(int32) == 0 || y.(int32) == 0 {
				return zValue
			}
			a := float64(x.(int32))
			b := float64(y.(int32))
			return fmt.Sprintf("%d", int((a/b)*float64(100)))
		case int64:
			if x.(int64) == 0 || y.(int64) == 0 {
				return zValue
			}
			a := float64(x.(int64))
			b := float64(y.(int64))
			return fmt.Sprintf("%d", int((a/b)*float64(100)))
		case uint64:
			if x.(uint64) == 0 || y.(uint64) == 0 {
				return zValue
			}
			a := float64(x.(uint64))
			b := float64(y.(uint64))
			return fmt.Sprintf("%d", int((a/b)*float64(100)))
		default:
			beego.Error("Unknown type:", v)
		}
		return "Mapping error"
	})
}

func num2str(n int64, sep rune) string {
	s := strconv.FormatInt(n, 10)
	startOffset := 0
	var buff bytes.Buffer
	if n < 0 {
		startOffset = 1
		buff.WriteByte('-')
	}
	l := len(s)
	commaIndex := 3 - ((l - startOffset) % 3)
	if commaIndex == 3 {
		commaIndex = 0
	}
	for i := startOffset; i < l; i++ {
		if commaIndex == 3 {
			buff.WriteRune(sep)
			commaIndex = 0
		}
		commaIndex++
		buff.WriteByte(s[i])
	}
	return buff.String()
}
