package base

import (
	"fmt"
	"net/http"
	"traintickets/base/contract"
)

//Client12306 ...
type client12306 struct {
	cookies     []*http.Cookie
	loginModule contract.ILogin
	vcodeModule contract.IVCode
}

//New12306Client ...
func New12306Client(url string, login contract.ILogin, vcode contract.IVCode) (contract.IClient12306, error) {

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	client := client12306{loginModule: login, vcodeModule: vcode}
	cookieCount := len(resp.Cookies())
	client.cookies = make([]*http.Cookie, cookieCount)
	copy(client.cookies, resp.Cookies())
	fmt.Println("1:", resp.Cookies())
	fmt.Println("2:", client.cookies)
	fmt.Println("3:", resp.Header)
	client.CaptureVCode(resp.Body)
	//
	return &client, nil
}

func (client *client12306) Login() error {
	return client.loginModule.Login()
}

func (client *client12306) CaptureVCode(resp contract.RespBody) (string, error) {
	return client.vcodeModule.CaptureVCode(resp)
}
