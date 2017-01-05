package main

import (
	"fmt"
	"time"
	"traintickets/base"
	"traintickets/base/contract"
	"traintickets/module/loginmod"
	"traintickets/module/piaomod"
)

func main() {

	//vcodeModule := &vcodemod.VCodeModule{}
	loginModule := &loginmod.LoginModule{}
	ticketModule := &piaomod.PIAO{}

	client, err := base.New12306Client("https://kyfw.12306.cn/otn/login/init", loginModule, ticketModule)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("初始化成功")
	}
	//client.Login("adminadmin", "adminadmin", vcodeModule)
	query := contract.TicketQuery{TrainDate: time.Now(), FromStation: "SHH", ToStation: "BJP", PurposeCodes: "ADULT"}
	_, err = client.QueryATicket(query)
	fmt.Println("error:", err)
}
