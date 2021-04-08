package strfmt

import (
	"errors"
	"reflect"
	"strconv"
	"time"
)

//error message
const (
	INPUT_STR_ERROR           = "string [{0}] format is not available"
	INPUT_INDEX_OUT_OF_RANGE  = "string [{0}] format count did not match args length"
	INPUT_DATA_ERROR          = "args type [{0}] is not available, expect type is Struct"
	INPUT_DATA_KEY_NOT_EXISTS = "string [{0}] format could not found key [{1}] in args"
	INPUT_TIME_FORMAT_ERROR   = "time format [{0}] is not available"
)

//handle unify error message
func format_error(errinfo string, formatter ...string) error {
	fmtResult, _ := Format(errinfo, formatter...)
	return errors.New(fmtResult)
}

//get sub struct data
func get_reflect_data(t *reflect.Type, v *reflect.Value) map[string]string {

	typ := *t
	val := *v
	kind := typ.Kind()
	args_map := make(map[string]string)

	if kind == reflect.Ptr {
		//get real type in Pointer
		typ = typ.Elem()
		val = val.Elem()
		kind = typ.Kind()
	}

	if kind != reflect.Struct {
		return args_map
	}

	for i := 0; i < typ.NumField(); i++ {
		name := typ.Field(i).Name
		field_value := val.Field(i)
		var value string

		field_type := field_value.Type()
		field_kind := field_type.Kind()

		if field_kind == reflect.Ptr {
			//get field type in Pointer
			field_type = field_type.Elem()
			field_kind = field_type.Kind()
			field_value = field_value.Elem()
		}

		if field_type.String() == "time.Time" {
			//format time first
			//convert to RFC1123Z as a default format(which will not lost too much information)
			value = field_value.Interface().(time.Time).Format(time.RFC1123Z)
		} else {
			switch field_kind {
			case reflect.Int, reflect.Int8, reflect.Int32, reflect.Int64:
				//decimalism for int convert
				value = strconv.FormatInt(field_value.Int(), 10)
			case reflect.Float32, reflect.Float64:
				//point = 2,bit size = 32
				value = strconv.FormatFloat(field_value.Float(), 'e', 2, 32)
			case reflect.Bool:
				value = strconv.FormatBool(field_value.Bool())
			case reflect.Uint, reflect.Uint8, reflect.Uint32, reflect.Uint64:
				value = strconv.FormatUint(field_value.Uint(), 10)
			case reflect.Struct:
				//recusively get struct data here
				resmap := get_reflect_data(&field_type, &field_value)
				for k, v := range resmap {
					args_map[k] = v
				}
			default: //normal type
				value = field_value.String()
			}
		}

		args_map[name] = value
	}
	return args_map
}

//Format Strings with struct type data
//	str:target string, args:struct
//	if args is nil or len(str) is zero, return itself
//	string format should be like : some description{field}
func FormatData(str string, args interface{}) (string, error) {
	if len(str) == 0 || args == nil {
		return str, nil
	}
	args_type := reflect.TypeOf(args)
	args_value := reflect.ValueOf(args)

	args_map := get_reflect_data(&args_type, &args_value)
	return FormatMap(str, &args_map)
}

//Format Strings with a map[string]string
//	str:target string, args:map
//	if args is nil or len(str) is zero, return itself
//	string format should be like : some description{field}
func FormatMap(str string, args *map[string]string) (string, error) {
	if len(str) == 0 || len(*args) == 0 {
		return str, nil
	}
	var result []byte
	pos := 0
	length := len(str)
	var ch byte

	for {
		for { //(pos<len)

			//escape loop
			if pos >= length {
				break
			}

			ch = str[pos]
			pos++

			if ch == '}' {
				//escape char for }}
				if pos < length && str[pos] == '}' {
					pos++
				} else {
					return str, format_error(INPUT_STR_ERROR, str)
				}
			}

			//escape char for {{
			if ch == '{' {
				if pos < length && str[pos] == '{' {
					pos++
				} else {
					pos--
					break
				}
			}

			result = append(result, ch)
		}

		if pos == length {
			break
		}
		pos++
		if pos == length {
			return str, format_error(INPUT_STR_ERROR, str)
		}

		if ch = str[pos]; ch < 'A' || (ch > 'Z' && (ch < 'a' || ch > 'z')) {
			return str, format_error(INPUT_STR_ERROR, str)
		}

		// get numbers in {}
		var key []byte
		for { //(ch >= '0' && ch <= '9' && index < 1000000)
			key = append(key, ch)
			pos++
			if pos == length {
				return str, format_error(INPUT_STR_ERROR, str)
			}
			ch = str[pos]

			//escape loop
			if ch < 'A' || (ch > 'Z' && (ch < 'a' || ch > 'z')) {
				break
			}
		}

		//remove all space
		for { //ch = str[pos]) | (pos < length && ch == ' ')
			if ch = str[pos]; pos >= length || ch != ' ' {
				break
			}
			pos++
		}

		leftJustify := false
		var width byte = 0

		//get number after ',' to leftpad or rightpad space
		if ch == ',' {
			pos++
			for { //(pos < len && str[pos] == ' ')
				if pos >= length || str[pos] != ' ' {
					break
				}
				pos++
			}

			if pos == length {
				return str, format_error(INPUT_STR_ERROR, str)
			}

			ch = str[pos]
			if ch == '-' {
				leftJustify = true
				pos++
				if pos == length {
					return str, format_error(INPUT_STR_ERROR, str)
				}
				ch = str[pos]
			}

			if ch < '0' || ch > '9' {
				return str, format_error(INPUT_STR_ERROR, str)
			}

			//get numbers after ','
			for {
				width = width*10 + ch - '0'
				pos++
				if pos == length {
					return str, format_error(INPUT_STR_ERROR, str)
				}
				ch = str[pos]

				if ch < '0' || ch > '9' || width > 255 {
					break
				}
			}
		}

		for { //ch = str[pos]) | (pos < length && ch == ' ')
			if ch = str[pos]; pos >= length || ch != ' ' {
				break
			}
			pos++
		}

		//get time format after :
		var timeFormatter []byte
		if ch == ':' {
			pos++
			for {
				if pos == length {
					return str, format_error(INPUT_STR_ERROR, str)
				}

				ch = str[pos]
				pos++

				if ch == '{' {
					//escape char for {{
					if pos < length && str[pos] == '{' {
						pos++
					} else {
						return str, format_error(INPUT_STR_ERROR, str)
					}
				}

				//escape char for }}
				if ch == '}' {
					if pos < length && str[pos] == '}' {
						pos++
					} else {
						pos--
						break
					}
				}

				timeFormatter = append(timeFormatter, ch)
			}
		}

		//already handle {number1,number2 , should get } here
		if ch != '}' {
			return str, format_error(INPUT_STR_ERROR, str)
		}
		pos++

		//if args map did not exists key
		if _, ok := (*args)[string(key)]; !ok {
			return str, format_error(INPUT_DATA_KEY_NOT_EXISTS, str, string(key))
		}

		arg := []byte((*args)[string(key)])

		if len(timeFormatter) > 0 {
			t_arg, err := time.Parse(time.RFC1123Z, string(arg))
			if err == nil {
				arg = []byte(t_arg.Format(string(timeFormatter)))
				//todo : should validate time format if is illegal
			} else {
				return str, format_error(INPUT_TIME_FORMAT_ERROR, string(timeFormatter))
			}
		}

		pad := int(width) - len(arg)

		//leftPad
		if !leftJustify && pad > 0 {
			for j := 0; j <= pad; j++ {
				result = append(result, ' ')
			}
		}

		//append arg
		result = append(result, arg...)

		//rightPad
		if leftJustify && pad > 0 {
			for j := 0; j <= pad; j++ {
				result = append(result, ' ')
			}
		}
	}

	return string(result), nil
}

//Format Strings with string args
//	str:target string, args: strings
//	if args is nil or len(str) is zero, return itself
//	string format should be like : some description{0}{1}
func Format(str string, args ...string) (string, error) {
	if len(str) == 0 || len(args) == 0 {
		return str, nil
	}

	var result []byte
	pos := 0
	length := len(str)
	var ch byte

	for {
		for { //(pos<len)

			//escape loop
			if pos >= length {
				break
			}

			ch = str[pos]
			pos++

			if ch == '}' {
				//escape char for }}
				if pos < length && str[pos] == '}' {
					pos++
				} else {
					return str, format_error(INPUT_STR_ERROR, str)
				}
			}

			//escape char for {{
			if ch == '{' {
				if pos < length && str[pos] == '{' {
					pos++
				} else {
					pos--
					break
				}
			}

			result = append(result, ch)
		}

		if pos == length {
			break
		}
		pos++
		if pos == length {
			return str, format_error(INPUT_STR_ERROR, str)
		}

		if ch = str[pos]; ch < '0' || ch > '9' {
			return str, format_error(INPUT_STR_ERROR, str)
		}

		// get numbers in {}
		var index byte = 0
		for { //(ch >= '0' && ch <= '9' && index < 1000000)
			index = index*10 + ch - '0'
			pos++
			if pos == length {
				return str, format_error(INPUT_STR_ERROR, str)
			}
			ch = str[pos]

			//escape loop
			if ch < '0' || ch > '9' || index >= 255 {
				break
			}
		}

		if int(index) >= len(args) {
			return str, format_error(INPUT_INDEX_OUT_OF_RANGE, str)
		}

		//remove all space
		for { //ch = str[pos]) | (pos < length && ch == ' ')
			if ch = str[pos]; pos >= length || ch != ' ' {
				break
			}
			pos++
		}

		leftJustify := false
		var width byte = 0

		//get number after ',' to leftpad or rightpad space
		if ch == ',' {
			pos++
			for { //(pos < len && str[pos] == ' ')
				if pos >= length || str[pos] != ' ' {
					break
				}
				pos++
			}

			if pos == length {
				return str, format_error(INPUT_STR_ERROR, str)
			}

			ch = str[pos]
			if ch == '-' {
				leftJustify = true
				pos++
				if pos == length {
					return str, format_error(INPUT_STR_ERROR, str)
				}
				ch = str[pos]
			}

			if ch < '0' || ch > '9' {
				return str, format_error(INPUT_STR_ERROR, str)
			}

			//get numbers after ','
			for {
				//don't convert ch to int here , for space quantity limitation
				width = width*10 + ch - '0'
				pos++
				if pos == length {
					return str, format_error(INPUT_STR_ERROR, str)
				}
				ch = str[pos]

				//support most 255 space only
				if ch < '0' || ch > '9' || width > 255 {
					break
				}
			}
		}

		for { //ch = str[pos]) | (pos < length && ch == ' ')
			if ch = str[pos]; pos >= length || ch != ' ' {
				break
			}
			pos++
		}

		//get time format after :
		var timeFormatter []byte
		if ch == ':' {
			pos++
			for {
				if pos == length {
					return str, format_error(INPUT_STR_ERROR, str)
				}

				ch = str[pos]
				pos++

				if ch == '{' {
					//escape char for {{
					if pos < length && str[pos] == '{' {
						pos++
					} else {
						return str, format_error(INPUT_STR_ERROR, str)
					}
				}

				//escape char for }}
				if ch == '}' {
					if pos < length && str[pos] == '}' {
						pos++
					} else {
						pos--
						break
					}
				}

				timeFormatter = append(timeFormatter, ch)
			}
		}

		//already handle {number1,number2 , should get } here
		if ch != '}' {
			return str, format_error(INPUT_STR_ERROR, str)
		}
		pos++

		//format time here
		arg := []byte(args[index])

		if len(timeFormatter) > 0 {
			t_arg, err := time.Parse(time.RFC1123Z, string(arg))
			if err == nil {
				arg = []byte(t_arg.Format(string(timeFormatter)))
			} else {
				return str, format_error(INPUT_TIME_FORMAT_ERROR, string(timeFormatter))
			}
		}

		pad := int(width) - len(arg)

		//leftPad
		if !leftJustify && pad > 0 {
			for j := 0; j <= pad; j++ {
				result = append(result, ' ')
			}
		}

		//append arg
		result = append(result, arg...)

		//rightPad
		if leftJustify && pad > 0 {
			for j := 0; j <= pad; j++ {
				result = append(result, ' ')
			}
		}
	}

	return string(result), nil
}
