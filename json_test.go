package jsontime

import (
	"testing"
	"time"
)

var json = ConfigWithCustomTimeFormat

type Book struct {
	Id          int        `json:"id"`
	PublishedAt *time.Time `json:"published_at" time_format:"sql_date" time_utc:"true"`
	UpdatedAt   *time.Time `json:"updated_at" time_format:"sql_date" time_utc:"true"`
	CreatedAt   time.Time  `json:"created_at" time_format:"sql_datetime" time_location:"UTC"`
}

func TestMarshalFormat(t *testing.T) {
	t2018 := time.Date(2018, 1, 1, 0, 0, 0, 0, time.UTC)
	book := Book{
		Id:        1,
		UpdatedAt: &t2018,
		CreatedAt: t2018,
	}

	if bytes, err := json.Marshal(book); err != nil {
		t.Error(err)
	} else if string(bytes) != `{"id":1,"published_at":null,"updated_at":"2018-01-01","created_at":"2018-01-01 00:00:00"}` {
		t.Errorf("got:%s\n", bytes)
	}

}

func TestUnmarshalFormat(t *testing.T) {
	t2018 := time.Date(2018, 1, 1, 0, 0, 0, 0, time.UTC)
	bytes := []byte(`{"id":1,"updated_at":"2018-01-01","created_at":"2018-01-01 00:00:00"}`)

	book := Book{}
	if err := json.Unmarshal(bytes, &book); err != nil {
		t.Error(err)
	} else if book.Id != 1 || book.CreatedAt != t2018 ||
		book.UpdatedAt == nil || *book.UpdatedAt != t2018 ||
		book.PublishedAt != nil {
		t.Errorf("got:%v", book)
	}
}

type User struct {
	Id        int        `json:"id"`
	UpdatedAt *time.Time `json:"updated_at" time_format:"sql_datetime" time_location:"Local"`
	CreatedAt time.Time  `json:"created_at" time_format:"sql_datetime" time_location:"Local"`
}

func TestLocale(t *testing.T) {
	user := User{
		Id:        0,
		UpdatedAt: nil,
		CreatedAt: time.Date(0, 1, 1, 0, 0, 0, 0, time.Local),
	}

	bytes, err := json.Marshal(user)
	if err != nil {
		t.Error(err.Error())
	}

	if string(bytes) != `{"id":0,"updated_at":null,"created_at":"0000-00-00 00:00:00"}` {
		t.Errorf("got: %s", bytes)
	}
}

func TestUnMarshalZero(t *testing.T) {
	user := User{}
	jsonBytes := []byte(`{"id":0,"updated_at":null,"created_at":"0000-00-00 00:00:00"}`)

	err := json.Unmarshal(jsonBytes, &user)
	if err != nil {
		t.Error(err.Error())
	}

	bytes, err := json.Marshal(user)
	if err != nil {
		t.Error(err.Error())
	}

	if string(bytes) != `{"id":0,"updated_at":null,"created_at":"0000-00-00 00:00:00"}` {
		t.Errorf("got: %s", bytes)
	}
}
