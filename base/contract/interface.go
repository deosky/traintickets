package contract

//IClient12306 ...
type IClient12306 interface {
	ILogin
	ITicket
}

//ILogin ...
type ILogin interface {
	Login(username, pwd string, vcp IVCode) (bool, error)
}

//IVCode ...
type IVCode interface {
	CaptureVCode(module, rand string) (*string, error)
	ResolveVCodeImg(base64Img *string) (string, error)
	CheckVCode(code string) (bool, error)
}

//ITicket ...
type ITicket interface {
	QueryATicket(query TicketQuery) (*string, error)
}
