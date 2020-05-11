package main

import (
	"fmt"
	"os"
)

var SIGNATURE = `
******************************************
*                                        *
*      ********         ********         *
*      **      **       **      **       *
*      **       **      **       **      *
*      **      **       **      **       *
*      ********         ********         *
*      **               **               *
*      **               **               *
*                                        *
******************************************
         ** Praveen Premaratne **

 * Project name: timesheet
 * Created by: Praveen Premaratne
 * Created on: 29/02/2020 17:51

`

type (
	App struct {
		Ticket        string
		Comment       string
		Started       string
		TimeSpent     string
		Help          bool
		Encode        string
		TimeRemaining bool
		History       bool
		PrintWeek     bool
		Version       bool
		Update        bool
		Configuration struct {
			Auth   string
			Domain string
		}
	}
)

type Application interface {
	Parser()
	CredentialEncode()
	GetTimeRemaining(domain string, auth string)
	GetHistory()
	GetWeekTimesheet(domain string, auth string)
}

var VERSION string
var AppName = "timesheet"

func main() {
	var app App

	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}()

	app.Parser()
	app.loadConf()
	app.updatable()

	if app.Update {
		fmt.Println("Checking for updates....")
		release := app.update()
		status := release.installUpdate()
		if status != nil {
			panic(status.Error())
		}
		os.Exit(0)
	}

	fmt.Println("This might take a moment....")

	if app.TimeRemaining {
		app.GetTimeRemaining(app.Configuration.Domain, app.Configuration.Auth)
		os.Exit(0)
	}

	if app.History {
		app.GetHistory()
		os.Exit(0)
	}

	if app.PrintWeek {
		app.GetWeekTimesheet(app.Configuration.Domain, app.Configuration.Auth)
		os.Exit(0)
	}

	LogTime(app.Ticket, app.TimeSpent, app.Started, app.Comment, app.Configuration.Domain, app.Configuration.Auth)
	app.GetTimeRemaining(app.Configuration.Domain, app.Configuration.Auth)

}
