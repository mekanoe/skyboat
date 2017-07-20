package util

import (
	"errors"
	"os"
	"strconv"
	"strings"
)

type EnvValue struct {
	Empty bool

	value        string
	defaultValue interface{}
}

var (
	ErrCannotConvert = errors.New("cannot convert type")
)

func (e EnvValue) String() (string, error) {
	if e.Empty {
		return e.defaultValue.(string), nil
	}
	return e.value, nil
}

func (e EnvValue) Bytes() ([]byte, error) {
	if e.Empty {
		switch e.defaultValue.(type) {
		case []byte:
			return e.defaultValue.([]byte), nil
		case string:
			return []byte(e.defaultValue.(string)), nil
		default:
			return []byte(""), ErrCannotConvert
		}
	}

	return []byte(e.value), nil
}

func (e EnvValue) StringSlice() ([]string, error) {
	if e.Empty {
		return e.defaultValue.([]string), nil
	}
	return strings.Split(e.value, ","), nil
}

func (e EnvValue) Int() (int, error) {
	if e.Empty {
		switch e.defaultValue.(type) {
		case string:
			return strconv.Atoi(e.defaultValue.(string))
		default:
			return e.defaultValue.(int), nil
		}
	}

	return strconv.Atoi(e.value)
}

func (e EnvValue) Bool() (bool, error) {
	if e.Empty {
		switch e.defaultValue.(type) {
		case string:
			return e.defaultValue.(string) != "0", nil
		case int:
			return e.defaultValue.(int) != 0, nil
		case bool:
			return e.defaultValue.(bool), nil
		default:
			return false, ErrCannotConvert
		}
	}

	return e.value != "0", nil
}

// Getenvdef gets an environment variable, then returns a conversion interface
func Getenvdef(key string, val interface{}) EnvValue {
	out := os.Getenv(key)
	ev := EnvValue{
		defaultValue: val,
		value:        out,
		Empty:        false,
	}

	if out == "" {
		ev.Empty = true
	}

	return ev
}
