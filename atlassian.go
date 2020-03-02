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

func (slot *TimeLog) post(issueId string, domain string, auth string) (*Response, error) {
	client := &http.Client{}
	var url = fmt.Sprintf("https://%s/rest/api/3/issue/%s/worklog", strings.TrimSuffix(domain, "\n"), issueId)
	//fmt.Println(url)

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
