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

const (
	//TICKETTYPEADULT 成人票
	TICKETTYPEADULT int = 1
	//TICKETTYPECHILD 儿童票
	TICKETTYPECHILD int = 2
	//TICKETTYPESTUDENT 学生票
	TICKETTYPESTUDENT int = 3
	//TICKETTYPECANJUN 残军票
	TICKETTYPECANJUN int = 4
)

//AccountInfo ...
type AccountInfo struct {
	UserName string
	Password string
	IDCards  []string
}

//TicketQuery ...
type TicketQuery struct {
	TrainDate        time.Time       //购票日期
	FromStation      string          //始发站
	ToStation        string          //到站
	PurposeCodes     string          //乘车人类别 查票页面 ADULT成人票  0X00学生票
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
	TourFlag             string //dc表示单程
	PurposeCodes         string
	QueryFromStationName string //中文
	QueryToStationName   string //中文
	FromStationTelecode  string //编号
	ToStationTelecode    string //编号
	SeatTypes            byte   //座位类型
}

//CheckOutOrderContext ...
type CheckOutOrderContext struct {
	VCodeMod          IVCode
	LoginMod          ILogin
	UserName          string
	Pwd               string
	PassengerIDCardNo []string
	SecretStr         string
	Train             TrainInfo
	SeatType          byte
	TicketType        int
}

//TrainInfo ...
type TrainInfo struct {
	StationTrainCode     string
	TrainDate            string
	BackTrainDate        string
	TourFlag             string //dc表示单程
	PurposeCodes         string
	QueryFromStationName string
	QueryToStationName   string
}
