package utils

import (
	"errors"
	"fmt"
	"time"
)

const base_time_format = "2006-01-02 15:04:05"

func NowUnix() int {
	return int(time.Now().Unix())
}

func StrParseMustTime(timeString string) time.Time {
	t, _ := time.Parse(base_time_format, timeString)
	return t
}

func FormatTime(t time.Time) string {
	return t.Format(base_time_format)
}

func StrParseTime(timeString string) (time.Time, error) {
	t, err := time.Parse(base_time_format, timeString)
	return t, err
}

func YmdStrParseTime(timeString string) (time.Time, error) {
	if len(timeString) != 8 {
		return time.Time{}, errors.New("参数出错！")
	}

	t, err := time.Parse(base_time_format, timeString[:4]+"-"+timeString[4:6]+"-"+timeString[6:8]+" 00:00:00")
	return t, err
}

func TimeFormatYmd(t time.Time) string {
	y, m, d := t.Date()
	return fmt.Sprintf("%d%02d%02d", y, m, d)
}

func DurationFormatHms(t time.Duration) string {
	h := int(t.Hours())
	m := int(t.Minutes()) - h*60
	s := int(t.Seconds()) - m*60
	fstring := ""
	if h < 10 {
		fstring += "0%d"
	} else {
		fstring += "%d"
	}
	fstring += ":"
	if m < 10 {
		fstring += "0%d"
	} else {
		fstring += "%d"
	}
	fstring += ":"
	if s < 10 {
		fstring += "0%d"
	} else {
		fstring += "%d"
	}
	return fmt.Sprintf(fstring, h, m, s)
}

func ToDayStartUnix() int {
	now := time.Now()
	newTime := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	return int(newTime.Unix())
}
