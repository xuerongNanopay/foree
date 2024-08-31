package reflect_util

import (
	"reflect"
	"strconv"
)

func SetStuctValueFromString(o interface{}, f, v string) {
	rValue := reflect.ValueOf(o)
	s := rValue.Elem()

	if s.Kind() == reflect.Struct {
		f := s.FieldByName(f)
		if f.IsValid() && f.CanSet() {
			switch f.Kind() {
			case reflect.String:
				f.SetString(v)
			case reflect.Bool:
				if s, err := strconv.ParseBool(v); err == nil {
					f.SetBool(s)
				}
			case reflect.Int | reflect.Int16 | reflect.Int32 | reflect.Int64 | reflect.Int8:
				if s, err := strconv.Atoi(v); err == nil {
					x := int64(s)
					if !f.OverflowInt(x) {
						f.SetInt(x)
					}
				}
			}
		}
	}
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
