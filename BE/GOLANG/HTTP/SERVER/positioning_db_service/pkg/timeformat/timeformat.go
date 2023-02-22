package timeformat

import (
	"errors"
	"fmt"
	"time"
)

const timeFormat = "2006-01-02 15:04:05"
const dateFormat = "2006-01-02"

func Parsetime(timestr string) (time.Time, error) {
	t, e := time.Parse(timeFormat, timestr)
	if e != nil {
		errmsg := fmt.Sprintf("Invalid time format:%v, it should like: %v\n", timestr, timeFormat)
		return time.Time{}, errors.New(errmsg)
	}
	return t, nil
}

func ParseDate(datestr string) (time.Time, error) {
	t, e := time.Parse(dateFormat, datestr)
	if e != nil {
		errmsg := fmt.Sprintf("Invalid date format:%v, it should like: %v\n", datestr, dateFormat)
		return time.Time{}, errors.New(errmsg)
	}
	return t, nil
}
