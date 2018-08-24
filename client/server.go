package client

import "log"
import "time"

type EmailServer struct {
	emails []*EmailMessage
}

func (e *EmailServer) CurrentMessage() *EmailMessage {
	return e.emails[0]
}

func (e *EmailServer) ListMessages() []*EmailMessage {
	return e.emails
}

func (e *EmailServer) Send(email *EmailMessage) {
	log.Println("Send email: ", email)
}

func NewTestServer() *EmailServer {
	return &EmailServer{
		[]*EmailMessage{
			NewMessage("Testing", "This is an email from our test email server",
				"me@example.com", "you@example.com", time.Now()),
			NewMessage("Older", "This is an old email, it is not the current one",
				"me@example.com", "admin@example.com", time.Now().AddDate(0, 0, -1)),
		},
	}
}
