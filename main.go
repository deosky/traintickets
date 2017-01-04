package main

import (
	"fmt"
	"time"
	"traintickets/base"
	"traintickets/module/loginmod"
	"traintickets/module/vcodemod"
)

func main() {
	for i := 0; i < 2; i++ {
		vcodeModule := &vcodemod.VCodeModule{}
		loginModule := &loginmod.LoginModule{}
		client, err := base.New12306Client("https://kyfw.12306.cn/otn/login/init", loginModule, vcodeModule)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("初始化成功")
		}
		client.Login()
		time.Sleep(2 * time.Second)
	}

	//client.Login()
}
