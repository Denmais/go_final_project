package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func NextDate(now time.Time, date string, repeat string) (string, error) {
	last, err := time.Parse("20060102", date)
	var deltaInt int
	if err != nil {
		return "", err
	}
	if repeat == "" {
		return "", nil
	}
	if len(strings.Split(repeat, " ")) != 0 {
		if strings.Split(repeat, " ")[0] != "d" && strings.Split(repeat, " ")[0] != "y" {
			return "error", fmt.Errorf("invalid date format")
		}
	}
	if len(strings.Split(repeat, " ")) == 0 {
		return "", fmt.Errorf("invalid date format")
	} else if len(strings.Split(repeat, " ")) > 1 {
		deltaInt, err = strconv.Atoi(strings.Split(repeat, " ")[1])
		if err != nil {
			return "", fmt.Errorf("invalid date format")
		}
		if deltaInt > 400 {
			return "", fmt.Errorf("invalid date format")
		}
	} else {
		deltaInt = 1
	}
	letter := strings.Split(repeat, " ")[0]

	if letter == "d" {
		if len(strings.Split(repeat, " ")) != 2 {
			return "", fmt.Errorf("invalid date format")
		}
		if now.After(last) {
			for now.After(last) {
				last = last.AddDate(0, 0, deltaInt)
			}
		} else {
			last = last.AddDate(0, 0, deltaInt)
		}
	} else if letter == "y" {
		if now.After(last) {
			for now.After(last) {
				last = last.AddDate(1, 0, 0)
			}
		} else {
			last = last.AddDate(1, 0, 0)
		}
	} else {
		return "", fmt.Errorf("invalid date format")
	}
	return last.Format("20060102"), nil
}
