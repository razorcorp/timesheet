package main

import (
	"fmt"
	"os"
)

/**

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

 * Package name: main
 * Project name: timesheet
 * Created by: Praveen Premaratne
 * Created on: 29/02/2020 17:51
*/

type App struct {
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
	Configuration struct {
		Auth   string
		Domain string
	}
}

type Application interface {
	Parser()
	CredentialEncode()
	GetTimeRemaining(domain string, auth string)
	GetWeekTimesheet(domain string, auth string)
}

var VERSION string

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

	if app.TimeRemaining {
		app.GetTimeRemaining(app.Configuration.Domain, app.Configuration.Auth)
		os.Exit(0)
	}

	if app.PrintWeek {
		app.GetWeekTimesheet(app.Configuration.Domain, app.Configuration.Auth)
		os.Exit(0)
	}

	LogTime(app.Ticket, app.TimeSpent, app.Started, app.Comment, app.Configuration.Domain, app.Configuration.Auth)
	app.TimeRemaining = true
	app.GetTimeRemaining(app.Configuration.Domain, app.Configuration.Auth)

}
