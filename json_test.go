package jsontime

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var json = ConfigWithCustomTimeFormat

func TestTimeFormat(t *testing.T) {
	type Book struct {
		Id          int        `json:"id"`
		PublishedAt *time.Time `json:"published_at"`
		UpdatedAt   *time.Time `json:"updated_at"`
		CreatedAt   time.Time  `json:"created_at"`
	}

	timeZone, err := time.LoadLocation("Asia/Shanghai")
	assert.Nil(t, err)
	SetDefaultTimeFormat(time.RFC3339, timeZone)
	t2018 := time.Date(2018, 1, 1, 0, 0, 0, 0, timeZone)
	book1 := Book{
		Id:        1,
		UpdatedAt: &t2018,
		CreatedAt: t2018,
	}
	bytes, err := json.Marshal(book1)
	assert.Nil(t, err)
	assert.Equal(t, `{"id":1,"published_at":null,"updated_at":"2018-01-01T00:00:00+08:00","created_at":"2018-01-01T00:00:00+08:00"}`, string(bytes))

	book2 := Book{}
	err = json.Unmarshal(bytes, &book2)
	assert.Nil(t, err)
	assert.Equal(t, book1, book2)
}

func TestLocale(t *testing.T) {
	timeZoneShanghai, _ := time.LoadLocation("Asia/Shanghai")
	AddLocaleAlias("shanghai", timeZoneShanghai)
	type Book struct {
		Id          int        `json:"id"`
		PublishedAt time.Time  `json:"published_at" time_location:"UTC"`
		UpdatedAt   *time.Time `json:"updated_at" time_location:"shanghai"`
		CreatedAt   time.Time  `json:"created_at" time_location:"shanghai"`
	}

	book := Book{
		Id:          0,
		PublishedAt: time.Date(2018, 1, 1, 0, 0, 0, 0, timeZoneShanghai),
		UpdatedAt:   nil,
		CreatedAt:   time.Date(2018, 1, 1, 0, 0, 0, 0, timeZoneShanghai),
	}

	bytes, err := json.Marshal(book)
	assert.Nil(t, err)
	assert.Equal(t, `{"id":0,"published_at":"2017-12-31T16:00:00Z","updated_at":null,"created_at":"2018-01-01T00:00:00+08:00"}`, string(bytes))
}

func TestUnMarshalZero(t *testing.T) {
	type Book struct {
		Id        int        `json:"id"`
		UpdatedAt *time.Time `json:"updated_at" time_location:"UTC"`
		CreatedAt time.Time  `json:"created_at" time_location:"UTC"`
	}
	book := Book{}
	jsonBytes := []byte(`{"id":0,"updated_at":null,"created_at":"0000-00-00 00:00:00"}`)

	err := json.Unmarshal(jsonBytes, &book)
	assert.NotNil(t, err)
}

func TestAlias(t *testing.T) {
	type Book struct {
		Id          int        `json:"id"`
		PublishedAt *time.Time `json:"published_at" time_format:"sql_datetime"`
		UpdatedAt   *time.Time `json:"updated_at" time_format:"sql_datetime"`
		CreatedAt   time.Time  `json:"created_at" time_format:"sql_datetime"`
	}

	timeZoneShanghai, err := time.LoadLocation("Asia/Shanghai")
	assert.Nil(t, err)

	AddTimeFormatAlias("sql_datetime", "2006-01-02 15:04:05")
	AddLocaleAlias("shanghai", timeZoneShanghai)
	t2018 := time.Date(2018, 1, 1, 0, 0, 0, 0, timeZoneShanghai)
	book1 := Book{
		Id:        1,
		UpdatedAt: &t2018,
		CreatedAt: t2018,
	}
	bytes, err := json.Marshal(book1)
	assert.Nil(t, err)
	assert.Equal(t, `{"id":1,"published_at":null,"updated_at":"2018-01-01 00:00:00","created_at":"2018-01-01 00:00:00"}`, string(bytes))

	book2 := Book{}
	err = json.Unmarshal(bytes, &book2)
	assert.Nil(t, err)
	assert.Equal(t, book1, book2)
}
