package client

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"strings"
	"time"

	mime "github.com/sloonz/go-mime-message"
	"github.com/sloonz/go-qprintable"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"google.golang.org/api/gmail/v1"
)

type gMailServer struct {
	user    string
	srv     *gmail.Service
	history uint64

	emails   []*EmailMessage
	incoming chan *EmailMessage
}

func (g *gMailServer) CurrentMessage() *EmailMessage {
	return g.emails[0]
}

func (g *gMailServer) ListMessages() []*EmailMessage {
	return g.emails
}

func (g *gMailServer) Send(email *EmailMessage) {
	email.From = "test@example.com"
	email.Date = time.Now()

	data := email.ToGMailEncoding()
	msg := &gmail.Message{Raw: data}

	g.srv.Users.Messages.Send(user, msg).Do()
}

func (g *gMailServer) Incoming() chan *EmailMessage {
	in := make(chan *EmailMessage)

	go func() {
		for {
			time.Sleep(10 * time.Second)

			for _, email := range g.downloadNewMessages() {
				g.emails = append([]*EmailMessage{email}, g.emails...)
				in <- email
			}
		}
	}()

	return in
}

func NewGMailServer() EmailServer {
	server := &gMailServer{user: "me"}
	server.srv = setupService()

	server.downloadMessages()

	return server
}

// toGMailEncoding is from Chapter12/cloud/post.go
func (e *EmailMessage) ToGMailEncoding() string {
	m := mime.NewMultipartMessage("alternative", "")
	m.SetHeader("Subject", mime.EncodeWord(e.Subject))
	m.SetHeader("From", mime.EncodeWord("Demo")+" <"+string(e.From)+">")
	m.SetHeader("To", mime.EncodeWord("Demo")+" <"+string(e.To)+">")
	plain := mime.NewTextMessage(qprintable.UnixTextEncoding, bytes.NewBufferString(e.Content))
	plain.SetHeader("Content-Type", "text/plain")
	m.AddPart(plain)

	var buf bytes.Buffer
	io.Copy(&buf, m)
	return base64.URLEncoding.EncodeToString(buf.Bytes())
}

// Helper methods from Google's GMail quickstart
// https://github.com/gsuitedevs/go-samples/blob/master/gmail/quickstart/quickstart.go

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

// end of Google helper functions

// setupService sets up a connection to GMail and configures the credentials
func setupService() *gmail.Service {
	b, err := ioutil.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	config, err := google.ConfigFromJSON(b, gmail.GmailReadonlyScope, gmail.GmailComposeScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(config)

	srv, err := gmail.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve Gmail client: %v", err)
	}

	return srv
}

const user = "me"

// downloadMessages caches all the messages in a user's Inbox
func (g *gMailServer) downloadMessages() {
	req := g.srv.Users.Messages.List(g.user)
	req.LabelIds("INBOX")
	resp, err := req.Do()
	if err != nil {
		log.Fatalf("Unable to retrieve Inbox items: %v", err)
	}

	for _, message := range resp.Messages {
		email := g.downloadMessage(message)
		g.emails = append(g.emails, email)
	}
}

func decodeBody(payload *gmail.MessagePart) string {
	data := payload.Body.Data
	if data == "" {
		data = payload.Parts[0].Body.Data
	}
	content, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		fmt.Println("Failed to decode body", err)
	}

	return string(content)
}

func (g *gMailServer) downloadMessage(message *gmail.Message) *EmailMessage {
	mail, err := g.srv.Users.Messages.Get(user, message.Id).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve message payload: %v", err)
	}
	g.history = uint64(math.Max(float64(g.history), float64(mail.HistoryId)))

	var subject string
	var to, from Email
	var date time.Time

	content := decodeBody(mail.Payload)
	for _, header := range mail.Payload.Headers {
		switch header.Name {
		case "Subject":
			subject = header.Value
		case "To":
			to = Email(header.Value)
		case "From":
			from = Email(header.Value)
		case "Date":
			value := strings.Replace(header.Value, "(UTC)", "", -1)
			date, err = time.Parse("Mon, _2 Jan 2006 15:04:05 -0700", strings.TrimSpace(value))
			if err != nil {
				log.Println("Error: Could not parse date", value)
				date = time.Now()
			}
		}
	}

	return NewMessage(subject, content, to, from, date)
}

// downloadNewMessages caches all messages in a user's Inbox since the most recent one seen
func (g *gMailServer) downloadNewMessages() []*EmailMessage{
	req := g.srv.Users.History.List(g.user)
	req.StartHistoryId(g.history)
	req.LabelId("INBOX")
	resp, err := req.Do()
	if err != nil {
		log.Fatalf("Unable to retrieve Inbox items: %v", err)
	}

	var emails []*EmailMessage
	for _, history := range resp.History {
		for _, message := range history.Messages {
			email := g.downloadMessage(message)
			emails = append(emails, email)
		}
	}

	return emails
}
