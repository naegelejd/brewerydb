package brewerydb

import (
	"net/url"
	"reflect"
	"strconv"
	"strings"
)

// type Foo struct {
// 	ID       int     `json:"id"`
// 	Name     string  `json:"name"`
// 	Min      float64 `json:"fMin"`
// 	Max      float64 `json:"fMax"`
// 	Mean     float64 `json:"average"`
// 	Results  []bool  `json:"results"`
// 	DescrPtr *string `json:"description"`
// 	Country  struct {
// 		Unused    []interface{}
// 		AlwaysNil *Foo
// 	} `json:"country"`
// }

// func example() {
// 	foo := Foo{}

// 	foo.Name = "Admiral"
// 	foo.Min = 5.7
// 	foo.Max = 13
// 	foo.Results = []bool{true, false, true, false}
// 	foo.DescrPtr = &foo.Name
// 	foo.Country.Unused = []interface{}{"US"}
// 	fmt.Println(encode(foo))
// }

func encode(data interface{}) url.Values {
	query := url.Values{}

	t := reflect.TypeOf(data)
	v := reflect.ValueOf(data)

	if t.Kind() == reflect.Ptr {
		if v.IsNil() {
			return query
		}
		v = reflect.Indirect(v)
		t = v.Type()
	}

	if t.Kind() != reflect.Struct {
		return query
	}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		var key string
		tag := field.Tag.Get("json")
		if tag != "" {
			// from src/encoding/json/tags.go
			if idx := strings.Index(tag, ","); idx != -1 {
				tag = tag[:idx]
			}
			key = tag
		} else {
			key = field.Name
		}

		val := v.Field(i)

		switch val.Kind() {
		case reflect.Array, reflect.Slice:
			for i := 0; i < val.Len(); i++ {
				sval := toString(val.Index(i))
				query.Add(key, sval)
			}
		default:
			if isEncodableValue(val) {
				sval := toString(val)
				query.Set(key, sval)
			}
		}

	}
	return query
}

func toString(v reflect.Value) string {
	var sval string
	switch v.Kind() {
	case reflect.String:
		sval = v.String()
	case reflect.Bool:
		if v.Bool() {
			sval = "true"
		} else {
			sval = "false"
		}
	case reflect.Float32, reflect.Float64:
		sval = strconv.FormatFloat(v.Float(), 'f', -1, 64)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		sval = strconv.FormatInt(v.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		sval = strconv.FormatUint(v.Uint(), 10)
	case reflect.Ptr:
		return toString(reflect.Indirect(v))
	}
	return sval
}

// isEncodableValue determines if a value is:
// 1. is of an encodable type, and
// 2. not the zero value
// TODO: stop caring about whether it's the zero-value and only
// omit it if tagged with "omitempty"
func isEncodableValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Array, reflect.Slice, reflect.String:
		return v.Len() != 0
	case reflect.Bool:
		return v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() != 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() != 0
	case reflect.Float32, reflect.Float64:
		return v.Float() != 0
	case reflect.Ptr:
		return isEncodableValue(reflect.Indirect(v))
	}
	return false
}
