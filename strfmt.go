package strfmt

import (
	"errors"
)

func format_error(str string) error {
	return errors.New("Format String [" + str + "] Is Not Correct Format")
}

//format string by reflection structs
func format_data(str string, args *interface{}) (string, error) {
	if len(str) == 0 || args == nil {
		return str, nil
	}
	return str, nil
}

//format string key like {key} {key2}
func format_map(str string, args map[string]string) (string, error) {
	if len(str) == 0 || len(args) == 0 {
		return str, nil
	}
	return str, nil
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
					return str, format_error(str)
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
			return str, format_error(str)
		}

		if ch = str[pos]; ch < '0' || ch > '9' {
			return str, format_error(str)
		}

		// get numbers in {}
		var index byte = 0
		for { //(ch >= '0' && ch <= '9' && index < 1000000)
			index = index*10 + ch - '0'
			pos++
			if pos == length {
				return str, format_error(str)
			}
			ch = str[pos]

			//escape loop
			if ch < '0' || ch > '9' || index >= 255 {
				break
			}
		}

		if int(index) >= len(args) {
			return str, format_error(str)
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
				return str, format_error(str)
			}

			ch = str[pos]
			if ch == '-' {
				leftJustify = true
				pos++
				if pos == length {
					return str, format_error(str)
				}
				ch = str[pos]
			}

			if ch < '0' || ch > '9' {
				return str, format_error(str)
			}

			for {
				width = width*10 + ch - '0'
				pos++
				if pos == length {
					return str, format_error(str)
				}
				ch = str[pos]

				if ch < '0' || ch > '9' || width >= 255 {
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

		arg := []byte(args[index])

		if ch != '}' {
			return str, format_error(str)
		}

		pos++

		pad := int(width) - len(arg)

		if !leftJustify && pad > 0 {
			for j := 0; j <= pad; j++ {
				result = append(result, ' ')
			}
		}
		result = append(result, arg...)

		if leftJustify && pad > 0 {
			for j := 0; j <= pad; j++ {
				result = append(result, ' ')
			}
		}

	}

	return string(result), nil
}
