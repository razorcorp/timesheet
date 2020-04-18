package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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

type WeekLog struct {
	Total  int
	Issues []Issue
}

type Issue struct {
	Key  string
	Logs []DayLog
}

type DayLog struct {
	Total     int
	WeekDay   string
	TimeSpent int
}

type Week struct {
	Total int
	Days  map[string]map[string][]int
}

var daysOfWeek = []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday"}

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
	issuesOfTheDay, iErr := getIssuesUpdatedToday(domain, auth, app.getDate())
	if iErr != nil {
		panic(iErr)
	}

	workLogs, wErr := issuesOfTheDay.getWorklogs(domain, auth)
	if wErr != nil {
		panic(wErr)
	}

	for _, wLog := range workLogs {
		if wLog.Total > 0 {
			for _, log := range wLog.Worklogs {
				if app.isDateMatch(log.Started) {
					if log.Author.EmailAddress == userEmail {
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

func (app *App) GetHistory() {
	var totalTimeSpent int
	userEmail, _ := basicAuth(app.Configuration.Auth)
	issuesOfTheDay, iErr := getIssuesUpdatedToday(app.Configuration.Domain, app.Configuration.Auth, app.getDate())
	if iErr != nil {
		panic(iErr)
	}

	workLogs, wErr := issuesOfTheDay.getWorklogs(app.Configuration.Domain, app.Configuration.Auth)
	if wErr != nil {
		panic(wErr)
	}

	fmt.Printf("Timesheet history: (%s):\n", app.getDate())

	for _, wLog := range workLogs {
		if wLog.Total > 0 {
			for _, log := range wLog.Worklogs {
				if app.isDateMatch(log.Started) {
					if log.Author.EmailAddress == userEmail {
						fmt.Printf("\t%s:\n\t\t%s: %s\n\t\t%s: %s\n\t\t%s: %s\n\t\t%s: %.2fh\n\n",
							wLog.Key,
							"Summary", wLog.Summary,
							"Author", log.Author.DisplayName,
							"Comment", log.Comment.Content[0].Content[0].Text,
							"Time spent", getInHours(log.TimeSpentSeconds),
						)
						totalTimeSpent += log.TimeSpentSeconds
					}
				}
			}
		}
	}

	fmt.Println(fmt.Sprintf("Total %.1fh", getInHours(totalTimeSpent)))

}

func (app *App) GetWeekTimesheet(domain string, auth string) {
	var weekLog WeekLog
	start, end := app.getWeek()
	userEmail, _ := basicAuth(auth)
	issuesOfTheWeek, iErr := getIssuesUpdatedBetweenDays(domain, auth,
		start.Format("2006-01-02"), end.Format("2006-01-02"))
	if iErr != nil {
		panic(iErr)
	}

	worklogs, wErr := issuesOfTheWeek.getWorklogs(domain, auth)
	if wErr != nil {
		panic(wErr)
	}

	for _, wLog := range worklogs {
		var issue Issue
		issue.Key = wLog.Key
		if wLog.Total > 0 {
			for _, log := range wLog.Worklogs {
				if app.isDateBetween(log.Started, start, end) {
					if log.Author.EmailAddress == userEmail {
						var dayLog DayLog
						dayLog.WeekDay = getDateOfWeek(log.Started)
						dayLog.TimeSpent = log.TimeSpentSeconds
						weekLog.Total += log.TimeSpentSeconds
						issue.Logs = append(issue.Logs, dayLog)
					}
				}
			}
		}
		if len(issue.Logs) > 0 {
			weekLog.Issues = append(weekLog.Issues, issue)
		}
	}

	weekLog.print()
}

func getIssuesUpdatedToday(domain string, auth string, date string) (*JiraSearchResult, error) {
	var client = &http.Client{}
	var query = fmt.Sprintf("jql=worklogDate%%20>%%3D%%20\"%s\"%%20AND%%20worklogDate%%20<%%3D%%20\"%s\"", date, date)
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

func getIssuesUpdatedBetweenDays(domain string, auth string, start string, end string) (*JiraSearchResult, error) {
	var query = fmt.Sprintf("jql=worklogDate%%20>%%3D%%20\"%s\"%%20AND%%20worklogDate%%20<%%3D%%20\"%s\"", start, end)
	var url = fmt.Sprintf("https://%s/rest/api/3/search?%s", strings.TrimSuffix(domain, "\n"), query)

	var headers []struct {
		key   string
		value string
	}

	headers = append(headers, struct {
		key   string
		value string
	}{key: "Content-Type", value: "application/json"})

	var response = new(JiraSearchResult)
	decodeErr := json.NewDecoder(httpReq("GET", url, auth, nil, headers).Body).Decode(&response)
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

func httpReq(method string, url string, auth string, body io.Reader, headers []struct {
	key   string
	value string
}) *http.Response {
	var client = &http.Client{}
	req, reqErr := http.NewRequest(method, url, body)
	if reqErr != nil {
		panic(reqErr)
	}

	for _, header := range headers {
		req.Header.Add(header.key, header.value)
	}
	req.SetBasicAuth(basicAuth(auth))
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	//defer resp.Body.Close()
	return resp

}

func (w *WeekLog) sort() Week {
	var sortedWeek Week

	sortedWeek.Days = make(map[string]map[string][]int)

	sortedWeek.Total = w.Total
	for _, issue := range w.Issues {
		sortedWeek.Days[issue.Key] = make(map[string][]int)
		for _, day := range issue.Logs {
			sortedWeek.Days[issue.Key][day.WeekDay] = append(sortedWeek.Days[issue.Key][day.WeekDay], day.TimeSpent)
		}
	}

	return sortedWeek.fillGaps()
}

func (w *Week) fillGaps() Week {
	for _, days := range w.Days {
		for _, d := range daysOfWeek {
			if _, found := days[d]; !found {
				days[d] = []int{0}
			}
		}
	}

	return *w
}

func (w *Week) sum() Week {
	for _, days := range w.Days {
		for _, d := range daysOfWeek {
			if times, found := days[d]; !found {
				sum := 0
				for _, time := range times {
					sum += time
				}
				days[d] = []int{sum}
			}
		}
	}

	return *w
}

func (w *WeekLog) print() {
	for i := 0; i <= 83; i++ {
		if i == 0 {
			fmt.Printf(" ")
		} else if i == 83 {
			fmt.Printf("\n")
		} else {
			fmt.Printf("_")
		}
	}

	fmt.Printf("| %-15s ", "Issue")
	for _, title := range daysOfWeek {
		fmt.Printf("| %-10s ", title)
	}
	fmt.Printf("|\n")
	for i := 0; i <= 82; i++ {
		if i == 0 {
			fmt.Printf("|")
		} else if i == 18 || i == 31 || i == 44 || i == 57 || i == 70 {
			fmt.Printf("|")
		} else if i == 82 {
			fmt.Printf("_|\n")
		} else {
			fmt.Printf("_")
		}
	}

	weekSorted := w.sort()
	var processedIssues int

	for issue, day := range weekSorted.sum().Days {
		if processedIssues > 0 {
			for i := 0; i <= 83; i++ {
				if i == 0 {
					fmt.Printf("|")
				} else if i == 83 {
					fmt.Printf("|\n")
				} else if i == 18 || i == 31 || i == 44 || i == 57 || i == 70 {
					fmt.Printf("|")
				} else {
					fmt.Printf("-")
				}
			}
		}
		fmt.Printf("| %-15s ", issue)
		for _, v := range daysOfWeek {
			if day[v][0] == 0 {
				fmt.Printf("| %-10s ", "")
			} else {
				fmt.Printf("| %-10.1f ", getInHours(day[v][0]))
			}
		}
		fmt.Println("|")
		processedIssues += 1
	}

	for i := 0; i <= 83; i++ {
		if i == 0 {
			fmt.Printf(" ")
		} else if i == 83 {
			fmt.Printf(" \n")
		} else {
			fmt.Printf("-")
		}
	}

	fmt.Println(fmt.Sprintf("Total %.1fh", getInHours(weekSorted.Total)))
}
