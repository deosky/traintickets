package vcodemod

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
	"traintickets/base/appconfig"
	"traintickets/base/piaohttputil"
)

//VCodeModule ...
type VCodeModule struct{}

//randGen ...
var (
	randGen = rand.New(rand.NewSource(time.Now().UnixNano()))
)

type checkRandCodeAnsynResult struct {
	ValidateMessagesShowID string                 `json:"validateMessagesShowId"`
	Status                 bool                   `json:"status"`
	Httpstatus             int                    `json:"httpstatus"`
	Data                   checkRandCodeAnsynData `json:"data"`
	Messages               []string               `json:"messages"`
	ValidateMessages       interface{}            `json:"validateMessages"`
}
type checkRandCodeAnsynData struct {
	Result string `json:"result"`
	Msg    string `json:"msg"`
}

var (
	appconf = appconfig.GetAppConfig()
)

//CaptureVCode ...touclick-randCode
func (vcode *VCodeModule) CaptureVCode(clientID int, module, rand string) (*string, error) {
	urlStr, _ := appconfig.Combine(appconf.MainURL, appconf.Ctx, "passcodeNew/getPassCodeNew")
	randNum := randGen.Float64()
	vcodeURL := urlStr + fmt.Sprintf("?module=%s&rand=%s&%s", module, rand, strconv.FormatFloat(randNum, 'f', 17, 64))
	fmt.Println("randnum:=", randNum)
	fmt.Println("vcodeUrl:=", vcodeURL)
	referer, _ := appconfig.Combine(appconf.MainURL, appconf.Ctx, "login/init")
	rep, err := piaohttputil.GetV(clientID, vcodeURL, referer, false)
	if err != nil {
		fmt.Println("CaptureVCode:=", err)
		return nil, err
	}
	defer rep.Body.Close()
	data, err := ioutil.ReadAll(rep.Body)
	if err != nil {
		return nil, err
	}
	file, err := os.Create(strconv.FormatFloat(randNum, 'f', 17, 64) + "vcode.png")
	if err != nil {
		return nil, err
	}
	defer file.Close()
	file.Write(data)
	base64Str := base64.StdEncoding.EncodeToString(data)
	//fmt.Println("data base64 str:=", base64Str)

	md5n := md5.New()
	md5n.Write(data)
	cipherStr := md5n.Sum(nil)
	fmt.Println("md5 bytes:=", cipherStr, "hex string:=", hex.EncodeToString(cipherStr))

	return &base64Str, nil
}

//CheckVCode ...
func (vcode *VCodeModule) CheckVCode(clientID int, code string) (bool, error) {
	//randCode:110,49,183,45,239,50
	//rand:sjrand

	urlStr, _ := appconfig.Combine(appconf.MainURL, appconf.Ctx, "passcodeNew/checkRandCodeAnsyn")
	fmt.Println("正在校验验证码")
	data := make(url.Values, 2)
	data.Add("randCode", code)
	data.Add("rand", "sjrand")
	referer, _ := appconfig.Combine(appconf.MainURL, appconf.Ctx, "login/init")
	resp, err := piaohttputil.PostV(clientID, urlStr, "application/x-www-form-urlencoded; charset=UTF-8", referer, true, strings.NewReader(data.Encode()))
	if err != nil {
		return false, err
	}
	fmt.Println("校验完成")
	buf, err := piaohttputil.ReadRespBody(resp.Body)
	if err != nil {
		return false, err
	}

	var result checkRandCodeAnsynResult
	fmt.Println("返回结果:", buf.String())
	err = json.Unmarshal(buf.Bytes(), &result)
	if err != nil {
		return false, err
	}
	if result.Data.Result == "1" {
		return true, nil
	}

	return false, errors.New("请点击正确的验证码")
}

//ResolveVCodeImg ...
func (vcode *VCodeModule) ResolveVCodeImg(clientID int, base64Img *string) (string, error) {

	vs := make(url.Values, 1)
	vs.Add("data", *base64Img)

	resp, err := piaohttputil.Post(clientID, "http://localhost:8988/getsign", "application/x-www-form-urlencoded; charset=UTF-8", strings.NewReader(vs.Encode()))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	resdata, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	rt := &CheckResult{}
	err = json.Unmarshal(resdata, rt)
	if err != nil {
		return "", err
	}
	if rt.Status == 1 {
		return rt.Data, nil
	}
	return "", errors.New(rt.Data)
}

//CheckResult ...
type CheckResult struct {
	Status int
	Data   string
}
