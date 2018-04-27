package json

import (
	"testing"
	"time"
)

type Book struct {
	Id        int       `json:"id"`
	UpdatedAt time.Time `json:"updated_at,sql_time"`
	CreatedAt time.Time `json:"created_at,sql_date"`
}

func TestMarshal(t *testing.T) {
	t2018 := time.Date(2018, 1, 1, 0, 0, 0, 0, time.Local)
	book := Book{
		Id:        1,
		UpdatedAt: t2018,
		CreatedAt: t2018,
	}

	if bytes, err := Marshal(book); err != nil {
		t.Error(err)
	} else if string(bytes) != `{"created_at":"2018-01-01","id":1,"updated_at":"2018-01-01 00:00:00"}` {
		t.Errorf("got:%s\n", bytes)
	}

}

func TestUnmarshal(t *testing.T) {
	t2018 := time.Date(2018, 1, 1, 0, 0, 0, 0, time.Local)
	bytes := []byte(`{"created_at":"2018-01-01","id":1,"updated_at":"2018-01-01 00:00:00"}`)

	book := Book{}
	if err := Unmarshal(bytes, &book); err != nil {
		t.Error(err)
	} else if book.Id != 1 || book.CreatedAt != t2018 || book.UpdatedAt != t2018 {
		t.Errorf("got:%s", book)
	}
}