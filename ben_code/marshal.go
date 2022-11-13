package ben_code

import (
	"BitTorrentClient/consts"
	"io"
	"reflect"
	"strings"
)

//Marshal 序列化
func Marshal(w io.Writer, s interface{}) error {
	v := reflect.ValueOf(s)
	if v.Kind() == reflect.Ptr {
		//若为指针读取出里层结构
		v = v.Elem()
	}
	marshalReflectValue(w, v)
	return nil
}
func marshalReflectValue(w io.Writer, v reflect.Value) {
	//拆出来 用于接下来的互递归
	switch v.Kind() {
	case reflect.String:
		EncodeStr(w, v.String())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint64, reflect.Uint32:
		EncodeInt(w, v.Int())
	case reflect.Slice:
		marshalList(w, v)
	case reflect.Struct:
		marshalDict(w, v)

	}
}
func marshalList(w io.Writer, v reflect.Value) {
	_, _ = w.Write([]byte{'l'})
	for i := 0; i < v.Len(); i++ {
		val := v.Index(i)
		marshalReflectValue(w, val)
	}
	_, _ = w.Write([]byte{'e'})
	return
}
func marshalDict(w io.Writer, v reflect.Value) {
	_, _ = w.Write([]byte{'d'})
	for i := 0; i < v.NumField(); i++ {
		vf := v.Field(i)
		vt := v.Type().Field(i)
		// 读取成员变量的tag 作为dict的key bencode:xxx
		tag := vt.Tag.Get(consts.BenCodeTag)
		if tag == "" {
			tag = strings.ToLower(vt.Name)
		}
		EncodeStr(w, tag)
		marshalReflectValue(w, vf)
	}
	_, _ = w.Write([]byte{'e'})
	return
}
