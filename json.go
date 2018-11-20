package jsontime

import (
	"github.com/json-iterator/go"
	"time"
	"strconv"
	"unsafe"
)

var ConfigWithCustomTimeFormat = jsoniter.ConfigCompatibleWithStandardLibrary

func init() {
	ConfigWithCustomTimeFormat.RegisterExtension(&CustomTimeExtension{})
}

type CustomTimeExtension struct {
	jsoniter.DummyExtension
}

func (extension *CustomTimeExtension) UpdateStructDescriptor(structDescriptor *jsoniter.StructDescriptor) {

	for _, binding := range structDescriptor.Fields {
		var typeErr error
		var isPtr bool
		typeName := binding.Field.Type().String()

		if typeName == "time.Time" {
			isPtr = false
		} else if typeName == "*time.Time" {
			isPtr = true
		} else {
			continue
		}

		formatTag := binding.Field.Tag().Get("time_format")
		timeFormat := formatTag
		if timeFormat == "sql_datetime" {
			timeFormat = "2006-01-02 15:04:05"
		} else if timeFormat == "sql_date" {
			timeFormat = "2006-01-02"
		}

		locale := time.Local
		if isUTC, _ := strconv.ParseBool(binding.Field.Tag().Get("time_utc")); isUTC {
			locale = time.UTC
		}
		if locTag := binding.Field.Tag().Get("time_location"); locTag != "" {
			loc, err := time.LoadLocation(locTag)
			if err != nil {
				typeErr = err
			} else {
				locale = loc
			}
		}

		var isSnap bool
		snapTag := binding.Field.Tag().Get("time_snap")
		if snapTag == "" && (formatTag == "sql_datetime" || formatTag == "sql_date") {
			isSnap = true
		} else {
			isSnap, _ = strconv.ParseBool(binding.Field.Tag().Get("time_snap"))
		}

		binding.Encoder = &funcEncoder{fun: func(ptr unsafe.Pointer, stream *jsoniter.Stream) {
			if typeErr != nil {
				stream.Error = typeErr
				return
			}

			var format string
			if timeFormat == "" {
				format = time.RFC3339Nano
			} else {
				format = timeFormat
			}

			var tp *time.Time
			if isPtr {
				tpp := (**time.Time)(ptr)
				tp = *(tpp)
			} else {
				tp = (*time.Time)(ptr)
			}

			if tp != nil {
				lt := tp.In(locale)
				str := lt.Format(format)
				if formatTag == "sql_date" && (str == "0000-01-01" || (isSnap && lt.Unix() <= 0)) {
					str = "0000-00-00"
				} else if formatTag == "sql_datetime" && (str == "0000-01-01 00:00:00" || (isSnap && lt.Unix() <= 0)) {
					str = "0000-00-00 00:00:00"
				}
				stream.WriteString(str)
			} else {
				stream.Write([]byte("null"))
			}
		}}

		binding.Decoder = &funcDecoder{fun: func(ptr unsafe.Pointer, iter *jsoniter.Iterator) {
			if typeErr != nil {
				iter.Error = typeErr
				return
			}

			var format string
			if timeFormat == "" {
				format = time.RFC3339
			} else {
				format = timeFormat
			}

			str := iter.ReadString()
			var t *time.Time
			if str != "" {
				var err error
				tmp, err := time.ParseInLocation(format, str, locale)
				if err != nil {
					if _, ok := err.(*time.ParseError); ok {
						tmp = time.Date(0, 1, 1, 0, 0, 0, 0, locale)
					} else {
						iter.Error = err
						return
					}
				}
				t = &tmp
			} else {
				t = nil
			}

			if isPtr {
				tpp := (**time.Time)(ptr)
				*tpp = t
			} else {
				tp := (*time.Time)(ptr)
				if tp != nil && t != nil{
					*tp = *t
				}
			}
		}}
	}
}

type funcDecoder struct {
	fun jsoniter.DecoderFunc
}

func (decoder *funcDecoder) Decode(ptr unsafe.Pointer, iter *jsoniter.Iterator) {
	decoder.fun(ptr, iter)
}

type funcEncoder struct {
	fun         jsoniter.EncoderFunc
	isEmptyFunc func(ptr unsafe.Pointer) bool
}

func (encoder *funcEncoder) Encode(ptr unsafe.Pointer, stream *jsoniter.Stream) {
	encoder.fun(ptr, stream)
}

func (encoder *funcEncoder) IsEmpty(ptr unsafe.Pointer) bool {
	if encoder.isEmptyFunc == nil {
		return false
	}
	return encoder.isEmptyFunc(ptr)
}
