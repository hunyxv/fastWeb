package fastweb

import (
	"bytes"
	"errors"
	"reflect"
	"regexp"
	"strconv"
	"time"
)

type field struct {
	name      string
	index     int
	t         reflect.Type
	required  bool   // 是否必须
	maxlength int    // 最大长度(字段必须是 string类型）
	minlength int    // 最小长度(字段必须是 string类型）
	strip     bool   // 是否自动去除值两侧的空白字符（字段必须是 string类型）
	regular   string // 自定义正则表达式（字段必须是 string类形象）
	format    string // 日期类型格式化
}

func (f *field) setvalue(obj reflect.Value, v []byte) error {
	if len(v) == 0 && f.required {
		return errors.New("Missing required parameters")
	}

	fv := obj.Field(f.index)
	switch f.t.Kind() {
	case reflect.String:
		if len(f.format) > 0 {
			datetime, err := time.Parse(f.format, b2s(v))
			if err != nil {
				return err
			}
			fv.Set(reflect.ValueOf(datetime))
			return nil
		}

		if f.strip {
			v = bytes.TrimSpace(v)
		}

		if f.maxlength > 0 && len(v) > f.maxlength {
			return errors.New("The value is too long")
		}

		if f.minlength > 0 && len(v) < f.minlength {
			return errors.New("The value is too short")
		}

		if len(f.regular) > 0 {
			ismatch, err := regexp.Match(f.regular, v)
			if err != nil {
				return err
			}
			if !ismatch {
				return errors.New("Regular match failed")
			}
		}
		fv.SetString(b2s(v))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		n, err := strconv.Atoi(b2s(v))
		if err != nil {
			return err
		}
		fv.SetInt(int64(n))
	case reflect.Float32, reflect.Float64:
		f, err := strconv.ParseFloat(b2s(v), len(v))
		if err != nil {
			return err
		}
		fv.SetFloat(f)
	}
	return nil
}
