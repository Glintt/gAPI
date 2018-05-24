package utils

import "time"

func CurrentDate() string {
	return time.Now().Format("2006-01-02")
}

func CurrentDateWithFormat(format string) string {
	return time.Now().Format(format)
}

func CurrentTimeMilliseconds() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}
func NanosecondsToMilliseconds(nano int64) int64 {
	return nano / int64(time.Millisecond)
}
