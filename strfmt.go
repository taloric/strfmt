package strfmt

import (
	"errors"
	"reflect"
	"strconv"
	"time"
)

const (
	INPUT_STR_ERROR           = "string [{0}] format is not available"
	INPUT_INDEX_OUT_OF_RANGE  = "string [{0}] format count did not match args length"
	INPUT_DATA_ERROR          = "args type [{0}] is not available, expect type is Struct"
	INPUT_DATA_KEY_NOT_EXISTS = "string [{0}] format could not found key [{1}] in args"
)

//handle unify error message
func format_error(errinfo string, formatter ...string) error {
	fmtResult, _ := format(errinfo, formatter...)
	return errors.New(fmtResult)
}

//format string by reflection structs
func format_data(str string, args interface{}) (string, error) {
	if len(str) == 0 || args == nil {
		return str, nil
	}
	args_type := reflect.TypeOf(args)
	args_value := reflect.ValueOf(args) //wrong

	kind := args_type.Kind()

	if kind == reflect.Ptr {
		args_type = args_type.Elem()
		args_value = args_value.Elem()
	}

	kind = args_type.Kind()
	if kind != reflect.Struct {
		return str, format_error(INPUT_DATA_ERROR, kind.String())
	}

	args_map := make(map[string]string)

	for i := 0; i < args_type.NumField(); i++ {
		name := args_type.Field(i).Name
		value_ref := args_value.Field(i)
		var value string

		t_field := value_ref.Type()
		if t_kind := t_field.Kind(); t_kind == reflect.Ptr {
			t_field = t_field.Elem()
		}

		switch t_field.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int32, reflect.Int64:
			value = strconv.FormatInt(value_ref.Int(), 10)
		case reflect.Float32, reflect.Float64:
			value = strconv.FormatFloat(value_ref.Float(), 'e', 2, 32)
		case reflect.Bool:
			value = strconv.FormatBool(value_ref.Bool())
		case reflect.Uint, reflect.Uint8, reflect.Uint32, reflect.Uint64:
			value = strconv.FormatUint(value_ref.Uint(), 10)
		default:
			{
				switch t_field.String() {
				case "time.Time":
					//todo : support self-def time format here
					value = value_ref.Interface().(time.Time).Format("2006-01-02 15:04:05")
				default:
					//todo : support recusively get fields here
					value = value_ref.String()
				}
			}
		}

		args_map[name] = value
	}
	return format_map(str, &args_map)
}

//format string key like {key} {key2}
func format_map(str string, args *map[string]string) (string, error) {
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

		//already handle {number1,number2 , should get } here
		if ch != '}' {
			return str, format_error(INPUT_STR_ERROR, str)
		}

		//if args map did not exists key, do nothing
		if _, ok := (*args)[string(key)]; !ok {
			return str, format_error(INPUT_DATA_KEY_NOT_EXISTS, str, string(key))
		}

		arg := (*args)[string(key)]

		pos++
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

//format numbers matching like {0} {1}
func format(str string, args ...string) (string, error) {
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

		//already handle {number1,number2 , should get } here
		if ch != '}' {
			return str, format_error(INPUT_STR_ERROR, str)
		}

		arg := []byte(args[index])
		pos++
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
