package main

import "encoding/json"
import "fmt"
import "io/ioutil"
import "net/http"

type Person struct {
	Title     string `json:"title,omitempty"`
	Firstname string `json:"firstname"`
	Surname   string `json:"surname"`

	Username string `json:"username"`
	Password string `json:"-"`
}

func readFromURL(url string) ([]byte, error) {
	var body []byte
	resp, err := http.Get(url)
	if err != nil {
		return body, err
	}

	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

func main() {
	person := &Person{
		"",
		"John",
		"Doe",
		"someuser",
		"somepassword",
	}
	fmt.Println("Struct:", person)

	data, _ := json.MarshalIndent(person, "", "  ")
	fmt.Println("JSON:", string(data))

	fmt.Println("Downloading...")
	data, _ = readFromURL("http://echo.jsontest.com/title/Sir/" +
		"firstname/Tom/surname/Jones/username/singer1/")
	fmt.Println("Download:", string(data))

	person = &Person{}
	json.Unmarshal(data, person)
	fmt.Println("Decoded:", person)
}
