package jsontime

import (
	"testing"
	"time"
)

var json = ConfigWithCustomTimeFormat

type Book struct {
	Id        int        `json:"id"`
	PublishAt *time.Time `json:"publish_at" time_format:"sql_date" time_utc:"true"`
	UpdatedAt *time.Time `json:"updated_at" time_format:"sql_date" time_utc:"true"`
	CreatedAt time.Time  `json:"created_at" time_format:"sql_datetime" time_location:"UTC"`
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
	} else if string(bytes) != `{"id":1,"publish_at":null,"updated_at":"2018-01-01","created_at":"2018-01-01 00:00:00"}` {
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
		book.PublishAt != nil {
		t.Errorf("got:%v", book)
	}
}
