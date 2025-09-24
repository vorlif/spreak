package util

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

func ToNumber(n any) (float64, error) {
	n = Indirect(n)
	if n == nil {
		return 0, errors.New("number is nil")
	}

	switch nt := n.(type) {
	case uint:
		return float64(nt), nil
	case uint8:
		return float64(nt), nil
	case uint16:
		return float64(nt), nil
	case uint32:
		return float64(nt), nil
	case uint64:
		return float64(nt), nil
	case int:
		return float64(nt), nil
	case int8:
		return float64(nt), nil
	case int16:
		return float64(nt), nil
	case int32:
		return float64(nt), nil
	case int64:
		return float64(nt), nil
	case float32:
		return float64(nt), nil
	case float64:
		return nt, nil
	case complex64:
		return float64(real(nt)), nil
	case complex128:
		return real(nt), nil
	case string:
		res, err := strconv.ParseFloat(nt, 64)
		if err != nil {
			return 0, fmt.Errorf("unable to cast %#v of type %T to float64", n, n)
		}
		return res, nil
	}

	if num, errC := ToNumber(fmt.Sprintf("%d", n)); errC == nil {
		return num, nil
	}

	return ToNumber(fmt.Sprintf("%v", n))
}

// Indirect returns the value, after dereferencing as many times
// as necessary to reach the base type (or nil).
// From html/template/content.go
// Copyright 2011 The Go Authors. All rights reserved.
func Indirect(a any) any {
	if a == nil {
		return nil
	}
	if t := reflect.TypeOf(a); t.Kind() != reflect.Ptr {
		// Avoid creating a reflect.Value if it's not a pointer.
		return a
	}
	v := reflect.ValueOf(a)
	for v.Kind() == reflect.Ptr && !v.IsNil() {
		v = v.Elem()
	}
	return v.Interface()
}
