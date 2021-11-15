package strfmt

import (
	"fmt"
	"strconv"
	"testing"
	"time"
)

const (
	format_student = "[{Id}]-Student Info:name-[{Name}],gender-[{Gender}],age-[{Age}],birthday-[{BirthDay}],is graduate-[{IsGraduate}]"
	format_people  = "[{Id}]-People Info:name-[{Name}],gender-[{Gender}],age-[{Age}],birthday-[{BirthDay:2006/01/02}]"

	format_today          = "Today is a {0} day"
	format_today_rightpad = "Today is a {0,-20} day"
	format_today_leftpad  = "Today is a {0,20} day"
	format_today_info     = "Today is {NotMapTest1} {DayofWeek} {NotMap_Test} {NotMap}"
	format_date           = "Today is {year}-{month}-{day},current time is {hour}-{minute}-{seconds}"

	format_error_only_left_brace         = "Today is { Day"
	format_error_without_condition       = "Today is {} Day"
	format_error_only_right_brace        = "Today is } Day"
	format_error_without_complete_format = "Today is {0, Day"
	format_error_index_out_of_range      = "Today is {1} Day"

	format_time_normal = "Current Time is {0:2006-01-02 15:04:05 Mon}"
	format_time_short  = "Current Time is {0:3:04PM}"
	format_time_map    = "Current Time is {day:2006-01-02 15:04:05 Mon}"

	//maybe should support yyyy-mm format?(convert it to golang time format)
	format_time_error_format       = "Current Time is {0:yyyy-mm}"
	format_time_error_format_lost  = "Current Time is {0:}"
	format_time_error_format_wrong = "Current Time is {0:2003-02-02}"
)

type People struct {
	Id       string
	Name     string
	Gender   string
	Age      int
	BirthDay time.Time
}

type Student struct {
	Info       *People
	IsGraduate bool
}

var g_student *Student
var g_people *People

func TestMain(m *testing.M) {
	g_people = &People{
		Id:       strconv.FormatInt(time.Now().Unix(), 10),
		Name:     "Mary",
		Gender:   "Women",
		Age:      43,
		BirthDay: time.Date(1978, 1, 1, 1, 1, 1, 1, time.Local),
	}

	g_student = &Student{
		Info: &People{
			Id:       strconv.FormatInt(time.Now().Unix(), 10),
			Name:     "Donald",
			Gender:   "Man",
			Age:      51,
			BirthDay: time.Date(1970, 1, 1, 1, 1, 1, 1, time.Local),
		},
		IsGraduate: true,
	}
	fmt.Println("----Test Main Begin----")
	m.Run()
	fmt.Println("----Test Main End----")
}

func Test_Format(t *testing.T) {
	res, err := Format(format_today, "wonderful")
	if err != nil {
		t.Error("test_format throw error " + err.Error())
	}
	fmt.Println(res)

	res, err = Format(format_today, "bad")
	if err != nil {
		t.Error("Test_format throw error " + err.Error())
	}
	fmt.Println(res)

	res, err = Format(format_today_rightpad, "wonderful")
	if err != nil {
		t.Error("Test_format throw error " + err.Error())
	}
	fmt.Println(res)

	res, err = Format(format_today_leftpad, "wonderful")
	if err != nil {
		t.Error("Test_format throw error " + err.Error())
	}
	fmt.Println(res)
}

func Test_FormatMap(t *testing.T) {
	args := make(map[string]string)
	args["DayofWeek"] = time.Now().Weekday().String()
	res, err := FormatMap(format_today_info, &args)
	if err != nil {
		t.Error("Test_FormatMap throw error " + err.Error())
	}
	fmt.Println(res)
}

func Test_FormatMap_2(t *testing.T) {
	args := make(map[string]string)
	time_now := time.Now()
	args["year"] = strconv.FormatInt(int64(time_now.Year()), 10)
	args["month"] = strconv.FormatInt(int64(time_now.Month()), 10)
	args["day"] = strconv.FormatInt(int64(time_now.Day()), 10)
	args["hour"] = strconv.FormatInt(int64(time_now.Hour()), 10)
	args["minute"] = strconv.FormatInt(int64(time_now.Minute()), 10)
	args["seconds"] = strconv.FormatInt(int64(time_now.Second()), 10)

	res, err := FormatMap(format_date, &args)
	if err != nil {
		t.Error("Test_FormatMap_2 throw error " + err.Error())
	}
	fmt.Println(res)
}

func Test_FormatData(t *testing.T) {
	res, err := FormatData(format_people, g_people)
	if err != nil {
		t.Error("Test_FormatData throw error " + err.Error())
	}
	fmt.Println(res)
}

func Test_format_sub_data(t *testing.T) {
	res, err := FormatData(format_student, g_student)
	if err != nil {
		t.Error("Test_format_sub_data throw error " + err.Error())
	}
	fmt.Println(res)
}

func Test_format_error(t *testing.T) {

	fmt_data := "wonderful"
	_, err := Format(format_error_only_left_brace, fmt_data)
	if err == nil {
		t.Error("Test_format_error [format_error_only_left_brace] should throw error ")
		t.FailNow()
	}
	fmt.Println("Test_format_error [format_error_only_left_brace] throw error", err.Error())

	_, err = Format(format_error_without_condition, fmt_data)
	if err == nil {
		t.Error("Test_format_error [format_error_without_condition] should throw error ")
		t.FailNow()
	}
	fmt.Println("Test_format_error [format_error_without_condition] throw error", err.Error())

	_, err = Format(format_error_only_right_brace, fmt_data)
	if err == nil {
		t.Error("Test_format_error [format_error_only_right_brace] should throw error ")
		t.FailNow()
	}
	fmt.Println("Test_format_error [format_error_only_right_brace] throw error", err.Error())

	_, err = Format(format_error_without_complete_format, fmt_data)
	if err == nil {
		t.Error("Test_format_error [format_error_without_complete_format] should throw error ")
		t.FailNow()
	}
	fmt.Println("Test_format_error [format_error_without_complete_format] throw error", err.Error())

	_, err = Format(format_error_index_out_of_range, fmt_data)
	if err == nil {
		t.Error("Test_format_error [format_error_index_out_of_range] should throw error ")
		t.FailNow()
	}
	fmt.Println("Test_format_error [format_error_index_out_of_range] throw error", err.Error())
}

func Test_FormatTime(t *testing.T) {
	current_time := time.Now()
	res, err := Format(format_time_normal, current_time.Format(time.RFC1123Z))
	if err != nil {
		t.Error("Test_FormatTime throw error " + err.Error())
	}
	fmt.Println(res)

	res, err = Format(format_time_short, current_time.Format(time.RFC1123Z))
	if err != nil {
		t.Error("Test_FormatTime throw error " + err.Error())
	}
	fmt.Println(res)

	time_map := make(map[string]string)
	time_map["day"] = current_time.Format(time.RFC1123Z)
	res, err = FormatMap(format_time_map, &time_map)
	if err != nil {
		t.Error("Test_FormatTime throw error " + err.Error())
	}
	fmt.Println(res)
}

func Test_FormatTimeError(t *testing.T) {
	//todo: should be able to recognize any format of time
	current_time := time.Now()
	res, err := Format(format_time_error_format, current_time.Format(time.RFC1123Z))
	if err != nil {
		t.Error("Test_FormatTimeError throw error " + err.Error())
	}
	fmt.Println(res)

	res, err = Format(format_time_error_format_lost, current_time.Format(time.RFC1123Z))
	if err != nil {
		t.Error("Test_FormatTimeError throw error " + err.Error())
	}
	fmt.Println(res)

	res, err = Format(format_time_error_format_wrong, current_time.Format(time.RFC1123Z))
	if err != nil {
		t.Error("Test_FormatTimeError throw error " + err.Error())
	}
	fmt.Println(res)
}

//------------------------------------------//
//           Benchmark Test Below           //
//------------------------------------------//

func Benchmark_Format(b *testing.B) {
	_, err := Format(format_today, "wonderful")
	if err != nil {
		b.Error("test_format throw error " + err.Error())
		return
	}
	//fmt.Println(res)
}

func Benchmark_FormatMap(b *testing.B) {
	args := make(map[string]string)
	time_now := time.Now()
	args["year"] = strconv.FormatInt(int64(time_now.Year()), 10)
	args["month"] = strconv.FormatInt(int64(time_now.Month()), 10)
	args["day"] = strconv.FormatInt(int64(time_now.Day()), 10)
	args["hour"] = strconv.FormatInt(int64(time_now.Hour()), 10)
	args["minute"] = strconv.FormatInt(int64(time_now.Minute()), 10)
	args["seconds"] = strconv.FormatInt(int64(time_now.Second()), 10)

	_, err := FormatMap(format_date, &args)
	if err != nil {
		b.Error("Test_FormatMap_2 throw error " + err.Error())
	}
}

func Benchmark_FormatData(b *testing.B) {
	_, err := FormatData(format_people, g_people)
	if err != nil {
		b.Error("Test_FormatData throw error " + err.Error())
	}
}
