package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type Email string

type EmailMessage struct {
	Subject string    `json:"subject"`
	Content string    `json:"content"`
	To      Email     `json:"to"`
	From    Email     `json:"from"`
	Date    time.Time `json:"sent"`
}

func (e *EmailMessage) String() string {
	format := "EmailMessage{\n  To:%s\n  From:%s\n  Subject:%s\n  Date:%s\n}"
	return fmt.Sprintf(format, e.To, e.From, e.Subject, e.Date.String())
}

func readStream(url string) io.ReadCloser {
	res, err := http.Get(url)
	if err != nil || res.StatusCode != 200 {
		log.Fatal("Error reading URL", err)
	}

	return res.Body
}

func main() {
	email := &EmailMessage{
		Subject: "Welcome to this message",
		Content: "This email is written directly in the Go file.",
		To:      "hi@example.com",
		From:    "me@example.com",
		Date:    time.Now(),
	}
	fmt.Println("Constructed:", email)

	fmt.Println("To JSON:")
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	encoder.Encode(email)

	stream := readStream("http://echo.jsontest.com/subject/Testing" +
		"/content/Hi, this is a test email/to/you@example.com/from/me@example.com" +
		"/date/2018-10-09T09:43:42.19673586+01:00/")
	defer stream.Close()

	email = &EmailMessage{}
	json.NewDecoder(stream).Decode(email)
	fmt.Println("Downloaded:", email)
}
