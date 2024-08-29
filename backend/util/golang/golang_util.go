package golang_util

import "reflect"

func SetStringValueIfFieldExist(input interface{}, f, v string) {
	rValue := reflect.ValueOf(input)
	s := rValue.Elem()

	if s.Kind() == reflect.Struct {
		f := s.FieldByName(f)
		if f.IsValid() && f.CanSet() && f.Kind() == reflect.String {
			f.SetString(v)
		}
	}
}
