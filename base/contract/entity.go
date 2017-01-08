package contract

import "io"
import "time"

//RespBody ...
type RespBody io.ReadCloser

//TicketQuery ...
type TicketQuery struct {
	TrainDate    time.Time
	FromStation  string
	ToStation    string
	PurposeCodes string
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
