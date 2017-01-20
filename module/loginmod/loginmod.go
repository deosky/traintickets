package loginmod

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
	"traintickets/base/appconfig"
	"traintickets/base/contract"
	"traintickets/base/piaohttputil"
)

var (
	appconf       = appconfig.GetAppConfig()
	usernamePat   = `var\s*user_name\s*=\s*'(.*');`
	usernameReg   = regexp.MustCompile(usernamePat)
	userregardPat = `var\s*user_regard\s*=\s*'(.*)';`
	userregardReg = regexp.MustCompile(userregardPat)
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

	urlStr, err := appconfig.Combine(appconf.MainURL, appconf.Ctx, "login/loginAysnSuggest")
	if err != nil {
		return false, nil
	}
	//time.Sleep(2 * time.Second)
	//捕获验证码
	base64Img, err := vcp.CaptureVCode(clientID, "login", "sjrand")
	if err != nil {
		return false, err
	}
	fmt.Println("请输入验证码:")
	var vcode string
	// conn, err := net.Dial("tcp", "127.0.0.1:8686")
	// if err != nil {
	// 	return false, err
	// }
	// v, err := ioutil.ReadAll(conn)
	// conn.Close()
	// if err != nil {
	// 	return false, err
	// }
	// vcode = string(v)

	//fmt.Scanf("%s\n", &vcode)
	//fmt.Printf("输入的验证码为%s\n", vcode)
	vcode, err = vcp.ResolveVCodeImg(clientID, base64Img)
	if err != nil {
		return false, err
	}

	time.Sleep(5 * time.Second)

	_, err = vcp.CheckVCode(clientID, vcode)
	if err != nil {
		return false, err
	}

	time.Sleep(2 * time.Second)

	vs := make(url.Values, 3)
	vs.Add("loginUserDTO.user_name", username)
	vs.Add("userDTO.password", pwd)
	vs.Add("randCode", vcode)

	fmt.Printf("正在登陆 %s , %s , vcode:%s\n", username, "******", vcode)

	referer, _ := appconfig.Combine(appconf.MainURL, appconf.Ctx, "login/init")
	resp, err := piaohttputil.PostV(clientID, urlStr, "application/x-www-form-urlencoded; charset=UTF-8", referer, true, strings.NewReader(vs.Encode()))

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

	// vs1 := make(url.Values, 1)
	// vs1.Add("_json_att=", "")
	//piaohttputil.PostV(clientID, "https://kyfw.12306.cn/otn/login/userLogin", "application/x-www-form-urlencoded", "https://kyfw.12306.cn/otn/login/init", false, strings.NewReader(vs1.Encode()))

	fmt.Println("登陆成功!!!")
	initMy12306UrlStr, _ := appconfig.Combine(appconf.MainURL, appconf.Ctx, "index/initMy12306")
	resp1, err := piaohttputil.Get(clientID, initMy12306UrlStr)
	if err != nil {
		log.Println(err)
		return false, fmt.Errorf("#102,%s", err.Error())
	}
	defer resp1.Body.Close()
	bodydata, err := ioutil.ReadAll(resp1.Body)
	if err != nil {
		return false, fmt.Errorf("解析initMy12306数据失败:%s", err.Error())
	}

	uname, err := getUserName(bodydata)
	if err != nil {
		return false, err
	}
	uregard, _ := getUserregard(bodydata)
	fmt.Println(uname, uregard)

	fmt.Println(resp1.Request.URL.Path)

	return true, nil
}

//Refresh 刷新页面用来保持登录的,如果返回异常则应该重新登录
func (lm *LoginModule) Refresh(clientID int) (bool, error) {

	initMy12306UrlStr, _ := appconfig.Combine(appconf.MainURL, appconf.Ctx, "index/initMy12306")
	resp1, err := piaohttputil.Get(clientID, initMy12306UrlStr)
	if err != nil {
		log.Println(err)
		return false, fmt.Errorf("#179,%s", err.Error())
	}
	defer resp1.Body.Close()

	print(resp1.Request.URL.Path)
	if resp1.Request.URL.Path != "/otn/index/initMy12306" {
		return false, fmt.Errorf("#185,%s", resp1.Request.URL.Path)
	}
	return true, nil
}

//CheckUser ...
func (lm *LoginModule) CheckUser(clientID int) (bool, error) {
	urlStr, _ := appconfig.Combine(appconf.MainURL, appconf.Ctx, "login/checkUser")
	referer, _ := appconfig.Combine(appconf.MainURL, appconf.Ctx, "leftTicket/init")
	vs := make(url.Values, 1)
	vs.Add("_json_att", "")
	resp, err := piaohttputil.PostV(clientID, urlStr, "application/x-www-form-urlencoded; charset=UTF-8", referer, true, strings.NewReader(vs.Encode()))
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

func getUserName(body []byte) (string, error) {
	group := usernameReg.FindSubmatch(body)
	if len(group) < 2 {
		return "", errors.New("获取用户名异常")
	}
	username := string(group[1])

	convertString := "\"" + username + "\""
	s, err := strconv.Unquote(convertString)
	if err != nil {
		return "", err
	}
	return s, nil
}

func getUserregard(body []byte) (string, error) {
	group := userregardReg.FindSubmatch(body)
	if len(group) < 2 {
		return "", errors.New("获取欢迎信息异常")
	}
	userregard := string(group[1])
	convertString := "\"" + userregard + "\""
	s, err := strconv.Unquote(convertString)
	if err != nil {
		return "", err
	}
	return s, nil
}
