package base

import (
	"traintickets/base/contract"
	"traintickets/base/piaohttputil"
)

//Client12306 ...
type client12306 struct {
	loginModule contract.ILogin
}

//New12306Client ...
func New12306Client(url string, login contract.ILogin) (contract.IClient12306, error) {

	resp, err := piaohttputil.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	client := client12306{loginModule: login}
	return &client, nil
}

func (client *client12306) Login(username, pwd string, vcp contract.IVCode) (bool, error) {
	return client.loginModule.Login(username, pwd, vcp)
}
