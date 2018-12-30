package client

import "time"

type Email string

type EmailMessage struct {
	Subject, Content string
	To, From         Email
	Date             time.Time
}

func (e *EmailMessage) ToEmailString() string {
	return string(e.To)
}

func (e *EmailMessage) FromEmailString() string {
	return string(e.From)
}

func (e *EmailMessage) DateString() string {
	return e.Date.Format(time.Stamp)
}

func (e *EmailMessage) String() string {
	return "Email to:" + e.ToEmailString() + " subject:\"" + e.Subject + "\""
}

func NewMessage(subject, content string, to, from Email, date time.Time) *EmailMessage {
	return &EmailMessage{subject, content, to, from, date}
}
