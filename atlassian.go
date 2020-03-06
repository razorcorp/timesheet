package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

/**
 * Package name: main
 * Project name: timesheet
 * Created by: Praveen Premaratne
 * Created on: 01/03/2020 18:34
 */

var secondsInDay = 28800

type TimeLog struct {
	Started   string   `json:"started"`
	TimeSpent string   `json:"timeSpent"`
	Comment   *Comment `json:"comment"`
}

type Comment struct {
	Version     int    `json:"version"`
	CommentType string `json:"type"`
	Content     []*Doc `json:"content"`
}

type Doc struct {
	ContentType string       `json:"type"`
	Content     []*Paragraph `json:"content"`
}
type Paragraph struct {
	Text     string `json:"text"`
	TextType string `json:"type"`
}

type Response struct {
	ErrorMessages []string `json:"errorMessages"`
}

type JiraSearchResult struct {
	Issues []struct {
		Id     string `json:"id"`
		Key    string `json:"key"`
		Fields struct {
			Summary string `json:"summary"`
		} `json:"fields"`
	} `json:"issues"`
}

type WorkLogs struct {
	Key      string
	Summary  string
	Total    int       `json:"total"`
	Worklogs []Worklog `json:"worklogs"`
}

type Worklog struct {
	TimeSpentSeconds int    `json:"timeSpentSeconds"`
	IssueId          string `json:"issueId"`
	Started          string `json:"started"`
	Author           struct {
		EmailAddress string `json:"emailAddress"`
		DisplayName  string `json:"displayName"`
	} `json:"author"`
}

func LogTime(reference string, time string, started string, comment string, domain string, auth string) {
	var slot = TimeLog{}
	slot.TimeSpent = time
	slot.Started = started

	if comment != "" {
		var slotComment = Comment{}
		var doc = Doc{}
		var paragraph = Paragraph{}
		paragraph.Text = comment
		paragraph.TextType = "text"
		doc.ContentType = "paragraph"
		doc.Content = append(doc.Content, &paragraph)
		slotComment.Version = 1
		slotComment.CommentType = "doc"
		slotComment.Content = append(slotComment.Content, &doc)
		slot.Comment = &slotComment
	}
	resp, err := slot.post(reference, domain, auth)
	if err != nil {
		panic(err)
	}

	if len(resp.ErrorMessages) > 0 {
		panic(resp.ErrorMessages)
	}

	fmt.Printf("%s booked to issue %s\n", slot.TimeSpent, reference)
}

func (app *App) GetTimeRemaining(domain string, auth string) {
	var totalTimeSpent int
	var timeRemaining float64
	userEmail, _ := basicAuth(auth)
	issuesOfTheDay, iErr := getIssuesUpdatedToday(domain, auth)
	if iErr != nil {
		panic(iErr)
	}

	worklogs, wErr := issuesOfTheDay.getWorklogs(domain, auth)
	if wErr != nil {
		panic(wErr)
	}

	for _, wLog := range worklogs {
		if wLog.Total > 0 {
			for _, log := range wLog.Worklogs {
				if app.isDateMatch(log.Started) {
					if log.Author.EmailAddress == userEmail {
						if app.History {
							fmt.Println("Timesheet history:")
							fmt.Printf("\t%s: %s\n\t%s: %s\n\t%s: %s\n\t%s: %.2fh\n\n",
								"Key", wLog.Key,
								"Summary", wLog.Summary,
								"Author", log.Author.DisplayName,
								"Time spent", getInHours(log.TimeSpentSeconds))
						}
						totalTimeSpent += log.TimeSpentSeconds
					}
				}
			}
		}
	}

	timeRemaining = getInHours(secondsInDay - totalTimeSpent)
	if timeRemaining < 0 {
		fmt.Printf("oops... Looks like you've booked %.2f hours more that what you supposed to!", timeRemaining)
	} else {
		fmt.Printf("You've %.2f hours ramaining!\n", timeRemaining)
	}

}

func getIssuesUpdatedToday(domain string, auth string) (*JiraSearchResult, error) {
	var client = &http.Client{}
	var query = "jql=worklogDate%20>%3D%20startOfDay()%20AND%20worklogDate%20<%3D%20endOfDay()"
	var url = fmt.Sprintf("https://%s/rest/api/3/search?%s", strings.TrimSuffix(domain, "\n"), query)
	req, reqErr := http.NewRequest("GET", url, nil)
	if reqErr != nil {
		panic(reqErr)
	}

	req.Header.Add("Content-Type", "application/json")
	req.SetBasicAuth(basicAuth(auth))

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var response = new(JiraSearchResult)
	decodeErr := json.NewDecoder(resp.Body).Decode(&response)
	if decodeErr != nil {
		panic(decodeErr)
	}
	return response, nil
}

func (issues *JiraSearchResult) getWorklogs(domain string, auth string) ([]WorkLogs, error) {
	var worklogs []WorkLogs
	for _, issue := range issues.Issues {
		var client = &http.Client{}
		var url = fmt.Sprintf("https://%s/rest/api/3/issue/%s/worklog", strings.TrimSuffix(domain, "\n"), issue.Key)

		req, reqErr := http.NewRequest("GET", url, nil)
		if reqErr != nil {
			panic(reqErr)
		}

		req.Header.Add("Content-Type", "application/json")
		req.SetBasicAuth(basicAuth(auth))

		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			panic(resp.Status)
		}

		var response = WorkLogs{}
		response.Key = issue.Key
		response.Summary = issue.Fields.Summary
		decodeErr := json.NewDecoder(resp.Body).Decode(&response)
		if decodeErr != nil {
			panic(decodeErr)
		}

		worklogs = append(worklogs, response)

	}
	return worklogs, nil
}

func (slot *TimeLog) post(issueId string, domain string, auth string) (*Response, error) {
	client := &http.Client{}
	var url = fmt.Sprintf("https://%s/rest/api/3/issue/%s/worklog", strings.TrimSuffix(domain, "\n"), issueId)

	req, reqErr := http.NewRequest("POST", url, bytes.NewBuffer(slot.json()))
	if reqErr != nil {
		panic(reqErr)
	}

	loginDetails := strings.Split(auth, ":")
	req.Header.Add("Content-Type", "application/json")
	req.SetBasicAuth(loginDetails[0], loginDetails[1])

	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	var response = new(Response)
	decodeErr := json.NewDecoder(resp.Body).Decode(&response)
	if decodeErr != nil {
		panic(decodeErr)
	}
	return response, nil
}

func (slot *TimeLog) json() []byte {
	body, err := json.Marshal(slot)
	if err != nil {
		panic(err)
	}
	return body
}

func basicAuth(token string) (string, string) {
	var loginDetails = strings.Split(token, ":")
	return loginDetails[0], loginDetails[1]
}

func getInHours(seconds int) float64 {
	return float64(seconds) / float64(3600)
}
