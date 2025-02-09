package v2

import (
	"time"
	"unsafe"
	"maps"

	"github.com/json-iterator/go"
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
}

var _localeAlias = map[string]*time.Location{
	Local: time.Local,
	UTC:   time.UTC,
}

type CustomTimeExtension struct {
	jsoniter.DummyExtension
	formatAlias 	map[string]string
	localeAlias 	map[string]*time.Location
	defaultFormat 	string
	defaultLocale 	*time.Location
}

func NewCustomTimeExtension() *CustomTimeExtension {
	return &CustomTimeExtension{
		formatAlias: maps.Clone(_formatAlias),
		localeAlias: maps.Clone(_localeAlias),
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
				str := lt.Format(timeFormat)
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

			str := iter.ReadString()
			var t *time.Time
			if str != "" {
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
