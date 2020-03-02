package main

import (
	"encoding/base64"
	"fmt"
	"os"
	"strings"
)

/**
 * Package name: main
 * Project name: timesheet
 * Created by: Praveen Premaratne
 * Created on: 01/03/2020 18:28
 */

func (app *App) loadConf() {
	if rawConf := os.Getenv("TIMESHEET"); rawConf == "" {
		panic("please export \"TIMESHEET\" with Base64 encoded Atlassian data in following format: email:token;domain")
	} else {
		if conf, err := base64.StdEncoding.DecodeString(rawConf); err != nil {
			panic("config is not Base64 encoded.")
		} else {
			var config = strings.Split(string(conf), ";")
			app.Configuration.Auth = config[0]
			app.Configuration.Domain = config[1]
		}
	}

}

func (app *App) CredentialEncode() {
	var token = base64.StdEncoding.EncodeToString([]byte(app.Encode))
	fmt.Println(token)
}
