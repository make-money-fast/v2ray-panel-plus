package helpers

import (
	"strconv"
	"time"
)

func NowDateTime() string {
	return time.Now().Format(`2006-01-02 15:04:05`)
}

func Str2Int(s string) int {
	a, _ := strconv.Atoi(s)
	return a
}
