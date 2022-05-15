package util

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"time"
)

func ToTime(i interface{}) (t time.Time, err error) {
	i = indirect(i)
	if i == nil {
		return time.Time{}, errors.New("time is nil")
	}

	switch v := i.(type) {
	case time.Time:
		return v, nil
	case time.Duration:
		return time.Now().Add(v), nil
	}

	if num, err := ToNumber(i); err == nil {
		return time.Unix(int64(num), 0), nil
	}
	return time.Time{}, fmt.Errorf("unable to cast %#v of type %T to Time", i, i)
}

func ToNumber(n interface{}) (res float64, err error) {
	n = indirect(n)
	if n == nil {
		return 0, errors.New("number is nil")
	}

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
		if err != nil {
			return 0, fmt.Errorf("unable to cast %#v of type %T to float64", n, n)
		}
		return
	}

	if num, errC := ToNumber(fmt.Sprintf("%d", n)); errC == nil {
		return num, nil
	}

	return ToNumber(fmt.Sprintf("%v", n))
}

// From html/template/content.go
// Copyright 2011 The Go Authors. All rights reserved.
// indirect returns the value, after dereferencing as many times
// as necessary to reach the base type (or nil).
func indirect(a interface{}) interface{} {
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
