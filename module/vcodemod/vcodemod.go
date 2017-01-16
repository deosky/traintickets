package vcodemod

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
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
	rep, err := piaohttputil.Get(clientID, vcodeURL)
	if err != nil {
		fmt.Println("CaptureVCode:=", err)
		return nil, err
	}
	defer rep.Body.Close()

	var buf bytes.Buffer
	data := make([]byte, 1024)
	for {
		n, err := rep.Body.Read(data)
		buf.Write(data[:n])
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return nil, err
			}
		}
	}
	file, err := os.Create(strconv.FormatFloat(randNum, 'f', 17, 64) + "vcode.png")
	if err != nil {
		return nil, err
	}
	defer file.Close()
	file.Write(buf.Bytes())
	base64Str := base64.StdEncoding.EncodeToString(buf.Bytes())
	//fmt.Println("data base64 str:=", base64Str)

	md5n := md5.New()
	md5n.Write(buf.Bytes())
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

	resp, err := piaohttputil.Post(clientID, urlStr, "application/x-www-form-urlencoded; charset=UTF-8", strings.NewReader(data.Encode()))
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

	return "", nil
}
