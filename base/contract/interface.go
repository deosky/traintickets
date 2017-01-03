package contract

import "io"

//RespBody ...
type RespBody io.ReadCloser

//IClient12306 ...
type IClient12306 interface {
	ILogin
	IVCode
}

//ILogin ...
type ILogin interface {
	Login() error
}

//IVCode ...
type IVCode interface {
	CaptureVCode(resp RespBody) (string, error)
}
