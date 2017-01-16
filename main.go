package main

import (
	"fmt"
	"time"
	"traintickets/base"
	"traintickets/base/appconfig"
	"traintickets/base/contract"
	"traintickets/module/loginmod"
	"traintickets/module/ticketmod"
	"traintickets/module/vcodemod"
)

var (
	appconf = appconfig.GetAppConfig()
)

func main() {
	vcodeModule := &vcodemod.VCodeModule{}
	loginModule := &loginmod.LoginModule{}
	ticketModule := &ticketmod.PIAO{}

	sign := make(chan bool, 1)

	mainContext := base.NewClientContext(loginModule, vcodeModule, ticketModule)

	client, err := base.New12306Client(mainContext, appconf.InitURL)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("初始化成功")
	}

	t, _ := time.Parse("2006-01-02", "2017-02-09")
	query := &contract.TicketQuery{
		TrainDate:    t,
		FromStation:  "FYH",
		ToStation:    "SHH",
		PurposeCodes: "ADULT",
		IntervalTime: 3 * time.Second,
		SeatTypes:    []byte{contract.YW},
	}
	client.Start(query)

	<-sign

}
