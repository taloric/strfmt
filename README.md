# strfmt
string formatter for golang

## usage
-------------
1. Format string with number index

    format like {0}{1}
```go
package main

import (
    "fmt"
    "github.com/taloric/strfmt"
)

func main(){
    format_string := "Today is a {0} {1}"
    res,err := strfmt.Format(format_string, "wonderful", "day")
    fmt.Println(res)
}
```

```
output: Today is a wonderful day
```

2. Format string with key

    format like {field}{desc}
```go
package main

import (
    "fmt"
    "github.com/taloric/strfmt"
)

func main(){
    format_string := "Today is a {what} {desc}"
    args := make(map[string]string)
    args["what"] = "wonderful"
    args["desc"] = "day"
    res,err := strfmt.FormatMap(format_string, &args)
    fmt.Println(res)
}
```

```
output: Today is a wonderful day
```

3. Format string with a struct

    same format as format a map
```go
package main

import (
    "fmt"
    "github.com/taloric/strfmt"
)

type Day struct{
    DayState string
}

func main() {
	format_string := "Today is a {DayState} day"
	args := &Day{DayState: "wonderful"}
	res, err := strfmt.FormatData(format_string, args)
	fmt.Println(res)
}
```

```
output: Today is a wonderful day
```

4. Format time with different format

    format like {0:2006/01/02 15:04:05}
```go
package main

import (
    "fmt"
    "github.com/taloric/strfmt"
)

type Day struct{
    DayState string
}

func main() {
	current_time := time.Now()

	format_time_normal := "Current Time is {0:2006-01-02 15:04:05 Mon}"
	format_time_short := "Current Time is {0:3:04PM}"
	format_time_map := "Current Time is {day:2006-01-02 15:04:05 Mon}"

    //use time.RFC1123Z to format a time string right now, maybe transform to another way not before long
	res, err := strfmt.Format(format_time_normal, current_time.Format(time.RFC1123Z))
	fmt.Println(res)

	res, err = strfmt.Format(format_time_short, current_time.Format(time.RFC1123Z))
	fmt.Println(res)

	time_map := make(map[string]string)
	time_map["day"] = current_time.Format(time.RFC1123Z)
	res, err = strfmt.FormatMap(format_time_map, &time_map)

	fmt.Println(res)
}
```

5. Format string and fill it with fixed length

    format like {0,20} will fill string length to 20 with space put on left

    format like {0,-10} will fill string length to 10 with space put on right

```go
package main

import (
    "fmt"
    "github.com/taloric/strfmt"
)

func main(){
    format_string := "Today is a {0,-20} {1,10}"
    res,err := strfmt.Format(format_string, "wonderful", "day")
    fmt.Println(res)
}
```

```
output: Today is a wonderful                   day
```