package jsontime

import (
	"maps"
	"strconv"
	"time"
	"unsafe"

	jsoniter "github.com/json-iterator/go"
)

// time format alias
const (
	ANSIC       = "ANSIC"
	UnixDate    = "UnixDate"
	RubyDate    = "RubyDate"
	RFC822      = "RFC822"
	RFC822Z     = "RFC822Z"
	RFC850      = "RFC850"
	RFC1123     = "RFC1123"
	RFC1123Z    = "RFC1123Z"
	RFC3339     = "RFC3339"
	RFC3339Nano = "RFC3339Nano"
	Kitchen     = "Kitchen"
	Stamp       = "Stamp"
	StampMilli  = "StampMilli"
	StampMicro  = "StampMicro"
	StampNano   = "StampNano"
	DateTime    = "DateTime"
	DateOnly    = "DateOnly"
	TimeOnly    = "TimeOnly"

	// Unix timestamp formats
	Unix      = "unix"
	UnixMilli = "unixmilli"
	UnixMicro = "unixmicro"
	UnixNano  = "unixnano"
)

// time zone alias
const (
	Local = "Local"
	UTC   = "UTC"
)

const (
	tagNameTimeFormat   = "time_format"
	tagNameTimeLocation = "time_location"
)

var _formatAlias = map[string]string{
	ANSIC:       time.ANSIC,
	UnixDate:    time.UnixDate,
	RubyDate:    time.RubyDate,
	RFC822:      time.RFC822,
	RFC822Z:     time.RFC822Z,
	RFC850:      time.RFC850,
	RFC1123:     time.RFC1123,
	RFC1123Z:    time.RFC1123Z,
	RFC3339:     time.RFC3339,
	RFC3339Nano: time.RFC3339Nano,
	Kitchen:     time.Kitchen,
	Stamp:       time.Stamp,
	StampMilli:  time.StampMilli,
	StampMicro:  time.StampMicro,
	StampNano:   time.StampNano,
	DateTime:    time.DateTime,
	DateOnly:    time.DateOnly,
	TimeOnly:    time.TimeOnly,

	// Unix timestamp formats use themselves as keys
	Unix:      Unix,
	UnixMilli: UnixMilli,
	UnixMicro: UnixMicro,
	UnixNano:  UnixNano,
}

var _localeAlias = map[string]*time.Location{
	Local: time.Local,
	UTC:   time.UTC,
}

// isUnixFormat returns whether the format is one of the Unix timestamp formats
func isUnixFormat(format string) bool {
	switch format {
	case Unix, UnixMilli, UnixMicro, UnixNano:
		return true
	default:
		return false
	}
}

type CustomTimeExtension struct {
	jsoniter.DummyExtension
	formatAlias   map[string]string
	localeAlias   map[string]*time.Location
	defaultFormat string
	defaultLocale *time.Location
}

func NewCustomTimeExtension() *CustomTimeExtension {
	return &CustomTimeExtension{
		formatAlias:   maps.Clone(_formatAlias),
		localeAlias:   maps.Clone(_localeAlias),
		defaultFormat: time.RFC3339,
		defaultLocale: time.Local,
	}
}

func (extension *CustomTimeExtension) AddTimeFormatAlias(alias, format string) {
	extension.formatAlias[alias] = format
}

func (extension *CustomTimeExtension) AddLocaleAlias(alias string, locale *time.Location) {
	extension.localeAlias[alias] = locale
}

func (extension *CustomTimeExtension) SetDefaultTimeFormat(timeFormat string, timeLocation *time.Location) {
	extension.defaultFormat = timeFormat
	extension.defaultLocale = timeLocation
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

		timeFormat := extension.defaultFormat
		formatTag := binding.Field.Tag().Get(tagNameTimeFormat)
		if format, ok := extension.formatAlias[formatTag]; ok {
			timeFormat = format
		} else if formatTag != "" {
			timeFormat = formatTag
		}
		locale := extension.defaultLocale
		if localeTag := binding.Field.Tag().Get(tagNameTimeLocation); localeTag != "" {
			if loc, ok := extension.localeAlias[localeTag]; ok {
				locale = loc
			} else {
				loc, err := time.LoadLocation(localeTag)
				if err != nil {
					typeErr = err
				} else {
					locale = loc
				}
			}
		}

		binding.Encoder = &funcEncoder{fun: func(ptr unsafe.Pointer, stream *jsoniter.Stream) {
			if typeErr != nil {
				stream.Error = typeErr
				return
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

				if isUnixFormat(timeFormat) {
					var timestamp int64
					switch timeFormat {
					case Unix:
						timestamp = lt.Unix()
					case UnixMilli:
						timestamp = lt.UnixMilli()
					case UnixMicro:
						timestamp = lt.UnixMicro()
					case UnixNano:
						timestamp = lt.UnixNano()
					}
					stream.WriteInt64(timestamp)
				} else {
					str := lt.Format(timeFormat)
					stream.WriteString(str)
				}
			} else {
				stream.Write([]byte("null"))
			}
		}}
		binding.Decoder = &funcDecoder{fun: func(ptr unsafe.Pointer, iter *jsoniter.Iterator) {
			if typeErr != nil {
				iter.Error = typeErr
				return
			}

			var t *time.Time

			if isUnixFormat(timeFormat) {
				switch iter.WhatIsNext() {
				case jsoniter.NumberValue:
					tt := unixTimeDecoder(iter.ReadInt64(), timeFormat)
					t = &tt
				case jsoniter.StringValue:
					if str := iter.ReadString(); str != "" {
						timestamp, _ := strconv.ParseInt(str, 10, 64)
						tt := unixTimeDecoder(timestamp, timeFormat)
						t = &tt
					}
				case jsoniter.NilValue:
					iter.ReadNil()
					t = nil
				default:
					iter.ReportError("time.Time", "expect number or string")
					return
				}
			} else {
				if str := iter.ReadString(); str != "" {
					var err error
					tmp, err := time.ParseInLocation(timeFormat, str, locale)
					if err != nil {
						iter.Error = err
						return
					}
					t = &tmp
				} else {
					t = nil
				}
			}

			if isPtr {
				tpp := (**time.Time)(ptr)
				*tpp = t
			} else {
				tp := (*time.Time)(ptr)
				if tp != nil && t != nil {
					*tp = *t
				}
			}
		}}
	}
}

func unixTimeDecoder(t int64, format string) time.Time {
	switch format {
	case Unix:
		return time.Unix(t, 0)
	case UnixMilli:
		return time.UnixMilli(t)
	case UnixMicro:
		return time.UnixMicro(t)
	case UnixNano:
		return time.Unix(0, t)
	default:
		return time.Unix(t, 0)
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
