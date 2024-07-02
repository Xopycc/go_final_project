package task

import (
	"fmt"
	"strconv"
	"strings"
	"time"
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
		if err != nil || d < 1 || d > 400 {
			return "", fmt.Errorf("%w: invalid days %w", ErrFormat, err)
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
