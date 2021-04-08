package strfmt

import (
	"fmt"
	"strconv"
	"testing"
	"time"
)

const (
	format_student = "[{Id}]-Student Info:name-[{Name}],gender-[{Gender}],age-[{Age}],birthday-[{BirthDay}],is graduate-[{IsGraduate}]"
	format_people  = "[{Id}]-People Info:name-[{Name}],gender-[{Gender}],age-[{Age}],birthday-[{BirthDay}]"

	format_today          = "Today is a {0} day"
	format_today_rightpad = "Today is a {0,-20} day"
	format_today_leftpad  = "Today is a {0,20} day"
	format_today_info     = "Today is {DayofWeek}"
	format_date           = "Today is {year}-{month}-{day},current time is {hour}-{minute}-{seconds}"

	format_error_only_left_brace         = "Today is { Day"
	format_error_without_condition       = "Today is {} Day"
	format_error_only_right_brace        = "Today is } Day"
	format_error_without_complete_format = "Today is {0, Day"
	format_error_index_out_of_range      = "Today is {1} Day"
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

func Test_format(t *testing.T) {
	res, err := format(format_today, "wonderful")
	if err != nil {
		t.Error("test_format throw error " + err.Error())
	}
	fmt.Println(res)

	res, err = format(format_today, "bad")
	if err != nil {
		t.Error("Test_format throw error " + err.Error())
	}
	fmt.Println(res)

	res, err = format(format_today_rightpad, "wonderful")
	if err != nil {
		t.Error("Test_format throw error " + err.Error())
	}
	fmt.Println(res)

	res, err = format(format_today_leftpad, "wonderful")
	if err != nil {
		t.Error("Test_format throw error " + err.Error())
	}
	fmt.Println(res)
}

func Test_format_map(t *testing.T) {
	args := make(map[string]string)
	args["DayofWeek"] = time.Now().Weekday().String()
	res, err := format_map(format_today_info, &args)
	if err != nil {
		t.Error("Test_format_map throw error " + err.Error())
	}
	fmt.Println(res)
}

func Test_format_map_2(t *testing.T) {
	args := make(map[string]string)
	time_now := time.Now()
	args["year"] = strconv.FormatInt(int64(time_now.Year()), 10)
	args["month"] = strconv.FormatInt(int64(time_now.Month()), 10)
	args["day"] = strconv.FormatInt(int64(time_now.Day()), 10)
	args["hour"] = strconv.FormatInt(int64(time_now.Hour()), 10)
	args["minute"] = strconv.FormatInt(int64(time_now.Minute()), 10)
	args["seconds"] = strconv.FormatInt(int64(time_now.Second()), 10)

	res, err := format_map(format_date, &args)
	if err != nil {
		t.Error("Test_format_map_2 throw error " + err.Error())
	}
	fmt.Println(res)
}

func Test_format_data(t *testing.T) {
	res, err := format_data(format_people, g_people)
	if err != nil {
		t.Error("Test_format_data throw error " + err.Error())
	}
	fmt.Println(res)
}

func Test_format_sub_data(t *testing.T) {
	res, err := format_data(format_student, g_student)
	if err != nil {
		t.Error("Test_format_sub_data throw error " + err.Error())
	}
	fmt.Println(res)
}

func Test_format_error(t *testing.T) {

	fmt_data := "wonderful"
	_, err := format(format_error_only_left_brace, fmt_data)
	if err == nil {
		t.Error("Test_format_error [format_error_only_left_brace] should throw error ")
	}
	fmt.Println("Test_format_error [format_error_only_left_brace] throw error", err.Error())

	_, err = format(format_error_without_condition, fmt_data)
	if err == nil {
		t.Error("Test_format_error [format_error_without_condition] should throw error ")
	}
	fmt.Println("Test_format_error [format_error_without_condition] throw error", err.Error())

	_, err = format(format_error_only_right_brace, fmt_data)
	if err == nil {
		t.Error("Test_format_error [format_error_only_right_brace] should throw error ")
	}
	fmt.Println("Test_format_error [format_error_only_right_brace] throw error", err.Error())

	_, err = format(format_error_without_complete_format, fmt_data)
	if err == nil {
		t.Error("Test_format_error [format_error_without_complete_format] should throw error ")
	}
	fmt.Println("Test_format_error [format_error_without_complete_format] throw error", err.Error())

	_, err = format(format_error_index_out_of_range, fmt_data)
	if err == nil {
		t.Error("Test_format_error [format_error_index_out_of_range] should throw error ")
	}
	fmt.Println("Test_format_error [format_error_index_out_of_range] throw error", err.Error())
}

//------------------------------------------//
//           Benchmark Test Below           //
//------------------------------------------//

func Benchmark_format(b *testing.B) {
	_, err := format(format_today, "wonderful")
	if err != nil {
		b.Error("test_format throw error " + err.Error())
		return
	}
	//fmt.Println(res)
}

func Benchmark_format_map(b *testing.B) {
	args := make(map[string]string)
	//Today is {year}-{month}-{day},current time is {hour}-{minute}-{seconds}
	time_now := time.Now()
	args["year"] = strconv.FormatInt(int64(time_now.Year()), 10)
	args["month"] = strconv.FormatInt(int64(time_now.Month()), 10)
	args["day"] = strconv.FormatInt(int64(time_now.Day()), 10)
	args["hour"] = strconv.FormatInt(int64(time_now.Hour()), 10)
	args["minute"] = strconv.FormatInt(int64(time_now.Minute()), 10)
	args["seconds"] = strconv.FormatInt(int64(time_now.Second()), 10)

	_, err := format_map(format_date, &args)
	if err != nil {
		b.Error("Test_format_map_2 throw error " + err.Error())
	}
	//fmt.Println(res)
}

func Benchmark_format_data(b *testing.B) {
	_, err := format_data(format_people, g_people)
	if err != nil {
		b.Error("Test_format_data throw error " + err.Error())
	}
	//fmt.Println(res)
}
