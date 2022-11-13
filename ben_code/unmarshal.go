package ben_code

import (
	"BitTorrentClient/consts"
	"errors"
	"fmt"
	"io"
	"reflect"
	"strings"
)

//Unmarshal 反序列化
func Unmarshal(r io.Reader, s interface{}) error {
	v := reflect.ValueOf(s)
	if v.Kind() != reflect.Ptr {
		//必须传入指针
		return errors.New("todo1")
	}
	ben, err := DecodeBencode(r)
	if err != nil {
		return errors.New("todo2")
	}
	switch ben.BenType {
	case consts.INTEGER:
		num := ben.IntegerWithNoErr()
		v.Elem().SetInt(num)
	case consts.STRING:
		str := ben.StrWithNoErr()
		v.Elem().SetString(str)
	case consts.LIST:
		//另开辟空间
		lists := ben.ListWithNoErr()
		fmt.Println("list:", lists)
		reflectList := reflect.MakeSlice(v.Elem().Type(), len(lists), len(lists))
		v.Elem().Set(reflectList)

		err = copyList(v, lists)
		if err != nil {
			return err
		}
	case consts.DICT:
		dict := ben.DictWithNoErr()
		err = copyDict(v, dict)
	}

	return nil
}
func copyList(reflectPtr reflect.Value, lists []*BenCode) error {
	vele := reflectPtr.Elem()
	if len(lists) == 0 {
		return nil
	}
	fmt.Println(len(lists))
	for i := range lists {
		switch lists[i].BenType {
		case consts.STRING:
			vele.Index(i).SetString(lists[i].StrWithNoErr())
		case consts.INTEGER:
			vele.Index(i).SetInt(lists[i].IntegerWithNoErr())

		case consts.LIST:
			//列表套列表
			if vele.Index(i).Kind() != reflect.Slice {
				//类型不正确 无法接收
				return errors.New("err type")
			}
			subList := lists[i].ListWithNoErr()
			subReflectType := vele.Index(i).Type()
			subReflectList := reflect.MakeSlice(subReflectType, len(subList), len(subList))
			subReflectListPtr := reflect.New(vele.Type().Elem()) //封装一层
			subReflectListPtr.Elem().Set(subReflectList)

			err := copyList(subReflectListPtr, subList)
			if err != nil {
				return err
			}
			vele.Index(i).Set(subReflectListPtr.Elem()) //vele是多维数组 给元素开辟空间
		case consts.DICT:
			if vele.Index(i).Kind() != reflect.Struct {
				return errors.New("err type")
			}
			dict := lists[i].DictWithNoErr()
			reflectStructPtr := reflect.New(vele.Index(i).Type())
			err := copyDict(reflectStructPtr, dict)
			if err != nil {
				return err
			}
			vele.Index(i).Set(reflectStructPtr.Elem())
		}
	}
	return nil
}

func reflectIsInteger(kind reflect.Kind) bool {
	switch kind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint64, reflect.Uint32:
		return true
	default:
		return false
	}
}
func copyDict(reflectPtr reflect.Value, dict map[string]*BenCode) error {
	vele := reflectPtr.Elem()

	for i := 0; i < vele.NumField(); i++ {
		field := vele.Field(i)
		if !field.CanSet() {
			continue
		}
		fieldType := vele.Type().Field(i)
		tag := fieldType.Tag.Get(consts.BenCodeTag)
		if tag == "" {
			tag = strings.ToLower(fieldType.Name)
		}
		dictField := dict[tag]
		if dictField == nil {
			continue
		}
		switch dictField.BenType {
		case consts.STRING:
			if fieldType.Type.Kind() != reflect.String {
				continue //类型不匹配
			}
			str := dictField.StrWithNoErr()
			field.SetString(str)

		case consts.INTEGER:
			if !reflectIsInteger(fieldType.Type.Kind()) {
				continue
			}
			num := dictField.IntegerWithNoErr()
			field.SetInt(num)

		case consts.LIST:
			if fieldType.Type.Kind() != reflect.Slice {
				continue
			}
			subList := dictField.ListWithNoErr()
			//开辟空间
			reflectList := reflect.MakeSlice(fieldType.Type, len(subList), len(subList))
			reflectListPtr := reflect.New(fieldType.Type)
			reflectListPtr.Elem().Set(reflectList) //包一层

			err := copyList(reflectListPtr, subList)
			if err != nil {
				return err
			}
			field.Set(reflectListPtr.Elem())

		case consts.DICT:
			if fieldType.Type.Kind() != reflect.Struct {
				continue
			}
			reflectStructPtr := reflect.New(fieldType.Type)
			subDict := dictField.DictWithNoErr()
			err := copyDict(reflectStructPtr, subDict)
			if err != nil {
				return err
			}
			field.Set(reflectStructPtr.Elem())

		}
	}
	return nil
}
