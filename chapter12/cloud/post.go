package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	mime "github.com/sloonz/go-mime-message"
	"github.com/sloonz/go-qprintable"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type Email string

type EmailMessage struct {
	Subject string
	Content string
	To      Email
	From    Email
	Date    time.Time
}

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

func (e *EmailMessage) String() string {
	format := "EmailMessage{\n  To:%s\n  From:%s\n  Subject:%s\n  Date:%s\n}"
	return fmt.Sprintf(format, e.To, e.From, e.Subject, e.Date.String())
}

var conf *oauth2.Config
var ctx = context.Background()
var client = &http.Client{}
var token *oauth2.Token

const TokenPath = "/tmp/gomail.token"

func loadToken() *oauth2.Token {
	if _, err := os.Stat(TokenPath); os.IsNotExist(err) {
		return nil
	}

	tok := &oauth2.Token{}
	file, _ := os.Open(TokenPath)
	enc := json.NewDecoder(file)
	enc.Decode(&tok)

	return tok
}

func authPost(url, mime string, body io.Reader) io.ReadCloser {
	if token != nil {
		fmt.Println("Reusing stored token")
		client = conf.Client(ctx, token)
	}
	ret, err := client.Post(url, mime, body)
	if err != nil || ret.StatusCode == 401 {
		log.Fatal("Could not authenticate to send", err)
	}

	return ret.Body
}

func setupOAuth() {
	// Your credentials should be obtained from the Google Developer Console
	// (https://console.developers.google.com).
	conf = &oauth2.Config{
		ClientID:     "CLIENT_ID",
		ClientSecret: "CLIENT_SECRET",
		Scopes:       []string{"https://www.googleapis.com/auth/gmail.compose"},
		Endpoint:     google.Endpoint,
		RedirectURL:  "http://127.0.0.1:19999/oauth/callback",
	}
	ctx = context.WithValue(context.Background(), oauth2.HTTPClient, client)
}

func postMessage(msg *EmailMessage) {
	raw := msg.ToGMailEncoding()
	body := fmt.Sprintf("{\"raw\": \"%s\"}", raw)

	ret := authPost("https://www.googleapis.com/gmail/v1/users/me/messages/send",
		"application/json", strings.NewReader(body))
	io.Copy(os.Stdout, ret)
	ret.Close()
}

func main() {
	setupOAuth()
	token = loadToken()

	msg := &EmailMessage{
		"GoMail Test Email",
		"This is a test email sent from a Go example",
		"andy@andywilliams.me",
		"handyande@gmail.com",
		time.Now()}
	postMessage(msg)

	time.Sleep(time.Second)
}
