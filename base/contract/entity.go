package contract

import "io"
import "time"

//RespBody ...
type RespBody io.ReadCloser

const (
	//TZ 特等座
	TZ = 'P'
	//ZY 一等座
	ZY = 'M'
	//ZE 二等座
	ZE = 'O'
	//GR   //高级软卧
	//RW 软卧
	RW = '4'
	//YW 硬卧
	YW = '3'
	//SRRB 动卧
	SRRB = 'F'
	//YYRW 高级动卧
	YYRW = 'A'
	//RZ 软座
	RZ = '2'
	//YZ 硬座
	YZ = '1'
	//WZ 无座
	WZ = '0'
)

//TicketQuery ...
type TicketQuery struct {
	TrainDate        time.Time       //购票日期
	FromStation      string          //始发站
	ToStation        string          //到站
	PurposeCodes     string          //类型
	StationTrainCode map[string]byte //车次
	SeatTypes        []byte          //座位类型
	IntervalTime     time.Duration   //查询的间隔时间
}

//TicketResult ...
type TicketResult struct {
	SecretStr            string
	StationTrainCode     string
	TrainDate            string
	BackTrainDate        string
	TourFlag             string
	PurposeCodes         string
	QueryFromStationName string
	QueryToStationName   string
	SeatTypes            byte //座位类型
}

//CheckOutOrderContext ...
type CheckOutOrderContext struct {
	Mod               CheckOrderMod
	UserName          string
	Pwd               string
	PassengerIDCardNo []string
	SecretStr         string
	Train             TrainInfo
	SeatType          string
	TicketType        string
}

//CheckOrderMod ...
type CheckOrderMod struct {
	VCode IVCode
	Login ILogin
}

//TrainInfo ...
type TrainInfo struct {
	StationTrainCode     string
	TrainDate            string
	BackTrainDate        string
	TourFlag             string
	PurposeCodes         string
	QueryFromStationName string
	QueryToStationName   string
}
