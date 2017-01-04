package contract

import "io"

//RespBody ...
type RespBody io.ReadCloser

//IClient12306 ...
type IClient12306 interface {
	ILogin
}

//ILogin ...
type ILogin interface {
	Login(username, pwd string, vcp IVCode) (bool, error)
}

//IVCode ...
type IVCode interface {
	CaptureVCode() (string, error)
	ResolveVCodeImg(base64Img string) (string, error)
	CheckVCode(code string) (bool, error)
}
