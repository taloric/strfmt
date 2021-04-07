package strfmt

import (
	"regexp"
)

var (
	regex_map_nums    *regexp.Regexp = regexp.MustCompile("{[0-9]*?}")
	regex_map_formats *regexp.Regexp = regexp.MustCompile("{[a-zA-Z0-9]*?}")
)

//format string key like {key} {key2}
func format_map(input string, data map[string]string) string {
	if len(input) == 0 || len(data) == 0 {
		return input
	}
	return input
}

//format numbers matching like {0} {1}
func format(input string, datas ...string) string {
	if len(input) == 0 || len(datas) == 0 {
		return input
	}
	return input
}

//format string by reflection structs
func format_data(input string, data *interface{}) string {
	if len(input) == 0 || data == nil {
		return input
	}
	return input
}
