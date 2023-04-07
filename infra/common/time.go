package common

import "time"

const (
	fileDateFormat = "2006_01_02"
)

func GetNowTime() time.Time {
	return time.Now().UTC()
}

func GetFileDateFormat(t time.Time) string {
	return t.Format(fileDateFormat)
}
