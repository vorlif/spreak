package cast

import (
	"errors"
	"fmt"
	"time"
)

func ToTime(i any) (t time.Time, err error) {
	i = Indirect(i)
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
