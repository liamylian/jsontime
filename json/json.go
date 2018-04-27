package json

import (
	gojson "encoding/json"
	"strings"
	"reflect"
	"time"
)

type tagOptions string

const (
	SqlDateTag    = "sql_date"
	SqlTimeTag    = "sql_time"
	SqlDateFormat = "2006-01-02"
	SqlTimeFormat = "2006-01-02 15:04:05"
)

func Marshal(v interface{}) ([]byte, error) {
	tmp := make(map[string]interface{})

	typ := reflect.TypeOf(v)
	val := reflect.ValueOf(v)
	for i := 0; i < typ.NumField(); i++ {
		typeField := typ.Field(i)
		structField := val.Field(i)
		jsonName, jsonOptions := parseTag(typeField.Tag.Get("json"))
		if jsonName == "" {
			jsonName = typeField.Name
		}

		if t, isTime := structField.Interface().(time.Time); isTime {
			if jsonOptions.Contains(SqlDateTag) {
				tmp[jsonName] = t.Format(SqlDateFormat)
			} else if jsonOptions.Contains(SqlTimeTag) {
				tmp[jsonName] = t.Format(SqlTimeFormat)
			}
		} else {
			tmp[jsonName] = structField.Interface()
		}
	}

	return gojson.Marshal(tmp)
}

func Unmarshal(data []byte, v interface{}) error {
	tmp := make(map[string]interface{})
	if err := gojson.Unmarshal(data, &tmp); err != nil {
		return err
	}

	typ := reflect.TypeOf(v).Elem()
	val := reflect.ValueOf(v).Elem()
	for i := 0; i < typ.NumField(); i++ {
		structField := val.Field(i)
		typeField := typ.Field(i)

		if !structField.CanSet() {
			continue
		}

		if _, isT := isTime(structField.Interface()); isT {
			jsonName, jsonOptions := parseTag(typeField.Tag.Get("json"))
			if jsonName == "" {
				jsonName = typeField.Name
			}

			var format string
			if jsonOptions.Contains(SqlDateTag) {
				format = SqlDateFormat
			} else if jsonOptions.Contains(SqlTimeTag) {
				format = SqlTimeFormat
			} else {
				continue
			}

			if _, ok := tmp[jsonName]; !ok {
				continue
			}

			timeVal, ok := tmp[jsonName].(string)
			if !ok {
				continue
			}

			if t, err := time.ParseInLocation(format, timeVal, time.Local); err == nil {
				tmp[jsonName] = t
			} else {
				return err
			}

		}
	}

	if tmpBytes, err := gojson.Marshal(tmp); err == nil {
		return gojson.Unmarshal(tmpBytes, v)
	} else {
		return err
	}
}

func isTime(t interface{}) (*time.Time, bool) {
	if val, isTime := t.(time.Time); isTime {
		return &val, true
	} else if ptr, isTime := t.(*time.Time); isTime {
		return ptr, true
	} else {
		return nil, false
	}
}

func parseTag(tag string) (string, tagOptions) {
	if idx := strings.Index(tag, ","); idx != -1 {
		return tag[:idx], tagOptions(tag[idx+1:])
	}
	return tag, tagOptions("")
}

func (o tagOptions) Contains(optionName string) bool {
	if len(o) == 0 {
		return false
	}
	s := string(o)
	for s != "" {
		var next string
		i := strings.Index(s, ",")
		if i >= 0 {
			s, next = s[:i], s[i+1:]
		}
		if s == optionName {
			return true
		}
		s = next
	}
	return false
}

