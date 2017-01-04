package piaohttputil

import (
	"io"
	"net/http"
	"net/http/cookiejar"
)

var (
	gCookieJar, _ = cookiejar.New(nil)
	client        = &http.Client{}
)

//Get ...
func Get(url string) (*http.Response, error) {

	// req, err := http.NewRequest("GET", url, nil)
	// if err != nil {
	// 	return nil, err
	// }
	client.Jar = gCookieJar
	return client.Get(url)
	//return http.Get(url)
}

//Post ...
func Post(url string, bodyType string, body io.Reader) (resp *http.Response, err error) {
	return http.Post(url, bodyType, body)
}
