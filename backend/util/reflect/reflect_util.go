package reflect_util

import (
	"fmt"
	"reflect"
	"strconv"
)

func TrySetStuctValueFromString(o interface{}, f, v string) {
	SetStuctValueFromString(o, f, v)
}

func SetStuctValueFromString(o interface{}, field, value string) error {
	rValue := reflect.ValueOf(o)
	s := rValue.Elem()

	if s.Kind() == reflect.Struct {
		f := s.FieldByName(field)
		if f.IsValid() && f.CanSet() {
			switch f.Kind() {
			case reflect.String:
				f.SetString(value)
			case reflect.Bool:
				if s, err := strconv.ParseBool(value); err == nil {
					f.SetBool(s)
				} else {
					return fmt.Errorf("value `%v` is not a bool", value)
				}
			case reflect.Int8:
				fallthrough
			case reflect.Int16:
				fallthrough
			case reflect.Int32:
				fallthrough
			case reflect.Int64:
				fallthrough
			case reflect.Int:
				if s, err := strconv.Atoi(value); err == nil {
					x := int64(s)
					if !f.OverflowInt(x) {
						f.SetInt(x)
					}
				} else {
					return fmt.Errorf("value `%v` is not a integer", value)
				}
			default:
				return fmt.Errorf("unsupport type: `%v`", f.Kind())
			}
		}
	}
	return nil
}

func GetAllFieldNamesOfStruct(o interface{}) []string {
	rType := reflect.ValueOf(o).Elem().Type()
	if rType.Kind() != reflect.Struct {
		return make([]string, 0)
	}
	ret := make([]string, rType.NumField())
	for i := 0; i < rType.NumField(); i++ {
		ret = append(ret, rType.Field(i).Name)
	}
	return ret
}

func GetTagOfStruct(o interface{}, fieldName string) (reflect.StructField, reflect.StructTag) {
	field, _ := reflect.TypeOf(o).Elem().FieldByName(fieldName)
	return field, field.Tag
}

func ContainField(o interface{}, fieldName string) bool {
	rType := reflect.ValueOf(o).Elem().Type()
	if rType.Kind() != reflect.Struct {
		return false
	}
	_, has := rType.FieldByName(fieldName)
	return has
}
