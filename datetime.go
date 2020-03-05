package main

import (
	"fmt"
	"regexp"
	"time"
)

/**
 * Package name: main
 * Project name: timesheet
 * Created by: Praveen Premaratne
 * Created on: 02/03/2020 21:50
 */

func (app *App) getDateTime() string {
	var now = time.Now()
	return fmt.Sprintf("%sT%s", now.Format("2006-01-02"), app.getTime())
}

func (app *App) getTime() string {
	var now = time.Now()
	return fmt.Sprintf("%s.000+0000", now.Format("15:04:05"))
}

func (app *App) getTimeFixed() string {
	return "09:00:00.000+0000"
}

func (app *App) isDateMatch(datetime string) bool {
	dateExpr, _ := regexp.Compile("([0-9]{4}-[0-9]{2}-[0-9]{2})")
	date, err := time.Parse("2006-01-02", dateExpr.FindString(datetime))
	if err != nil {
		panic(err)
	}

	if date.Format("2006-01-02") == time.Now().Format("2006-01-02") {
		return true
	}
	return false
}
