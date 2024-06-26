package utils

import (
	"fmt"
	"time"
)

var Months = map[string][]string{
	"en": {"January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December"},
	"ru": {"января", "февраля", "марта", "апреля", "мая", "июня", "июля", "августа", "сентября", "октября", "ноября", "декабря"},
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

	date, err = time.Parse("2006.01.02", dateString)
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

func ParseNameDateTimeCityString(nameDateTimeCityString string) (string, string, string, string, error) {
	//var regexpVariants = []string{
	//	`(\d{4}.\d{2}.\d{2})`,
	//	`(\d{2}.\d{2}.\d{4})`,
	//	`(\d{4}-\d{2}-\d{2})`,
	//	`(\d{2}/\d{2}/\d{4})`,
	//}
	//
	//var date string
	//for _, regexpVariant := range regexpVariants {
	//	re := regexp.MustCompile(regexpVariant)
	//	date = re.FindString(nameDateTimeCityString)
	//	if date != "" {
	//		break
	//	}
	//}
	//
	//var sDate = &database.ShortDate{}
	//var sTime = &database.ShortTime{}
	//
	//x := strings.Split(nameDateTimeCityString, date)
	//
	//parsedDate, err := ParseDateFromString(date)
	//if err != nil {
	//	return "", sDate, sTime, "", err
	//}
	//date = parsedDate.Format("2006-01-02")
	//sDate.UnmarshalJSON([]byte(date))
	//
	//println("date:", date)
	//
	//if date == "" {
	//	return "", nil, nil, "", fmt.Errorf("can't parse nameDateTimeCityString: no date found")
	//}
	//
	//fmt.Println(x)
	//
	//if len(x) != 2 {
	//	return "", nil, nil, "", fmt.Errorf("can't parse nameDateTimeCityString: split error")
	//}
	//
	//name := strings.TrimSpace(x[0])
	//rightPart := strings.TrimSpace(x[1])
	//
	//re := regexp.MustCompile(`\d{2}:\d{2}`)
	//t := re.FindString(rightPart)
	//
	//var city string
	//if t == "" {
	//	city = rightPart
	//} else {
	//	cityParts := strings.SplitN(rightPart, t, 2)
	//	if len(cityParts) == 2 {
	//		city = strings.TrimSpace(cityParts[1])
	//	} else {
	//		return "", nil, nil, "", fmt.Errorf("can't parse nameDateTimeCityString: time split error")
	//	}
	//}
	//
	//if t != "" {
	//	parsedTime, err := ParseTimeFromString(t)
	//
	//	if err != nil {
	//		return "", nil, nil, "", err
	//	}
	//	t = parsedTime.Format("15:04")
	//	err = sTime.UnmarshalJSON([]byte(t))
	//	if err != nil {
	//		return "", nil, nil, "", err
	//	}
	//
	//	return name, sDate, sTime, city, nil
	//}
	//
	//return name, sDate, nil, city, nil
	return "", "", "", "", fmt.Errorf("not implemented")
}
