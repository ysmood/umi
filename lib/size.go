package lib

import (
	"reflect"
)

// Size ...
func Size(v interface{}) uintptr {
	t := reflect.TypeOf(v)

	var s uintptr

	if t == nil {
		return s
	}

	switch t.Kind() {
	case reflect.Invalid:

	case reflect.Slice:
		fallthrough
	case reflect.Array:
		s += t.Size()

		b, ok := v.([]byte)

		if ok {
			s += uintptr(len(b))
			break
		}

		val := reflect.ValueOf(v)
		l := val.Len()
		for i := 0; i < l; i++ {
			s += Size(val.Index(i).Interface())
		}

	// TODO: The size of the hash map should also contains the hash keys and collision linked lists,
	// but the internal data structure is invisible.
	case reflect.Map:
		s += t.Size()
		val := reflect.ValueOf(v)
		keys := val.MapKeys()

		for _, i := range keys {
			s += Size(i.Interface())
			s += Size(val.MapIndex(i).Interface())
		}

	case reflect.String:
		s += t.Size()
		s += uintptr(len(v.(string)))

	case reflect.Struct:
		val := reflect.ValueOf(v)
		reflect.TypeOf(v).Kind()

		l := val.NumField()

		for i := 0; i < l; i++ {
			field := val.Field(i)
			s += Size(field.Interface())
		}

	case reflect.Ptr:
		s += t.Size()
		val := reflect.ValueOf(v)
		s += Size(val.Elem().Interface())

	case reflect.Interface:
		s += t.Size()
		s += Size(v)

	default:
		s += t.Size()
	}

	return s
}
