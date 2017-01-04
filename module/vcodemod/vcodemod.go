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
	ValidateMessages       string                 `json:"validateMessages"`
}
type checkRandCodeAnsynData struct {
	Result string `json:"result"`
	Msg    string `json:"msg"`
}

//CaptureVCode ...touclick-randCode
func (vcode *VCodeModule) CaptureVCode() (string, error) {
	randNum := randGen.Float64()
	vcodeURL := fmt.Sprintf("https://kyfw.12306.cn/otn/passcodeNew/getPassCodeNew?module=login&rand=sjrand&%s", strconv.FormatFloat(randNum, 'f', 17, 64))
	fmt.Println("randnum:=", randNum)
	fmt.Println("vcodeUrl:=", vcodeURL)
	rep, err := piaohttputil.Get(vcodeURL)
	if err != nil {
		fmt.Println("CaptureVCode:=", err)
		return "", err
	}
	defer rep.Body.Close()

	var buf bytes.Buffer
	data := make([]byte, 1024)
	for {
		n, err := rep.Body.Read(data)
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return "", err
			}
		}
		buf.Write(data[:n])
	}
	file, err := os.Create(strconv.FormatFloat(randNum, 'f', 17, 64) + "vcode.png")
	if err != nil {
		return "", err
	}
	defer file.Close()
	file.Write(buf.Bytes())
	base64Str := base64.StdEncoding.EncodeToString(buf.Bytes())
	fmt.Println("data base64 str:=", base64Str)

	md5n := md5.New()
	md5n.Write(buf.Bytes())
	cipherStr := md5n.Sum(nil)
	fmt.Println("md5 bytes:=", cipherStr, "string:=", string(cipherStr), "hex string:=", hex.EncodeToString(cipherStr))

	return base64Str, nil
}

//CheckVCode ...
func (vcode *VCodeModule) CheckVCode(code string) (bool, error) {
	//randCode:110,49,183,45,239,50
	//rand:sjrand
	fmt.Println("正在校验验证码")
	data := make(url.Values, 2)
	data.Add("randCode", code)
	data.Add("rand", "sjrand")

	resp, err := piaohttputil.Post("https://kyfw.12306.cn/otn/passcodeNew/checkRandCodeAnsyn", "application/x-www-form-urlencoded; charset=UTF-8", strings.NewReader(data.Encode()))
	if err != nil {
		return false, err
	}
	fmt.Println("校验完成")
	buf, err := readRespBody(resp.Body)
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
func (vcode *VCodeModule) ResolveVCodeImg(base64Img string) (string, error) {

	return "", nil
}

//readRespBody ...
func readRespBody(resp io.ReadCloser) (*bytes.Buffer, error) {
	buf := &bytes.Buffer{}
	data := make([]byte, 1024)
	for {
		n, err := resp.Read(data)
		buf.Write(data[:n])
		if err != nil {
			if err == io.EOF {
				buf.Write(data[:n])
				break
			} else {
				return buf, err
			}
		}
	}
	return buf, nil
}
