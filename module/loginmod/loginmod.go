package loginmod

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/url"
	"strings"
	"time"
	"traintickets/base/contract"
	"traintickets/base/piaohttputil"
)

//checkUserResult ...
type checkUserResult struct {
	ValidateMessagesShowID string        `json:"validateMessagesShowId"`
	Status                 bool          `json:"status"`
	Httpstatus             int           `json:"httpstatus"`
	Data                   checkUserData `json:"data"`
	Messages               []string      `json:"messages"`
	ValidateMessages       interface{}   `json:"validateMessages"`
}

//checkUserData ...
type checkUserData struct {
	Flag bool `json:"flag"`
}

//LoginModule ...
type LoginModule struct{}

//Login ...
func (lm *LoginModule) Login(clientID int, username, pwd string, vcp contract.IVCode) (bool, error) {
	urlStr := "https://kyfw.12306.cn/otn/login/loginAysnSuggest"

	time.Sleep(2 * time.Second)
	//捕获验证码
	_, err := vcp.CaptureVCode(clientID, "login", "sjrand")
	if err != nil {
		return false, err
	}
	fmt.Println("请输入验证码:")
	var vcode string
	conn, err := net.Dial("tcp", "127.0.0.1:8686")
	if err != nil {
		return false, err
	}
	v, err := ioutil.ReadAll(conn)
	conn.Close()
	if err != nil {
		return false, err
	}
	vcode = string(v)

	// fmt.Scanf("%s", &vcode)
	fmt.Printf("输入的验证码为%s\n", vcode)
	_, err = vcp.CheckVCode(clientID, vcode)
	if err != nil {
		return false, err
	}

	time.Sleep(5 * time.Second)

	vs := make(url.Values, 3)
	vs.Add("loginUserDTO.user_name", username)
	vs.Add("userDTO.password", pwd)
	vs.Add("randCode", vcode)

	fmt.Printf("正在登陆 %s , %s , vcode:%s\n", username, pwd, vcode)
	resp, err := piaohttputil.PostV(clientID, urlStr, "application/x-www-form-urlencoded; charset=UTF-8", "https://kyfw.12306.cn/otn/login/init", true, strings.NewReader(vs.Encode()))

	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return false, fmt.Errorf("Error StatusCode %d", resp.StatusCode)
	}
	fmt.Printf("开始解析结果")
	buf, err := piaohttputil.ReadRespBody(resp.Body)
	if err != nil {
		return false, err
	}
	fmt.Printf("结果:%s\n", buf.String())
	res := &loginResp{}
	err = json.Unmarshal(buf.Bytes(), res)
	if err != nil {
		return false, err
	}

	if res.Data.LoginCheck != "Y" {
		return false, fmt.Errorf("%v\r\n%v", res.Messages, res.Data)
	}
	vs1 := make(url.Values, 1)
	vs1.Add("_json_att=", "")

	//piaohttputil.PostV(clientID, "https://kyfw.12306.cn/otn/login/userLogin", "application/x-www-form-urlencoded", "https://kyfw.12306.cn/otn/login/init", false, strings.NewReader(vs1.Encode()))

	fmt.Println("登陆成功!!!")

	resp1, _ := piaohttputil.Get(clientID, "https://kyfw.12306.cn/otn/index/initMy12306")
	defer resp.Body.Close()
	buf, err = piaohttputil.ReadRespBody(resp1.Body)
	fmt.Println("登陆页面展示")
	fmt.Println(buf.String())

	return true, nil
}

//CheckUser ...
func (lm *LoginModule) CheckUser(clientID int) (bool, error) {

	vs := make(url.Values, 1)
	vs.Add("_json_att", "")
	resp, err := piaohttputil.Post(clientID, "https://kyfw.12306.cn/otn/login/checkUser", "application/x-www-form-urlencoded; charset=UTF-8", strings.NewReader(vs.Encode()))
	if err != nil {
		return false, err
	}
	if resp.StatusCode != 200 {
		return false, fmt.Errorf("CheckUser错误的http状态:%d", resp.StatusCode)
	}

	defer resp.Body.Close()
	buf, err := piaohttputil.ReadRespBody(resp.Body)
	if err != nil {
		return false, fmt.Errorf("CheckUser:%s", err.Error())
	}
	rs := &checkUserResult{}
	err = json.Unmarshal(buf.Bytes(), rs)
	if err != nil {
		return false, err
	}

	return rs.Data.Flag, nil
}
