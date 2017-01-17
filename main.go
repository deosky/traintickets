package main

import (
	"fmt"
	"syscall"
	"time"
	"traintickets/base"
	"traintickets/base/appconfig"
	"traintickets/base/contract"
	"traintickets/module/loginmod"
	"traintickets/module/ticketmod"
	"traintickets/module/vcodemod"

	"log"

	"golang.org/x/crypto/ssh/terminal"
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

	username := ""
	fmt.Println("输入12306账户: ")
	fmt.Scanf("%s\n", &username)
	fmt.Println("输入密码: ")
	password, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		log.Fatal(err)
		return
	}
	println("******")

	trainDate := ""
	fmt.Println("输入乘车日期(2006-01-02):")
	fmt.Scanf("%s\n", &trainDate)
	t, _ := time.Parse("2006-01-02", trainDate)

	fmt.Println("请输入身份证号码:")
	idcard := ""
	fmt.Scanf("%s\n", &idcard)

	fmt.Println("出发地代码:")
	fromStation := ""
	fmt.Scanf("%s\n", &fromStation)

	fmt.Println("目的地代码:")
	toStation := ""
	fmt.Scanf("%s\n", &toStation)

	fmt.Println(`请输入座位编码:特等座(P),一等座(M),二等座(O),软卧(4),硬卧(3),软座(2),硬座(1),无座(0)`)
	var seattype byte
	fmt.Scanf("%c\n", &seattype)

	amount := &contract.AccountInfo{
		UserName: username,
		Password: string(password),
		IDCards:  []string{idcard},
	}

	query := &contract.TicketQuery{
		TrainDate:    t,
		FromStation:  fromStation,
		ToStation:    toStation,
		PurposeCodes: "ADULT",
		IntervalTime: 3 * time.Second,
		SeatTypes:    []byte{seattype},
	}
	err = client.Start(amount, query)
	if err != nil {
		log.Println("发生错误:", err.Error())
	}

	<-sign

}
