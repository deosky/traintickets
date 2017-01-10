package piaohttputil

import (
	"bytes"
	"io"
	"log"
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
	log.Println("req get:", urlStr)

	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		return nil, err
	}

	return getDo(clientID, req)
}

//GetV ...
func GetV(clientID int, urlStr, referer string, isXhr bool) (*http.Response, error) {
	log.Println("req getv:", urlStr)
	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		return nil, err
	}

	SetReqHeader(req)
	req.Header.Add("Cache-Control", "no-cache")
	req.Header.Add("Referer", referer)
	if isXhr {
		req.Header.Add("X-Requested-With", "XMLHttpRequest")
	}

	return getDo(clientID, req)
}

//getDo ...
func getDo(clientID int, req *http.Request) (*http.Response, error) {
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
	return resp, err
}

//Post ...
func Post(clientID int, url string, bodyType string, body io.Reader) (*http.Response, error) {

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}
	return postDo(clientID, bodyType, req)
}

//PostV ...
func PostV(clientID int, url, bodyType, referer string, isXhr bool, body io.Reader) (*http.Response, error) {

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}

	SetReqHeader(req)

	req.Header.Add("Referer", referer)
	req.Header.Add("Content-Length", string(req.ContentLength))
	log.Println("req.ContentLength", req.ContentLength)
	if isXhr {
		req.Header.Add("X-Requested-With", "XMLHttpRequest")
	}

	return postDo(clientID, bodyType, req)
}

//postDo ...
func postDo(clientID int, bodyType string, req *http.Request) (*http.Response, error) {

	jar, err := jarPool.GetJar(clientID)
	if err != nil {
		return nil, err
	}
	for _, cookie := range jar.Cookies(req.URL) {
		req.AddCookie(cookie)
	}

	req.Header.Set("Content-Type", bodyType)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if rc := resp.Cookies(); len(rc) > 0 {
		jar.SetCookies(req.URL, rc)
	}
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

//SetReqHeader 设置消息头
func SetReqHeader(req *http.Request) {
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Accept-Language", "zh-CN,zh;q=0.8")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/52.0.2743.116 Safari/537.36")
	req.Host = "kyfw.12306.cn"
}
