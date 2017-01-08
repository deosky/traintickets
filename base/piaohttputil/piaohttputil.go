package piaohttputil

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"sync"
)

var (
	client  = &http.Client{}
	jarPool = &cJar{jars: make(map[int]*cookiejar.Jar)}
)

//cJar ...
type cJar struct {
	jars map[int]*cookiejar.Jar
	mux  sync.Mutex
}

func (c *cJar) GetJar(id int) (*cookiejar.Jar, error) {
	c.mux.Lock()
	defer c.mux.Unlock()

	if _, ok := c.jars[id]; !ok {
		jar, err := cookiejar.New(nil)
		if err != nil {
			return nil, err
		}
		c.jars[id] = jar
	}
	return c.jars[id], nil
}

//Get ...
func Get(clientID int, urlStr string) (*http.Response, error) {
	fmt.Println("req get:", urlStr)

	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		return nil, err
	}
	jar, err := jarPool.GetJar(clientID)
	if err != nil {
		return nil, err
	}
	for _, cookie := range jar.Cookies(req.URL) {
		req.AddCookie(cookie)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if rc := resp.Cookies(); len(rc) > 0 {
		jar.SetCookies(req.URL, rc)
	}
	fmt.Println("print cookies")
	for _, cookie := range jar.Cookies(req.URL) {
		fmt.Println(cookie)
	}
	fmt.Println("print cookies end")

	return resp, err
}

//Post ...
func Post(clientID int, url string, bodyType string, body io.Reader) (*http.Response, error) {

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}
	return Do(clientID, bodyType, req)
}

//Do ...
func Do(clientID int, bodyType string, req *http.Request) (*http.Response, error) {

	jar, err := jarPool.GetJar(clientID)
	if err != nil {
		return nil, err
	}
	for _, cookie := range jar.Cookies(req.URL) {
		req.AddCookie(cookie)
	}

	req.Header.Set("Content-Type", bodyType)

	// u, _ := url.Parse("https://192.168.0.2:8888")
	// tr := &http.Transport{
	// 	Proxy: http.ProxyURL(u),
	// }
	// client.Transport = tr
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if rc := resp.Cookies(); len(rc) > 0 {
		jar.SetCookies(req.URL, rc)
	}
	fmt.Println("print cookies")
	for _, cookie := range jar.Cookies(req.URL) {
		fmt.Println(cookie)
	}
	fmt.Println("print cookies end")
	return resp, err
}

//ReadRespBody ...
func ReadRespBody(resp io.ReadCloser) (*bytes.Buffer, error) {
	buf := &bytes.Buffer{}
	data := make([]byte, 1024)
	for {
		n, err := resp.Read(data)
		buf.Write(data[:n])
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return buf, err
			}
		}
	}
	return buf, nil
}
