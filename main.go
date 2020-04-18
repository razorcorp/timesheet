package main

import (
	"encoding/json"
	"fmt"
	"net/http"
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
	GetHistory()
	GetWeekTimesheet(domain string, auth string)
}

var VERSION string

func (app *App) upgrade() {
	var client = &http.Client{}
	req, rErr := http.NewRequest("GET", "https://api.github.com/repos/praveenprem/timesheet/releases/latest", nil)
	if rErr != nil {
		panic(rErr)
	}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	var response struct {
		Name string `json:"name"`
		URL  string `json:"html_url"`
	}
	decodeErr := json.NewDecoder(resp.Body).Decode(&response)
	if decodeErr != nil {
		panic(decodeErr)
	}

	if fmt.Sprintf("v%s", VERSION) != response.Name {
		fmt.Println("New version available! Please download the latest release from", response.URL)
	}
}

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
	app.upgrade()

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
	app.TimeRemaining = true
	app.GetTimeRemaining(app.Configuration.Domain, app.Configuration.Auth)

}
