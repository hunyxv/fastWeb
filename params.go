package fastweb

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

var cache sync.Map

//map[string]*params = make(map[string]*params)

// TODO add gt lt qe le ge ne 比较符
type field struct {
	name      string
	key       string
	index     int
	t         reflect.Type
	required  bool   // 是否必须
	maxlength int    // 最大长度(字段必须是 string类型）
	minlength int    // 最小长度(字段必须是 string类型）
	strip     bool   // 是否自动去除值两侧的空白字符（字段必须是 string类型）
	re        string // 自定义正则表达式（字段必须是 string类形象）
	format    string // 日期类型格式化
}

func (f *field) setvalue(obj reflect.Value, v []byte) error {
	if len(v) == 0 && f.required {
		return errors.New("Missing required parameters")
	}

	fv := obj.Field(f.index)
	switch f.t.Kind() {
	case reflect.String:
		if f.strip {
			v = bytes.TrimSpace(v)
		}

		if f.maxlength > 0 && len(v) > f.maxlength {
			return errors.New("The value is too long")
		}

		if f.minlength > 0 && len(v) < f.minlength {
			return errors.New("The value is too short")
		}

		if len(f.re) > 0 {
			ismatch, err := regexp.Match(f.re, v)
			if err != nil {
				return err
			}
			if !ismatch {
				return errors.New("Regular match failed")
			}
		}
		fv.SetString(b2s(v))
	case reflect.Struct:
		if len(f.format) > 0 && f.t.Name() == "Time" {
			datetime, err := time.Parse(f.format, b2s(v))
			if err != nil {
				return err
			}
			fv.Set(reflect.ValueOf(datetime))
			return nil
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		n, err := strconv.Atoi(b2s(v))
		if err != nil {
			return err
		}
		fv.SetInt(int64(n))
	case reflect.Float32, reflect.Float64:
		float, err := strconv.ParseFloat(b2s(v), len(v))
		if err != nil {
			return err
		}
		fv.SetFloat(float)
	}
	return nil
}

func (f *field) tagparse(tag string) error {
	opts := strings.Split(tag, ",")
	key, opts := opts[0], opts[1:]
	if len(key) > 0 {
		f.key = key
	} else {
		f.key = f.name
	}

	for _, opt := range opts {
		switch opt[0] {
		case 'r':
			switch opt[2] {
			case 'q':
				f.required = true
			case '=':
				re := strings.SplitN(opt, "=", 2)
				if len(re) > 1 {
					f.re = re[1]
				}
			}
		case 'm':
			switch opt[1] {
			case 'a':
				sl := strings.SplitN(opt, "=", 2)
				l, err := strconv.ParseInt(sl[1], 10, 64)
				if err != nil {
					return err
				}
				f.maxlength = int(l)
			case 'i':
				sl := strings.SplitN(opt, "=", 2)
				l, err := strconv.ParseInt(sl[1], 10, 64)
				if err != nil {
					return err
				}
				f.minlength = int(l)
			}
		case 's':
			f.strip = true
		case 'f':
			format := strings.SplitN(opt, "=", 2)
			if len(format) > 1 {
				f.format = format[1]
			}
		}
	}
	return nil
}

type params struct {
	name   string
	fields map[string]*field
}

func (p *params) padding(key, value []byte, obj interface{}) error {
	vobj := reflect.ValueOf(obj)
	if vobj.Kind() != reflect.Ptr {
		return fmt.Errorf("non-pointer: %s", vobj.Type().Name())
	}

	f, ok := p.fields[b2s(key)]
	if ok {
		err := f.setvalue(vobj.Elem(), value)
		return err
	}
	return nil
}

func (p *params) valid(obj interface{}) error {
	v := reflect.ValueOf(obj).Elem()
	for _, f := range p.fields {
		if f.required {
			fv := v.FieldByName(f.name)
			if fv.IsZero() {
				return fmt.Errorf("%s[%s] is required", f.name, f.key)
			}
		}
	}
	return nil
}

func scan(obj interface{}) (*params, error) {
	t := reflect.TypeOf(obj)
	if t.Kind() != reflect.Ptr {
		return nil, fmt.Errorf("non-pointer: %s", t.Name())
	}
	t = t.Elem()
	pname := t.Name()

	if p, ok := cache.Load(pname); ok { // cache[pname]; ok {
		return p.(*params), nil
	}

	p := &params{name: pname, fields: make(map[string]*field, t.NumField())}
	for i := 0; i < t.NumField(); i++ {
		sf := t.Field(i)
		tag := sf.Tag.Get("valid")
		f := &field{name: sf.Name, index: i, t: sf.Type}
		err := f.tagparse(tag)
		if err != nil {
			return nil, err
		}
		p.fields[f.key] = f
	}
	cache.Store(pname, p)
	return p, nil
}
