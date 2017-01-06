package loginmod

import (
	"fmt"
	"traintickets/base/contract"
)

//LoginModule ...
type LoginModule struct{}

//Login ...
func (lm *LoginModule) Login(clientID int, username, pwd string, vcp contract.IVCode) (bool, error) {
	//捕获验证码
	_, err := vcp.CaptureVCode(clientID, "login", "sjrand")
	if err != nil {
		return false, err
	}
	fmt.Println("请输入验证码:")
	var vcode string
	fmt.Scanf("%s", &vcode)
	fmt.Printf("输入的验证码为%s\n", vcode)
	vcp.CheckVCode(clientID, vcode)
	//fmt.Scanln(vcode)
	return true, nil
}
