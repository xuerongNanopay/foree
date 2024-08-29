package reflect_util

import "reflect"

func SetStringValueIfFieldExist(o interface{}, f, v string) {
	rValue := reflect.ValueOf(o)
	s := rValue.Elem()

	if s.Kind() == reflect.Struct {
		f := s.FieldByName(f)
		if f.IsValid() && f.CanSet() && f.Kind() == reflect.String {
			f.SetString(v)
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
