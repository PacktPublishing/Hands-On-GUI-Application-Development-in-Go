package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type Email string

type EmailMessage struct {
	Subject string
	Content string
	To      Email
	From    Email
	Date    time.Time `xml:"Sent"`
}

func (e *EmailMessage) String() string {
	format := "EmailMessage{\n  To:\t%s\n  From:\t%s\n  Subject:\t%s\n  Date:\t%s\n\n  %s\n}"
	return fmt.Sprintf(format, e.To, e.From, e.Subject, e.Date.String(), e.Content)
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

	fmt.Println("To XML:")
	encoder := xml.NewEncoder(os.Stdout)
	encoder.Indent("", "  ")
	encoder.Encode(email)

	stream := readStream("http://urlecho.appspot.com/echo?status=200&Content-Type=text%2Fxml&" +
		"body=%3CEmailMessage%3E%0A%20%20%3CSubject%3ETesting%3C%2FSubject%3E%0A%20%20" +
		"%3CContent%3EHi,%20this%20is%20a%20test%20email%3C%2FContent%3E%0A%20%20%3CTo%3Eyou%40example.com%3C" +
		"%2FTo%3E%0A%20%20%3CFrom%3Eme%40example.com%3C%2FFrom%3E%0A%20%20%3CSent%3E2018-10-09T09%3A" +
		"55%3A45.568777985%2B01%3A00%3C%2FSent%3E%0A%3C%2FEmailMessage%3E")
	defer stream.Close()

	email = &EmailMessage{}
	xml.NewDecoder(stream).Decode(email)
	fmt.Println("Downloaded:", email)
}
