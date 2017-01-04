package piaohttputil

import (
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
)

var (
	jar, _ = cookiejar.New(nil)
	client = &http.Client{Jar: jar}
)

//Get ...
func Get(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	fmt.Println("print cookies")
	for _, cookie := range jar.Cookies(req.URL) {
		fmt.Println(cookie)
	}
	fmt.Println("print cookies end")
	return resp, err
}

//Post ...
func Post(url string, bodyType string, body io.Reader) (*http.Response, error) {

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", bodyType)
	resp, err := client.Do(req)
	fmt.Println("print cookies")
	for _, cookie := range jar.Cookies(req.URL) {
		fmt.Println(cookie)
	}
	fmt.Println("print cookies end")
	return resp, err
}
