package main

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

/**
 * Package name: main
 * Project name: timesheet
 * Created by: Praveen Premaratne
 * Created on: 02/03/2020 21:50
 */

var (
	DateFormat, _      = regexp.Compile(`[0-9]{4}-[0-9]{2}-[0-9]{2}`)
	RelativeDateFormat = regexp.MustCompile(`^(?P<Operator>[\-|\+])(?P<Days>[0-9]+)`)
	YmdFormat          = "2006-01-02"
	HmsFormat          = "15:04:05"
)

func (app *App) getDateTime() string {
	var now = time.Now()
	return fmt.Sprintf("%sT%s", now.Format(YmdFormat), app.getTime())
}

func (app *App) getTime() string {
	var now = time.Now()
	return fmt.Sprintf("%s.000+0000", now.Format(HmsFormat))
}

func (app *App) getDate() string {
	return strings.Split(app.Started, "T")[0]
}

func (app *App) getTimeFixed() string {
	return "09:00:00.000+0000"
}

func (app *App) isDateMatch(datetime string) bool {
	dateExpr, _ := regexp.Compile("([0-9]{4}-[0-9]{2}-[0-9]{2})")
	date, err := time.Parse(YmdFormat, dateExpr.FindString(datetime))
	if err != nil {
		panic(err)
	}

	if date.Format(YmdFormat) == app.getDate() {
		return true
	}
	return false
}

func (app *App) isDateBetween(datetime string, start time.Time, end time.Time) bool {
	dateExpr := regexp.MustCompile(`([0-9]{4})-([0-9]{2})-([0-9]{2})T([0-9]{2}):([0-9]{2}):([0-9]{2})`)
	timeArray := dateExpr.FindAllStringSubmatch(datetime, -1)
	year, month, day, hour, minute, second := timeArray[0][1], timeArray[0][2], timeArray[0][3], timeArray[0][4], timeArray[0][5], timeArray[0][6]

	date := time.Date(toInt(year), time.Month(toInt(month)), toInt(day), toInt(hour), toInt(minute), toInt(second), 0, start.Location())

	var normalizedStartDate, _ = fullDay(start)
	var _, normalizedEndDate = fullDay(end)

	if date.After(normalizedStartDate) && date.Before(normalizedEndDate) {
		return true
	}
	return false
}

func (app *App) getWeek() (time.Time, time.Time) {
	var now, err = time.Parse(YmdFormat, app.getDate())
	if err != nil {
		panic(err)
	}
	var weekBegin time.Time
	var weekEnd time.Time
	var dayPos = int(now.Weekday())
	var spanLeft = 1 - dayPos
	var spanRight = 5 - dayPos
	weekBegin = now.AddDate(0, 0, spanLeft)
	weekEnd = now.AddDate(0, 0, spanRight)
	return weekBegin, weekEnd
}

func (app *App) GetDateFromRelative() (*time.Time, error) {
	var now = time.Now()
	var date time.Time
	match := RelativeDateFormat.FindStringSubmatch(app.Started)
	numDays := toInt(match[2])

	if match[1] == "-" {
		date = now.AddDate(0, 0, -numDays)
	} else if match[1] == "+" {
		date = now.AddDate(0, 0, +numDays)
	} else {
		return nil, errors.New("invalid operation given")
	}
	return &date, nil
}

func fullDay(date time.Time) (time.Time, time.Time) {
	yeah, month, day := date.Date()
	return time.Date(yeah, month, day, 0, 0, 0, 0, date.Location()), time.Date(yeah, month, day, 23, 59, 59, 0, date.Location())
}

func getInHours(seconds int) float64 {
	return float64(seconds) / float64(3600)
}

func getDateOfWeek(datetime string) string {
	dateExpr, _ := regexp.Compile("([0-9]{4}-[0-9]{2}-[0-9]{2})")
	date, err := time.Parse(YmdFormat, dateExpr.FindString(datetime))
	if err != nil {
		panic(err)
	}
	return date.Weekday().String()
}

func toInt(s string) int {
	value, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		panic(err)
	}
	return int(value)
}
