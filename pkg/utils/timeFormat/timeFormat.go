package timeFormat

import (
	"strings"
	"time"
)

func GetTimeFormatForRest(month int, day int) (string, string) {
	startDay := strings.Replace(time.Now().AddDate(0, month, day).Format("2006-01-02T15:04:05.000-0700"), "+", "-", 1)
	endDay := strings.Replace(time.Now().Format("2006-01-02T15:04:05.000-0700"), "+", "-", 1)

	return startDay, endDay
}
