package appconfig

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path"
	"path/filepath"
)

//AppConfig ...
type AppConfig struct {
	CdnIPAddr []string `json:"cdnIpAddr"`
	InitURL   string   `json:"initUrl"`
	MainURL   string   `json:"mainUrl"`
	Ctx       string   `json:"ctx"`
}

var (
	config = &AppConfig{}
)

func init() {
	os.Getwd()
	dir, err := filepath.Abs("app.json")
	if err != nil {
		log.Fatal("获取app.json路径失败:", err)
	}
	file, err := os.Open(dir)
	defer file.Close()
	if err != nil {
		log.Fatal("打开文件app.json失败:", err)
	}
	configData, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal("读取文件app.json失败:", err)
	}
	conf := &AppConfig{}
	err = json.Unmarshal(configData, conf)
	if err != nil {
		log.Fatal("反序列化文件app.json失败:", err)
	}
	config = conf

}

//GetAppConfig ...
func GetAppConfig() *AppConfig {
	return config
}

//Combine ...
func Combine(mainurl string, ps ...string) (string, error) {
	u, err := url.Parse(mainurl)
	if err != nil {
		return "", nil
	}
	p := ""
	for _, v := range ps {
		p = path.Join(p, v)
	}
	u.Path = p

	return u.String(), nil
}

// {
//     "cdnIpAddr":[
//         "61.155.162.122",
//         "211.144.7.86",
//         "183.134.10.85"
//     ],
//     "initUrl":"https://kyfw.12306.cn/otn/login/init",
//     "mainUrl":"https://kyfw.12306.cn/",
//     "ctx":"otn"
// }
