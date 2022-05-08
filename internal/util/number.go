package util

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

func ToNumber(n interface{}) (res float64, err error) {
	switch nt := n.(type) {
	case uint:
		res = float64(nt)
		return
	case uint8:
		res = float64(nt)
		return
	case uint16:
		res = float64(nt)
		return
	case uint32:
		res = float64(nt)
		return
	case uint64:
		res = float64(nt)
		return
	case int:
		res = float64(nt)
		return
	case int8:
		res = float64(nt)
		return
	case int16:
		res = float64(nt)
		return
	case int32:
		res = float64(nt)
		return
	case int64:
		res = float64(nt)
		return
	case float32:
		res = float64(nt)
		return
	case float64:
		res = nt
		return
	case string:
		res, err = strconv.ParseFloat(nt, 64)
		return
	}

	if IsNilInterface(n) {
		return 0, errors.New("number is nil interface")
	}

	switch nt := n.(type) {
	case *uint:
		res = float64(*nt)
		return
	case *uint8:
		res = float64(*nt)
		return
	case *uint16:
		res = float64(*nt)
		return
	case *uint32:
		res = float64(*nt)
		return
	case *uint64:
		res = float64(*nt)
		return
	case *int:
		res = float64(*nt)
		return
	case *int8:
		res = float64(*nt)
		return
	case *int16:
		res = float64(*nt)
		return
	case *int32:
		res = float64(*nt)
		return
	case *int64:
		res = float64(*nt)
		return
	case *float32:
		res = float64(*nt)
		return
	case *float64:
		res = *nt
		return
	case *string:
		return ToNumber(*nt)
	}

	return ToNumber(fmt.Sprintf("%d", n))
}

func IsNilInterface(i interface{}) bool {
	return i == nil || (reflect.ValueOf(i).Kind() == reflect.Ptr && reflect.ValueOf(i).IsNil())
}
