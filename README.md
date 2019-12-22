# jsontime

![](https://github.com/liamylian/jsontime/workflows/Test/badge.svg)

A [json iterator](https://github.com/json-iterator/go) extension that support custom time format.

# Install

`go get github.com/liamylian/jsontime/v2`

or

`go mod edit -require=github.com/liamylian/jsontime/v2@v2.0.0`


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
import "github.com/liamylian/jsontime/v2"

var json = jsontime.ConfigWithCustomTimeFormat

json.Marshal(&data)
json.Unmarshal(input, &data)
```

## Example

```go
package main

import (
	"fmt"
	"time"

	"github.com/liamylian/jsontime/v2"
)

var json = jsontime.ConfigWithCustomTimeFormat

func init() {
	jsontime.SetDefaultTimeFormat(time.RFC3339, time.Local)
}

type Book struct {
	Id          int        `json:"id"`
	PublishedAt time.Time  `json:"published_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
	CreatedAt   time.Time  `json:"created_at"`
}

func main() {
	book := Book{
		Id:          1,
		PublishedAt: time.Now(),
		UpdatedAt:   nil,
		CreatedAt:   time.Now(),
	}

	bytes, _ := json.Marshal(book)
	fmt.Printf("%s", bytes)
}
```

## Advanced Usage

```go
var json = jsontime.ConfigWithCustomTimeFormat
timeZoneShanghai, _ := time.LoadLocation("Asia/Shanghai")
jsontime.AddTimeFormatAlias("sql_datetime", "2006-01-02 15:04:05")
jsontime.AddLocaleAlias("shanghai", timeZoneShanghai)

type Book struct {
	Id          int        `json:"id"`
	PublishedAt time.Time  `json:"published_at" time_format:"sql_datetime" time_location:"shanghai"`
	UpdatedAt   *time.Time `json:"updated_at" time_format:"sql_datetime" time_location:"shanghai"`
	CreatedAt   time.Time  `json:"created_at" time_format:"sql_datetime" time_location:"shanghai"`
}

bytes, _ := json.Marshal(book)
```

## Notice

Because there are two maps in the package, that is not thread safe. So it is suggested to call `AddTimeFormatAlias`, `AddLocaleAlias` only once, right before calling `Marshal` and `UnMarshal`. Or you should make sure not to call them at the same time (writing and reading map at the same time will panic).
