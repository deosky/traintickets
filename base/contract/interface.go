package contract

//IClient12306 ...
type IClient12306 interface {
	Context() IClientContext
	QueryTicket(query TicketQuery) error
	TicketSResult() <-chan (*TicketResult)
}

//IClientContext ...
type IClientContext interface {
	LoginModule() ILogin
	VCodeModule() IVCode
	TicketModule() ITicket
}

//ILogin ...
type ILogin interface {
	Login(clientID int, username, pwd string, vcp IVCode) (bool, error)
}

//IVCode ...
type IVCode interface {
	CaptureVCode(clientID int, module, rand string) (*string, error)
	ResolveVCodeImg(clientID int, base64Img *string) (string, error)
	CheckVCode(clientID int, code string) (bool, error)
}

//ITicket ...
type ITicket interface {
	QueryATicket(clientID int, query *TicketQuery) error
	TicketSResult() <-chan (*TicketResult)
}
