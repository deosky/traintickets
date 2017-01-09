package main

import (
	"fmt"
	"traintickets/base"
	"traintickets/module/loginmod"
	"traintickets/module/piaomod"
	"traintickets/module/vcodemod"
)

func main() {
	vcodeModule := &vcodemod.VCodeModule{}
	loginModule := &loginmod.LoginModule{}
	ticketModule := &piaomod.PIAO{}

	sign := make(chan bool, 1)

	mainContext := base.NewClientContext(loginModule, vcodeModule, ticketModule)

	client, err := base.New12306Client(mainContext, "https://kyfw.12306.cn/otn/login/init")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("初始化成功")
	}

	client.Start()

	<-sign
	// a := 0
	// go func() {

	// 	for i := 0; i < 10; i++ {
	// 		go func() {

	// 			a++
	// 			fmt.Println(a)
	// 		}()
	// 	}
	// 	//time.Sleep(2 * time.Second)
	// 	fmt.Print("结束\n")
	// }()

	// time.Sleep(10 * time.Second)

}

func main1() {

	// vcodeModule := &vcodemod.VCodeModule{}
	// loginModule := &loginmod.LoginModule{}
	// ticketModule := &piaomod.PIAO{}

	// mainContext := base.NewClientContext(loginModule, vcodeModule, ticketModule)

	// client, err := base.New12306Client(mainContext, "https://kyfw.12306.cn/otn/login/init")
	// if err != nil {
	// 	fmt.Println(err)
	// } else {
	// 	fmt.Println("初始化成功")
	// }
	// t, _ := time.Parse("2006-01-02", "2017-02-02")
	// query := contract.TicketQuery{TrainDate: t, FromStation: "FYH", ToStation: "SHH", PurposeCodes: "ADULT"}

	// stopsign := make(chan bool, 1)
	// ticker := time.NewTicker(time.Second * 3)
	// go func() {
	// 	for {
	// 		// f := false
	// 		// err = client.QueryTicket(query)
	// 		// if err != nil {
	// 		// 	println("query ticket err :", err.Error())
	// 		// }
	// 		// select {
	// 		// case <-stopsign:
	// 		// 	f = true
	// 		// case <-ticker.C:
	// 		// }
	// 		// if f {
	// 		// 	break
	// 		// }
	// 	}
	// 	fmt.Println("刷票流程结束")
	// }()
	// go func() {
	// 	// r := client.TicketSResult()
	// 	// for t := range r {
	// 	// 	stopsign <- true
	// 	// 	fmt.Println("抢到票啦啦啦啦！！！！", "车次:", t.StationTrainCode)
	// 	// }
	// }()

	// //client.Login("adminadmin", "adminadmin", vcodeModule)
	// //query := contract.TicketQuery{TrainDate: time.Now(), FromStation: "SHH", ToStation: "BJP", PurposeCodes: "ADULT"}
	// //_, err = client.QueryATicket(query)
	// //fmt.Println("error:", err)
	// fmt.Println("输入exit退出程序")
	// for {
	// 	s := ""
	// 	fmt.Scanln(&s)
	// 	if s == "exit" {
	// 		break
	// 	}
	// }

}
