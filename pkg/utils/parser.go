package utils

import (
	"fmt"
	"time"
)

var Months = map[string][]string{
	"en": {"January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December"},
	"ru": {"января", "февраля", "марта", "апреля", "мая", "июня", "июля", "августа", "сентября", "октября", "ноября", "декабря"},
}

func DateFormat(date time.Time) string {
	return fmt.Sprintf("%d %s %d", date.Day(), Months["ru"][date.Month()-1], date.Year())
}

func ParseDateFromString(dateString string) (time.Time, error) {
	date, err := time.Parse("2006-01-02", dateString)
	if err == nil {
		return date, nil
	}
	date, err = time.Parse("02.01.2006", dateString)
	if err == nil {
		return date, nil
	}
	date, err = time.Parse("02/01/2006", dateString)
	if err == nil {
		return date, nil
	}

	date, err = time.Parse("02-01-2006", dateString)
	if err == nil {
		return date, nil
	}
	date, err = time.Parse("02 January 2006", dateString)
	if err == nil {
		return date, nil
	}
	return date, err
}

func ParseTimeFromString(timeString string) (time.Time, error) {
	t, err := time.Parse("15:04", timeString)
	if err == nil {
		return t, nil
	}
	return t, err
}
