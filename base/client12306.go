package base

import (
	"log"
	"traintickets/base/contract"
	"traintickets/base/piaohttputil"
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
	lgm := client.Context().LoginModule()
	vcp := client.Context().VCodeModule()
	_, err := lgm.Login(client.id, "", "", vcp)
	if err != nil {
		log.Println("登陆失败:", err.Error())
		return
	}
	log.Println("登陆成功")

	log.Println("开始自动刷票中")
	ticketMod := client.Context().TicketModule()
	// t, _ := time.Parse("2006-01-02", "2017-02-07")
	// query := &contract.TicketQuery{
	// 	TrainDate:    t,
	// 	FromStation:  "FYH",
	// 	ToStation:    "SHH",
	// 	PurposeCodes: "ADULT",
	// 	IntervalTime: 3 * time.Second,
	// 	SeatTypes:    []byte{contract.YW},
	// }
	log.Println("查票参数:", query)
	res := ticketMod.QueryTicket(query)

	//开始下单
	for t := range res {
		ck := &contract.CheckOutOrderContext{
			// 		Mod               CheckOrderMod
			// UserName          string
			// Pwd               string
			// PassengerIDCardNo []string
			SecretStr: t.SecretStr,
			// Train             TrainInfo
			// SeatType          string
			// TicketType        string
		}
		cf, err := ticketMod.CheckOutOrder(0, ck)
		if err != nil {
			log.Println("chekout err :", err.Error())
			continue
		}
		if cf {
			break
		}
		log.Println("chekout err")
	}

}

//New12306Client ...
func New12306Client(context contract.IClientContext, urlStr string) (contract.IClient12306, error) {
	clientid := 1
	resp, err := piaohttputil.Get(clientid, urlStr)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	client := client12306{id: clientid, context: context}
	return &client, nil
}
