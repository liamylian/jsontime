# jsontime

A [json iterator](https://github.com/json-iterator/go) extension that support custom time format.

## Usage
100% compatibility with standard lib

Replace
```go
import "encoding/json"

json.Marshal(&data)
json.Unmarshal(input, &data)
```

with
```go
import "github.com/liamylian/jsontime"

var json = jsontime.ConfigWithCustomTimeFormat

json.Marshal(&data)
json.Unmarshal(input, &data)
```

## Example

```go
package main

import(
	"fmt"
	"time"
	"github.com/liamylian/jsontime"
)

var json = jsontime.ConfigWithCustomTimeFormat

type Book struct {
	Id        int           `json:"id"`
	UpdatedAt *time.Time    `json:"updated_at" time_format:"sql_date" time_utc:"true"`
	CreatedAt time.Time     `json:"created_at" time_format:"sql_datetime" time_location:"UTC"`
}

func main() {
	book := Book {
		Id:         1,
		UpdatedAt:  nil,
		CreatedAt:  time.Now(),
	}
	
	bytes, _ := json.Marshal(book)
	fmt.Printf("%s", bytes)
}

```