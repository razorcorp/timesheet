package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"regexp"
	"time"
)

/**
 * Package name: main
 * Project name: timesheet
 * Created by: Praveen Premaratne
 * Created on: 29/02/2020 17:52
 */

var dateFormat, _ = regexp.Compile("[0-9]{2}-[0-9]{2}-[0-9]{4}")

func (app *App) Parser() {
	var day = time.Now().Format("02-01-2006")

	flag.BoolVar(&app.Help, "h", false, "HELP: Print usage")
	flag.StringVar(&app.Ticket, "r", "",
		"REQUIRED: Jira ticket reference. E.g. DDSP-4")
	flag.StringVar(&app.TimeSpent, "t", "",
		"REQUIRED: The time spent as days (#d), hours (#h), or minutes (#m or #). E.g. 8h")
	flag.StringVar(&app.Started, "d", "",
		fmt.Sprintf("OPTIONAL: The date on which the worklog effort was started in DD-MM-YYYY format. Default %s", day))
	flag.StringVar(&app.Comment, "m", "",
		"OPTIONAL: A comment about the worklog")
	flag.StringVar(&app.Encode, "e", "", "HELP: Base64 encode the given credentials."+
		" Format: email:token;domain. e.g. example@example.com:abcThisIsFake;xyz.atlassian.net")
	flag.BoolVar(&app.TimeRemaining, "remaining", false, "Time remaining")
	flag.BoolVar(&app.History, "history", false, "Print the timesheet of the day")
	flag.Parse()
	app.validate()
}

func (app *App) validate() {
	if app.Help {
		fmt.Println("This tool can be used to log time spent on a specific Jira ticket on a project.")
		app.usage()
		os.Exit(0)
	}

	if app.TimeRemaining {
		return
	}

	if app.Encode != "" {
		app.CredentialEncode()
		os.Exit(0)
	}

	if app.Ticket == "" {
		panic(errors.New("please provide a ticket reference. -r"))
	}

	if app.TimeSpent == "" {
		panic(errors.New("no time given. -t"))
	}

	if app.Started != "" {
		if !dateFormat.MatchString(app.Started) {
			panic("provided date didn't match expected format. try -h for help")
		} else {
			app.Started = fmt.Sprintf("%sT%s", app.Started, app.getTimeFixed())
		}
	} else {
		app.Started = app.getDateTime()
	}
}

func (app *App) usage() {
	flag.PrintDefaults()
}
