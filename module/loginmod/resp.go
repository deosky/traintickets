package loginmod

//RespHead ...
type RespHead struct {
	ValidateMessagesShowID string      `json:"validateMessagesShowId"`
	Status                 bool        `json:"status"`
	Httpstatus             int         `json:"httpstatus"`
	Messages               []string    `json:"messages"`
	ValidateMessages       interface{} `json:"validateMessages"`
}

type loginResp struct {
	RespHead
	Data loginRespData `json:"data"`
}

type loginRespData struct {
	LoginAddress string `json:"loginAddress"`
	OtherMsg     string `json:"otherMsg"`
	LoginCheck   string `json:"loginCheck"`
}
