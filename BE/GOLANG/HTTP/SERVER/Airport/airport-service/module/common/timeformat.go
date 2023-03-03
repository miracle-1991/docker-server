package common

import "time"

const timeFormat = "2006-01-02 15:04"
const dateFormat = "2006-01-02"
const hourFormat = "15"

func Date(timestr string) (string, error) {
	t, err := time.Parse(timeFormat, timestr)
	if err != nil {
		return "", err
	}
	return t.Format(dateFormat), nil
}

func Hour(timestr string) (string, error) {
	t, err := time.Parse(timeFormat, timestr)
	if err != nil {
		return "", err
	}
	return t.Format(hourFormat), nil
}
