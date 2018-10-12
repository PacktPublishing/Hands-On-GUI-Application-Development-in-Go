package main

import (
	"context"
	"encoding/json"
	"fmt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"runtime"
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

func (e *EmailMessage) String() string {
	format := "EmailMessage{\n  To:%s\n  From:%s\n  Subject:%s\n  Date:%s\n}"
	return fmt.Sprintf(format, e.To, e.From, e.Subject, e.Date.String())
}

func openbrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unknown platform")
	}
	if err != nil {
		log.Fatal(err)
	}
}

var conf *oauth2.Config
var ctx = context.Background()
var client = &http.Client{}
var retReader io.ReadCloser
var server *http.Server
var token *oauth2.Token

const TokenPath = "/tmp/gomail.token"

func saveToken(tok *oauth2.Token) {
	file, _ := os.Create(TokenPath)
	enc := json.NewEncoder(file)
	enc.Encode(tok)
}

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

func callbackHandler(w http.ResponseWriter, r *http.Request) {
	queryParts, _ := url.ParseQuery(r.URL.RawQuery)

	tok, err := conf.Exchange(ctx, queryParts["code"][0])
	if err != nil {
		log.Fatal(err)
	}
	saveToken(tok)
	client = conf.Client(ctx, tok)

	ret, err := client.Get("https://www.googleapis.com/gmail/v1/users/me/messages")
	if err != nil {
		fmt.Fprint(w, "<p><strong>Authentication Failed</strong></p>")
		fmt.Fprintf(w, "<p style=\"color: red\">%s</p>", err.Error())
		fmt.Fprint(w, "<p>Please close this window and try again.</p>")
		log.Fatal(err)
	} else {
		fmt.Fprint(w, "<p><strong>Authentication Completed</strong></p>")
		fmt.Fprint(w, "<p>Please close this window.</p>")

		retReader = ret.Body
	}

	if f, ok := w.(http.Flusher); ok {
		f.Flush()
	}

	server.Shutdown(context.Background())
}

func authStream(url string) io.ReadCloser {
	if token != nil {
		fmt.Println("Reusing stored token")
		client = conf.Client(ctx, token)
	}
	ret, err := client.Get(url)
	if err == nil && ret.StatusCode != 401 {
		return ret.Body
	}

	fmt.Println("Requesting authorization")
	openbrowser(conf.AuthCodeURL("state", oauth2.AccessTypeOffline))

	http.HandleFunc("/oauth/callback", callbackHandler)
	server = &http.Server{Addr: ":19999", Handler: nil}
	server.ListenAndServe()

	return retReader
}

func setupOAuth() {
	// Your credentials should be obtained from the Google Developer Console
	// (https://console.developers.google.com).
	conf = &oauth2.Config{
		ClientID:     "CLIENT_ID",
		ClientSecret: "CLIENT_SECRET",
		Scopes: []string{"https://www.googleapis.com/auth/gmail.readonly",
			"https://www.googleapis.com/auth/gmail.compose"},
		Endpoint:    google.Endpoint,
		RedirectURL: "http://127.0.0.1:19999/oauth/callback",
	}
	ctx = context.WithValue(context.Background(), oauth2.HTTPClient, client)
}

func countMessages() {
	in := authStream("https://www.googleapis.com/gmail/v1/users/me/messages")
	defer in.Close()

	var content interface{}
	decoder := json.NewDecoder(in)
	decoder.Decode(&content)

	if body, ok := content.(map[string]interface{}); ok {
		list := body["messages"].([]interface{})
		fmt.Println(len(list), "messages found")
	}
}

func main() {
	setupOAuth()
	token = loadToken()

	countMessages()

	time.Sleep(time.Second)
}
