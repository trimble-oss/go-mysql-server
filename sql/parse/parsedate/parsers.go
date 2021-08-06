package parsedate

import (
	"fmt"
	"time"
)

type Parser func(result *datetime, chars string) (rest string, err error)

func trimPrefix(count int, str string) string {
	if len(str) > count {
		return str[count:]
	}
	return ""
}

func literalParser(literal byte) Parser {
	return func(dt *datetime, chars string) (rest string, err error) {
		if len(chars) < 1 && literal != ' ' {
			return "", fmt.Errorf("expected literal \"%c\", found empty string", literal)
		}
		chars = takeAllSpaces(chars)
		if literal == ' ' {
			return chars, nil
		}
		if chars[0] != literal {
			return "", fmt.Errorf("expected literal \"%c\", got \"%c\"", literal, chars[0])
		}
		return trimPrefix(1, chars), nil
	}
}

func parseAmPm(result *datetime, chars string) (rest string, err error) {
	if len(chars) < 2 {
		return "", fmt.Errorf("expected > 2 chars, found %d", len(chars))
	}
	switch chars[:2] {
	case "am":
		result.am = boolPtr(true)
	case "pm":
		result.am = boolPtr(false)
	default:
		return "", err
	}
	return trimPrefix(2, chars), nil
}

func parseWeedayAbbreviation(result *datetime, chars string) (rest string, err error) {
	if len(chars) < 3 {
		return "", err
	}
	weekday, ok := weekdayAbbrev(chars[:3])
	if !ok {
		return "", err
	}
	result.weekday = &weekday
	return trimPrefix(3, chars), nil
}

func parseMonthAbbreviation(result *datetime, chars string) (rest string, err error) {
	if len(chars) < 3 {
		return "", err
	}
	month, ok := monthAbbrev(chars[:3])
	if !ok {
		return "", err
	}
	result.month = &month
	return trimPrefix(3, chars), nil
}

func parseMonthNumeric(result *datetime, chars string) (rest string, err error) {
	num, rest, err := takeNumber(chars)
	if err != nil {
		return "", err
	}
	month := time.Month(num)
	result.month = &month
	return rest, nil
}

func parseDayOfMonthNumeric(result *datetime, chars string) (rest string, err error) {
	num, rest, err := takeNumber(chars)
	if err != nil {
		return "", err
	}
	result.day = uintPtr(uint(num))
	return rest, nil
}

func parseMicrosecondsNumeric(result *datetime, chars string) (rest string, err error) {
	num, rest, err := takeNumber(chars)
	if err != nil {
		return "", err
	}
	result.microseconds = uintPtr(uint(num))
	return rest, nil
}

func parse24HourNumeric(result *datetime, chars string) (rest string, err error) {
	hour, rest, err := takeNumber(chars)
	if err != nil {
		return "", err
	}
	result.hours = uintPtr(uint(hour))
	return rest, nil
}

func parse12HourNumeric(result *datetime, chars string) (rest string, err error) {
	num, rest, err := takeNumber(chars)
	if err != nil {
		return "", err
	}
	result.hours = uintPtr(uint(num))
	return rest, nil
}

func parseMinuteNumeric(result *datetime, chars string) (rest string, err error) {
	min, rest, err := takeNumber(chars)
	if err != nil {
		return "", err
	}
	result.minutes = uintPtr(uint(min))
	return rest, nil
}

func parseMonthName(result *datetime, chars string) (rest string, err error) {
	month, charCount, ok := monthName(chars)
	if !ok {
		return "", err
	}
	result.month = &month
	return trimPrefix(charCount, chars), nil
}

func parse12HourTimestamp(result *datetime, chars string) (rest string, err error) {
	hour, rest, err := takeNumber(chars)
	if err != nil {
		return "", err
	}
	rest, err = literalParser(':')(result, rest)
	if err != nil {
		return "", err
	}
	min, rest, err := takeNumber(rest)
	if err != nil {
		return "", err
	}
	rest, err = literalParser(':')(result, rest)
	if err != nil {
		return "", err
	}
	sec, rest, err := takeNumber(rest)
	if err != nil {
		return "", err
	}
	rest = takeAllSpaces(rest)
	rest, err = parseAmPm(result, rest)
	if err != nil {
		return "", err
	}
	result.seconds = uintPtr(uint(sec))
	result.minutes = uintPtr(uint(min))
	result.hours = uintPtr(uint(hour))
	return rest, nil
}

func parseSecondsNumeric(result *datetime, chars string) (rest string, err error) {
	sec, rest, err := takeNumber(chars)
	if err != nil {
		return "", err
	}
	result.seconds = uintPtr(uint(sec))
	return rest, nil
}

func parse24HourTimestamp(result *datetime, chars string) (rest string, err error) {
	hour, rest, err := takeNumber(chars)
	if err != nil {
		return "", err
	}
	rest, err = literalParser(':')(result, rest)
	if err != nil {
		return "", err
	}
	minute, rest, err := takeNumber(rest)
	if err != nil {
		return "", err
	}
	rest, err = literalParser(':')(result, rest)
	if err != nil {
		return "", err
	}
	seconds, rest, err := takeNumber(rest)
	if err != nil {
		return "", err
	}
	result.hours = uintPtr(uint(hour))
	result.minutes = uintPtr(uint(minute))
	result.seconds = uintPtr(uint(seconds))
	return rest, err
}

func parseYear2DigitNumeric(result *datetime, chars string) (rest string, err error) {
	if len(chars) < 2 {
		return "", err
	}
	year, rest, err := takeNumber(chars)
	if err != nil {
		return "", err
	}
	if year >= 100 {
		return "", err
	}
	if year >= 70 {
		year += 1900
	} else {
		year += 2000
	}
	result.year = uintPtr(uint(year))
	return rest, nil
}

func parseYear4DigitNumeric(result *datetime, chars string) (rest string, err error) {
	if len(chars) < 4 {
		return "", err
	}
	year, rest, err := takeNumber(chars)
	if err != nil {
		return "", err
	}
	result.year = uintPtr(uint(year))
	return rest, nil
}

func parseDayNumericWithEnglishSuffix(result *datetime, chars string) (rest string, err error) {
	num, rest, err := takeNumber(chars)
	if err != nil {
		return "", err
	}
	result.day = uintPtr(uint(num))
	return trimPrefix(2, rest), nil
}
