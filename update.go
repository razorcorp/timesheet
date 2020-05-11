package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

/**
 * Package name: main
 * Project name: timesheet
 * Created by: Praveen Premaratne
 * Created on: 10/05/2020 22:21
 */
type (
	Update struct {
		Name   string `json:"name"`
		URL    string `json:"html_url"`
		Assets []struct {
			Url         string `json:"url"`
			Name        string `json:"name"`
			ContentType string `json:"content_type"`
		}
	}
)

func (app *App) update() Update {
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
	var response Update
	decodeErr := json.NewDecoder(resp.Body).Decode(&response)
	if decodeErr != nil {
		panic(decodeErr)
	}
	return response
}

func (app *App) updatable() {
	remoteVersion := app.update()

	if fmt.Sprintf("v%s", VERSION) != remoteVersion.Name {
		fmt.Printf("New version %s available! Use -update to download the new vresion\n\n", remoteVersion.Name)
	}
}
func (u *Update) installUpdate() error {
	var binary []byte
	if fmt.Sprintf("v%s", VERSION) != u.Name {
		fmt.Printf("New version available!\nInstalling %s...\n", u.Name)
		for _, asset := range u.Assets {
			if asset.Name == AppName {
				client := &http.Client{}
				if req, err := http.NewRequest("GET", asset.Url, nil); err != nil {
					panic(err)
				} else {
					req.Header.Add("Accept", asset.ContentType)
					if resp, err := client.Do(req); err != nil {
						panic(err)
					} else {
						if data, err := ioutil.ReadAll(resp.Body); err != nil {
							resp.Body.Close()
							panic(err)
						} else {
							binary = data
						}
						resp.Body.Close()
					}
				}
			}
		}
	} else {
		return errors.New(fmt.Sprintf("You’ve got the latest version of %s\n", AppName))
	}

	return u.write(binary)
}

func (u *Update) write(binary []byte) error {
	if len(binary) < 1 {
		return errors.New("no data to write")
	}

	appFile, err := os.OpenFile(fmt.Sprintf("/tmp/%s", AppName), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0751)
	if err != nil {
		return err
	}
	defer appFile.Close()
	if _, err := appFile.Write(binary); err != nil {
		return err
	}

	return nil
}
