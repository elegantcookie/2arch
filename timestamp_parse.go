package main

import (
	"fmt"
	"strconv"
	"time"
)

type DateTime struct {
	Date string
	Time string
}

func timestampToDatetime(timestamp int) *DateTime {
	i, err := strconv.ParseInt(fmt.Sprintf("%d", timestamp), 10, 60)
	handleError(err)

	tm := time.Unix(i, 0)

	daysOfWeek := map[string]string{
		"Sunday":    "Вскр",
		"Monday":    "Пон",
		"Tuesday":   "Втрк",
		"Wednesday": "Срд",
		"Thursday":  "Чтв",
		"Friday":    "Птн",
		"Saturday":  "Суб",
	}
	date := tm.Format("02/01/06")
	weekDay := tm.Weekday()
	date += fmt.Sprintf(" %s", daysOfWeek[weekDay.String()])
	time := tm.Format("15:04:05")
	return &DateTime{Time: time, Date: date}

}
