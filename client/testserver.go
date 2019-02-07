package client

import "log"
import "time"

type TestEmailServer struct {
	emails   []*EmailMessage
	incoming chan *EmailMessage

	Outbox []*EmailMessage
}

func (e *TestEmailServer) CurrentMessage() *EmailMessage {
	return e.emails[0]
}

func (e *TestEmailServer) ListMessages() []*EmailMessage {
	return e.emails
}

func (e *TestEmailServer) Send(email *EmailMessage) {
	log.Println("Send", email)
	e.Outbox = append(e.Outbox, email)
}

func (e *TestEmailServer) Incoming() chan *EmailMessage {
	e.incoming = make(chan *EmailMessage)

	return e.incoming
}

func (e *TestEmailServer) simulateArrival() {
	timer := time.NewTimer(time.Second * 10)
	<-timer.C

	newmail := &EmailMessage{
		Subject: "Recently arrived",
		Content: "This email was delivered after the email application loaded.\n\n" +
			"It is just a test email but it arrived.",
		To:   "me@example.com",
		From: "automation@example.com",
		Date: time.Now()}

	e.emails = append([]*EmailMessage{newmail}, e.emails...)

	if e.incoming != nil {
		e.incoming <- newmail
	}
}

func NewTestServer() *TestEmailServer {
	server := &TestEmailServer{emails: []*EmailMessage{
		NewMessage("Testing", "This is an email from our test email server",
			"me@example.com", "you@example.com", time.Now()),
		NewMessage("Older", "This is an old email, it is not the current one",
			"me@example.com", "admin@example.com", time.Now().AddDate(0, 0, -1)),
	},
	}

	go server.simulateArrival()
	return server
}
