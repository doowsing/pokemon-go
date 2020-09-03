package utils

import "time"

func ToDayStartUnix() int {
	now := time.Now()
	return int(time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).Unix())
}

func NowUnix() int {
	return int(time.Now().Unix())
}
