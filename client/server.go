package client

type EmailServer interface {
	CurrentMessage() *EmailMessage
	ListMessages() []*EmailMessage
	Send(email *EmailMessage)
	Incoming() chan *EmailMessage
}
