package base

import (
	"fmt"
	"log"
	"traintickets/base/appconfig"
	"traintickets/base/contract"
	"traintickets/base/piaohttputil"
)

var (
	appconf = appconfig.GetAppConfig()
)

//ClientContext ...
type ClientContext struct {
	loginModule  contract.ILogin
	vcodeModule  contract.IVCode
	ticketModule contract.ITicket
}

//LoginModule ...
func (cc *ClientContext) LoginModule() contract.ILogin {
	return cc.loginModule
}

//VCodeModule ...
func (cc *ClientContext) VCodeModule() contract.IVCode {
	return cc.vcodeModule
}

// TicketModule ...
func (cc *ClientContext) TicketModule() contract.ITicket {
	return cc.ticketModule
}

//NewClientContext ...
func NewClientContext(login contract.ILogin, vcode contract.IVCode, ticket contract.ITicket) contract.IClientContext {
	return &ClientContext{loginModule: login, vcodeModule: vcode, ticketModule: ticket}
}

//Client12306 ...
type client12306 struct {
	id      int
	context contract.IClientContext
}

func (client *client12306) Context() contract.IClientContext {
	return client.context
}

//Start 开始刷票
func (client *client12306) Start(query *contract.TicketQuery) {

	username := ""
	pwd := ""

	fmt.Println("12306用户名:")
	fmt.Scanf("%s\n", &username)
	fmt.Println("12306密码:")

	lgm := client.Context().LoginModule()
	vcp := client.Context().VCodeModule()
	_, err := lgm.Login(client.id, username, pwd, vcp)
	if err != nil {
		log.Println("登陆失败:", err.Error())
		return
	}
	log.Println("登陆成功")

	log.Println("开始自动刷票中")
	ticketMod := client.Context().TicketModule()
	log.Println("查票参数:", query)
	res := ticketMod.QueryTicket(query)

	//开始下单
	for t := range res {
		ck := &contract.CheckOutOrderContext{
			VCodeMod:          client.Context().VCodeModule(),
			LoginMod:          client.Context().LoginModule(),
			UserName:          username,
			Pwd:               pwd,
			PassengerIDCardNo: []string{""},
			SecretStr:         t.SecretStr,
			Train: contract.TrainInfo{
				StationTrainCode:     t.StationTrainCode,
				TrainDate:            t.TrainDate,
				BackTrainDate:        t.BackTrainDate,
				TourFlag:             t.TourFlag,
				PurposeCodes:         t.PurposeCodes,
				QueryFromStationName: t.QueryFromStationName,
				QueryToStationName:   t.QueryToStationName,
			},
			SeatType:   t.SeatTypes,
			TicketType: contract.TICKETTYPEADULT,
		}
		cf, err := ticketMod.CheckOutOrder(0, ck)
		if err != nil {
			log.Println("chekout err :", err.Error())
			ticketMod.ReStart()
			continue
		}
		if cf {
			break
		}
		log.Println("chekout err")
	}
	log.Println("checkout success.")

}

//New12306Client ...
func New12306Client(context contract.IClientContext, urlStr string) (contract.IClient12306, error) {
	clientid := 0
	referer, err := appconfig.Combine(appconf.MainURL, appconf.Ctx)
	if err != nil {
		return nil, fmt.Errorf("mainUrl error:%s", err.Error())
	}
	resp, err := piaohttputil.GetV(clientid, urlStr, referer, false)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	client := client12306{id: clientid, context: context}
	return &client, nil
}
