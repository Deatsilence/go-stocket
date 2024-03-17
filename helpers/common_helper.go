package helpers

import (
	"reflect"
)

func HasValue(values ...interface{}) bool {
	for _, value := range values {
		v := reflect.ValueOf(value)
		switch v.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if v.Int() == 0 {
				return false
			}
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			if v.Uint() == 0 {
				return false
			}
		case reflect.Float32, reflect.Float64:
			if v.Float() == 0 {
				return false
			}
		case reflect.Bool:
			if !v.Bool() {
				return false
			}
		case reflect.String:
			if v.String() == "" {
				return false
			}
		default:
			if v.IsZero() {
				return false
			}
		}
	}
	return true
}
