package base

import (
	"time"
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

func (client *client12306) QueryTicket(query *contract.TicketQuery) (<-chan *contract.TicketResult, []chan<- bool) {

	return nil, nil
}

// func (client *client12306) QueryTicket(query contract.TicketQuery) error {
// 	ticketMod := client.Context().TicketModule()
// 	return ticketMod.QueryATicket(client.id, &query)
// }

// func (client *client12306) TicketSResult() <-chan (*contract.TicketResult) {
// 	ticketMod := client.Context().TicketModule()
// 	return ticketMod.TicketSResult()
// }

func (client *client12306) CheckOrderInfo(ticket *contract.TicketResult, vcp contract.IVCode, lgm contract.ILogin) (bool, error) {

	return true, nil
}

//Start ...
func (client *client12306) Start() {
	// lgm := client.Context().LoginModule()
	// vcp := client.Context().VCodeModule()
	// f, err := lgm.Login(client.id, "", "", vcp)
	// fmt.Println(f)
	// fmt.Println(err)

	ticketMod := client.Context().TicketModule()
	t, _ := time.Parse("2006-01-02", "2017-02-07")
	query := &contract.TicketQuery{
		TrainDate:    t,
		FromStation:  "FYH",
		ToStation:    "SHH",
		PurposeCodes: "ADULT",
		IntervalTime: 3 * time.Second,
		SeatTypes:    []byte{contract.YW},
	}
	ticketMod.QueryTicket(query)

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
