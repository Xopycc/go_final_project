package task

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	MinDays = 1
	MaxDays = 400
)

func NextDate(now time.Time, startDate string, repeat string) (string, error) {

	repeat = strings.ToLower(strings.TrimSpace(repeat))
	startDt, err := time.Parse(DateFormat, startDate)
	if err != nil {
		return "", fmt.Errorf("%w: unexpected date value", ErrFormat)
	}

	var nextDate time.Time

	switch {
	case repeat == "y":
		nextDate = startDt
		nextDate = nextDate.AddDate(1, 0, 0)
		for nextDate.Before(now) {
			nextDate = nextDate.AddDate(1, 0, 0)
		}

	case strings.HasPrefix(repeat, "d"):
		d, err := strconv.Atoi(strings.TrimSpace(strings.TrimPrefix(repeat, "d ")))
		if err != nil || d < MinDays || d > MaxDays {
			return "", fmt.Errorf("%v: invalid days %w", ErrFormat, err) // %v для вывода err
		}
		nextDate = startDt
		nextDate = nextDate.AddDate(0, 0, d)
		for nextDate.Before(now) {
			nextDate = nextDate.AddDate(0, 0, d)
		}
	default:
		return "", fmt.Errorf("%w: unexpected repeat value", ErrFormat)

	}
	return nextDate.Format(DateFormat), nil
}
