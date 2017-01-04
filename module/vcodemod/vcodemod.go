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
	"os"
	"strconv"
	"time"
	"traintickets/base/contract"
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
func (vcode *VCodeModule) CaptureVCode(resp contract.RespBody) (string, error) {
	//0.5872982159059681
	//https://kyfw.12306.cn/otn/passcodeNew/getPassCodeNew?module=login&rand=sjrand&0.12693779092491142
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
	fmt.Println("data string:=", buf.String())
	fmt.Println("data base64 str:=", base64.StdEncoding.EncodeToString(buf.Bytes()))
	md5n := md5.New()
	md5n.Write(buf.Bytes())
	cipherStr := md5n.Sum(nil)
	fmt.Println("md5 bytes:=", cipherStr, "string:=", string(cipherStr), "hex string:=", hex.EncodeToString(cipherStr))

	return "", nil
}

//CheckVCode ...
func (vcode *VCodeModule) CheckVCode(code string) (bool, error) {

	resp, err := piaohttputil.Get("https://kyfw.12306.cn/otn/passcodeNew/checkRandCodeAnsyn")
	buf, err := readRespBody(resp.Body)
	if err != nil {
		return false, err
	}
	var result checkRandCodeAnsynResult
	err = json.Unmarshal(buf.Bytes(), &result)
	if err != nil {
		return false, err
	}
	if result.Data.Result == "1" {
		return true, nil
	}

	return false, errors.New("请点击正确的验证码")
}

//readRespBody ...
func readRespBody(resp io.ReadCloser) (*bytes.Buffer, error) {
	buf := &bytes.Buffer{}
	data := make([]byte, 1024)
	for {
		n, err := resp.Read(data)
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return buf, err
			}
		}
		buf.Write(data[:n])
	}
	return buf, nil
}
