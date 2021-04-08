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
func FormatError(errinfo string, formatter ...string) error {
	fmtResult, _ := Format(errinfo, formatter...)
	return errors.New(fmtResult)
}

//not complete
func get_reflect_data(t *reflect.Type, v *reflect.Value) map[string]string {
	kind := (*t).Kind()
	typ := *t
	val := *v
	args_map := make(map[string]string)
	if kind == reflect.Ptr {
		typ = (*t).Elem()
		val = (*v).Elem()
	}
	kind = typ.Kind()
	if kind != reflect.Struct {
		return args_map
	}
	for i := 0; i < typ.NumField(); i++ {
		name := typ.Field(i).Name
		value_ref := val.Field(i)
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
		case reflect.Struct:
			resmap := get_reflect_data(&t_field, &value_ref)
			for k, v := range resmap {
				args_map[k] = v
			}
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
	return args_map
}

//format string by reflection structs
func FormatData(str string, args interface{}) (string, error) {
	if len(str) == 0 || args == nil {
		return str, nil
	}
	args_type := reflect.TypeOf(args)
	args_value := reflect.ValueOf(args)
	args_map := get_reflect_data(&args_type, &args_value)
	return FormatMap(str, &args_map)
}

//format string key like {key} {key2}
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
					return str, FormatError(INPUT_STR_ERROR, str)
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
			return str, FormatError(INPUT_STR_ERROR, str)
		}

		if ch = str[pos]; ch < 'A' || (ch > 'Z' && (ch < 'a' || ch > 'z')) {
			return str, FormatError(INPUT_STR_ERROR, str)
		}

		// get numbers in {}
		var key []byte
		for { //(ch >= '0' && ch <= '9' && index < 1000000)
			key = append(key, ch)
			pos++
			if pos == length {
				return str, FormatError(INPUT_STR_ERROR, str)
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
				return str, FormatError(INPUT_STR_ERROR, str)
			}

			ch = str[pos]
			if ch == '-' {
				leftJustify = true
				pos++
				if pos == length {
					return str, FormatError(INPUT_STR_ERROR, str)
				}
				ch = str[pos]
			}

			if ch < '0' || ch > '9' {
				return str, FormatError(INPUT_STR_ERROR, str)
			}

			//get numbers after ','
			for {
				width = width*10 + ch - '0'
				pos++
				if pos == length {
					return str, FormatError(INPUT_STR_ERROR, str)
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
			return str, FormatError(INPUT_STR_ERROR, str)
		}

		//if args map did not exists key, do nothing
		if _, ok := (*args)[string(key)]; !ok {
			return str, FormatError(INPUT_DATA_KEY_NOT_EXISTS, str, string(key))
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
					return str, FormatError(INPUT_STR_ERROR, str)
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
			return str, FormatError(INPUT_STR_ERROR, str)
		}

		if ch = str[pos]; ch < '0' || ch > '9' {
			return str, FormatError(INPUT_STR_ERROR, str)
		}

		// get numbers in {}
		var index byte = 0
		for { //(ch >= '0' && ch <= '9' && index < 1000000)
			index = index*10 + ch - '0'
			pos++
			if pos == length {
				return str, FormatError(INPUT_STR_ERROR, str)
			}
			ch = str[pos]

			//escape loop
			if ch < '0' || ch > '9' || index >= 255 {
				break
			}
		}

		if int(index) >= len(args) {
			return str, FormatError(INPUT_INDEX_OUT_OF_RANGE, str)
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
				return str, FormatError(INPUT_STR_ERROR, str)
			}

			ch = str[pos]
			if ch == '-' {
				leftJustify = true
				pos++
				if pos == length {
					return str, FormatError(INPUT_STR_ERROR, str)
				}
				ch = str[pos]
			}

			if ch < '0' || ch > '9' {
				return str, FormatError(INPUT_STR_ERROR, str)
			}

			//get numbers after ','
			for {
				//don't convert ch to int here , for space quantity limitation
				width = width*10 + ch - '0'
				pos++
				if pos == length {
					return str, FormatError(INPUT_STR_ERROR, str)
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
			return str, FormatError(INPUT_STR_ERROR, str)
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
