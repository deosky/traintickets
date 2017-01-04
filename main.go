package main

import (
	"fmt"
	"traintickets/base"
	"traintickets/module/loginmod"
	"traintickets/module/vcodemod"
)

func main() {

	vcodeModule := &vcodemod.VCodeModule{}
	loginModule := &loginmod.LoginModule{}

	client, err := base.New12306Client("https://kyfw.12306.cn/otn/login/init", loginModule)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("初始化成功")
	}
	client.Login("adminadmin", "adminadmin", vcodeModule)
	
}
