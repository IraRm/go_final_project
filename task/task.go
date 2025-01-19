package task

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

const TimeFormat string = "20060102"

type Task struct {
	Id      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func NextTime(now time.Time, repeatRule string) (time.Time, error) {
	if repeatRule == "" {
		return now, nil
	}

	parts := strings.Split(repeatRule, " ")
	switch parts[0] {
	case "d":
		if len(parts) != 2 {
			return time.Time{}, fmt.Errorf("unknown repeat rule type")
		}
		days, err := strconv.Atoi(parts[1])
		if err != nil || days <= 0 || days > 400 {
			return time.Time{}, fmt.Errorf("invalid number of days")
		}
		return now.AddDate(0, 0, days), nil
	case "y":
		return now.AddDate(1, 0, 0), nil
	default:
		return time.Time{}, fmt.Errorf("unknown repeat rule type")
	}
}

func (t *Task) NextDate() (string, error) {
	currentTime, err := time.Parse(TimeFormat, t.Date)
	if err != nil {
		return "", err
	}

	nT, err := NextTime(currentTime, t.Repeat)
	if err != nil {
		return "", err
	}

	return nT.Format(TimeFormat), nil
}
