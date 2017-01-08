package piaomod

//submitOrderReqInfo ...
type submitOrderReqInfo struct {
	SecretStr            string
	StationTrainCode     string
	TrainDate            string
	BackTrainDate        string
	TourFlag             string
	PurposeCodes         string
	QueryFromStationName string
	QueryToStationName   string
}

//getPassengerDTOsReqInfo ..._json_att=&REPEAT_SUBMIT_TOKEN=02b853c516d144427f39c393fd0fe159
type getPassengerDTOsReqInfo struct {
	JSONAtt           string
	RepeatSubmitToken string
}

//checkOrderReqInfo
type checkOrderReqInfo struct {
	CancelFlag         string
	BedLevelOrderNum   string
	PassengerTicketStr string
	OldPassengerStr    string
	TourFlag           string
	RandCode           string
	JSONAtt            string
	RepeatSubmitToken  string
}

//getQueueCountReq
type getQueueCountReq struct {
	TrainDate           string
	TrainNo             string
	StationTrainCode    string
	SeatType            string
	FromStationTelecode string
	ToStationTelecode   string
	LeftTicket          string
	PurposeCodes        string
	TrainLocation       string
	JSONAtt             string
	RepeatSubmitToken   string
}

//confirmSingleForQueue
type confirmSingleForQueueReq struct {
	PassengerTicketStr string `name:"passengerTicketStr"`
	OldPassengerStr    string `name:"oldPassengerStr"`
	RandCode           string `name:"randCode"`
	PurposeCodes       string `name:"purpose_codes"`
	KeyCheckIsChange   string `name:"key_check_isChange"`
	LeftTicketStr      string `name:"leftTicketStr"`
	TrainLocation      string `name:"train_location"`
	ChooseSeats        string `name:"choose_seats"`
	SeatDetailType     string `name:"seatDetailType"`
	RoomType           string `name:"roomType"`
	DwAll              string `name:"dwAll"`
	JSONAtt            string `name:"_json_att"`
	RepeatSubmitToken  string `name:"REPEAT_SUBMIT_TOKEN"`
}

//queryOrderWaitTimeReq
type queryOrderWaitTimeReq struct {
	Random            string
	TourFlag          string
	JSONAtt           string
	RepeatSubmitToken string
}
