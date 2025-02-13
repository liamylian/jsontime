# jsontime

![](https://github.com/yatzek/jsontime/workflows/Test/badge.svg)

A [json iterator](https://github.com/json-iterator/go) extension that support custom time format.

# Install

`go get github.com/yatzek/jsontime/v2`

or

`go mod edit -require=github.com/yatzek/jsontime/v2@v2.0.0`


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
import (
	jsoniter "github.com/json-iterator/go"
	jsontime "github.com/yatzek/jsontime/v2"
)

json := jsoniter.ConfigCompatibleWithStandardLibrary
timeExtension := jsontime.NewCustomTimeExtension()
json.RegisterExtension(timeExtension)

json.Marshal(&data)
json.Unmarshal(input, &data)
```

## Example

```go
package main

import (
	"fmt"
	"time"
	jsoniter "github.com/json-iterator/go"
	jsontime "github.com/yatzek/jsontime/v2"
)

type Book struct {
	Id          int        `json:"id"`
	PublishedAt time.Time  `json:"published_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
	CreatedAt   time.Time  `json:"created_at"`
}

func main() {
	json := jsoniter.ConfigCompatibleWithStandardLibrary
	timeExtension := jsontime.NewCustomTimeExtension()
	timeExtension.SetDefaultTimeFormat(time.RFC3339, time.Local)
	json.RegisterExtension(timeExtension)

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
json := jsoniter.ConfigCompatibleWithStandardLibrary
timeExtension := jsontime.NewCustomTimeExtension()

timeZoneShanghai, _ := time.LoadLocation("Asia/Shanghai")
timeExtension.AddTimeFormatAlias("sql_datetime", "2006-01-02 15:04:05")
timeExtension.AddLocaleAlias("shanghai", timeZoneShanghai)
json.RegisterExtension(timeExtension)

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

