package ticketmod

//RespHead ...
type RespHead struct {
	ValidateMessagesShowID string      `json:"validateMessagesShowId"`
	Status                 bool        `json:"status"`
	Httpstatus             int         `json:"httpstatus"`
	Messages               []string    `json:"messages"`
	ValidateMessages       interface{} `json:"validateMessages"`
}

//ticketResult ...
type ticketResult struct {
	RespHead
	Data []queryData `json:"data"`
}

type queryData struct {
	Dto            queryLeftNewDTO `json:"queryLeftNewDTO"`
	SecretStr      string          `json:"secretStr"`
	ButtonTextInfo string          `json:"buttonTextInfo"`
}

//queryLeftNewDTO ...
type queryLeftNewDTO struct {
	TrainNo                string `json:"train_no"`
	StationTrainCode       string `json:"station_train_code"`
	StartStationTelecode   string `json:"start_station_telecode"`
	StartStationName       string `json:"start_station_name"`
	EndStationTelecode     string `json:"end_station_telecode"`
	EndStationName         string `json:"end_station_name"`
	FromStationTelecode    string `json:"from_station_telecode"`
	FromStationName        string `json:"from_station_name"`
	ToStationTelecode      string `json:"to_station_telecode"`
	ToStationName          string `json:"to_station_name"`
	StartTime              string `json:"start_time"`
	ArriveTime             string `json:"arrive_time"`
	DayDifference          string `json:"day_difference"`
	TrainClassName         string `json:"train_class_name"`
	Lishi                  string `json:"lishi"`
	CanWebBuy              string `json:"canWebBuy"`
	LishiValue             string `json:"lishiValue"`
	YpInfo                 string `json:"yp_info"`
	ControlTrainDay        string `json:"control_train_day"`
	StartTrainDate         string `json:"start_train_date"`
	SeatFeature            string `json:"seat_feature"`
	YpEx                   string `json:"yp_ex"`
	TrainSeatFeature       string `json:"train_seat_feature"`
	SeatTypes              string `json:"seat_types"`
	LocationCode           string `json:"location_code"`
	FromStationNo          string `json:"from_station_no"`
	ToStationNo            string `json:"to_station_no"`
	ControlDay             int    `json:"control_day"`
	SaleTime               string `json:"sale_time"`
	IsSupportCard          string `json:"is_support_card"`
	ControlledTrainFlag    string `json:"controlled_train_flag"`
	ControlledTrainMessage string `json:"controlled_train_message"`
	TrainTypeCode          string `json:"train_type_code"`
	StartProvinceCode      string `json:"start_province_code"`
	StartCityCode          string `json:"start_city_code"`
	EndProvinceCode        string `json:"end_province_code"`
	EndCityCode            string `json:"end_city_code"`
	YzNum                  string `json:"yz_num"` //硬座
	RzNum                  string `json:"rz_num"`
	YwNum                  string `json:"yw_num"` //硬卧
	RwNum                  string `json:"rw_num"`
	GrNum                  string `json:"gr_num"`
	ZyNum                  string `json:"zy_num"` //一等座
	ZeNum                  string `json:"ze_num"` //二等座
	TzNum                  string `json:"tz_num"` //特等座位
	GgNum                  string `json:"gg_num"`
	YbNum                  string `json:"yb_num"`
	WzNum                  string `json:"wz_num"` //无座
	QtNum                  string `json:"qt_num"`
	SwzNum                 string `json:"swz_num"`
}

//submitOrderResp ...
type submitOrderResp struct {
	RespHead
	Data string `json:"data"`
}

//getPassengerDTOResp ..
type getPassengerDTOResp struct {
	RespHead
	Data passengerData `json:"data"`
}

//passengerData
type passengerData struct {
	IsExist          bool              `json:"isExist"`
	ExMsg            string            `json:"exMsg"`
	TwoIsOpenClick   []string          `json:"two_isOpenClick"`
	OtherIsOpenClick []string          `json:"other_isOpenClick"`
	NormalPassengers []normalPassenger `json:"normal_passengers"`
}

//normalPassengers
type normalPassenger struct {
	Code                string `json:"code"`
	PassengerName       string `json:"passenger_name"`
	SexCode             string `json:"sex_code"`
	SexName             string `json:"sex_name"`
	BornDate            string `json:"born_date"`
	CountryCode         string `json:"country_code"`
	PassengerIDTypeCode string `json:"passenger_id_type_code"`
	PassengerIDTypeName string `json:"passenger_id_type_name"`
	PassengerIDNo       string `json:"passenger_id_no"`
	PassengerType       string `json:"passenger_type"`
	PassengerFlag       string `json:"passenger_flag"`
	PassengerTypeName   string `json:"passenger_type_name"`
	MobileNo            string `json:"mobile_no"`
	PhoneNo             string `json:"phone_no"`
	Email               string `json:"email"`
	Address             string `json:"address"`
	PostalCode          string `json:"postalcode"`
	FirstLetter         string `json:"first_letter"`
	RecordCount         string `json:"recordCount"`
	TotalTimes          string `json:"total_times"`
	IndexID             string `json:"index_id"`
}

//checkOrderResp ...
type checkOrderResp struct {
	RespHead
	Data checkOrderData `json:"data"`
}

type checkOrderData struct {
	IfShowPassCode     string `json:"ifShowPassCode"`
	CanChooseBeds      string `json:"canChooseBeds"`
	CanChooseSeats     string `json:"canChooseSeats"`
	ChooseSeats        string `json:"choose_Seats"`
	IsCanChooseMid     string `json:"isCanChooseMid"`
	IfShowPassCodeTime string `json:"ifShowPassCodeTime"`
	SubmitStatus       bool   `json:"submitStatus"`
	SmokeStr           string `json:"smokeStr"`
	IsRelogin          bool   `json:"isRelogin"`
	IsNoActive         bool   `json:"isNoActive"`
	CheckSeatNum       bool   `json:"checkSeatNum"`
	ErrMsg             string `json:"errMsg"`
}

type getQueueCountResp struct {
	RespHead
	Data getQueueCountData `json:"data"`
}
type getQueueCountData struct {
	Count  string `json:"count"`
	Ticket string `json:"ticket"`
	Op2    string `json:"op_2"`
	CountT string `json:"countT"`
	Op1    string `json:"op_1"`
}

//confirmSingleForQueueResp ...
type confirmSingleForQueueResp struct {
	RespHead
	Data confirmSingleForQueueData `json:"data"`
}

type confirmSingleForQueueData struct {
	SubmitStatus bool   `json:"submitStatus"`
	ErrMsg       string `json:"errMsg"`
}

type queryOrderWaitTimeResp struct {
	RespHead
	Data queryOrderWaitTimeData `json:"data"`
}

type queryOrderWaitTimeData struct {
	QueryOrderWaitTimeStatus bool   `json:"queryOrderWaitTimeStatus"`
	Count                    int    `json:"count"`
	WaitTime                 int    `json:"waitTime"`
	RequestID                uint64 `json:"requestId"`
	WaitCount                int    `json:"waitCount"`
	TourFlag                 string `json:"tourFlag"`
	Errorcode                string `json:"errorcode"`
	Msg                      string `json:"msg"`
	OrderID                  string `json:"orderId"`
}

//confirmPassengerInitDcResp ...
type confirmPassengerInitDcResp struct {
	GlobalRepeatSubmitToken    string
	TicketInfoForPassengerForm *ticketInfoForPassengerForm
	QueryLeftTicketRequest     *queryLeftTicketRequestDTO
}

//ticketInfoForPassengerForm
type ticketInfoForPassengerForm struct {
	CardTypes        []cardType      `json:"cardTypes"`
	IsAsync          string          `json:"isAsync"`
	KeyCheckIsChange string          `json:"key_check_isChange"`
	LeftDetails      []string        `json:"leftDetails"`
	LeftTicketStr    string          `json:"leftTicketStr"`
	MaxTicketNum     string          `json:"maxTicketNum"`
	OrderRequestDTO  orderRequestDTO `json:"orderRequestDTO"`
	PurposeCodes     string          `json:"purpose_codes"`
	TourFlag         string          `json:"tour_flag"`
	TrainLocation    string          `json:"train_location"`
}

type cardType struct {
	EndStationName   string `json:"end_station_name"`
	EndTime          string `json:"end_time"`
	ID               string `json:"id"`
	StartStationName string `json:"start_station_name"`
	StartTime        string `json:"start_time"`
	Value            string `json:"value"`
}

type orderRequestDTO struct {
	AdultNum            int                `json:"adult_num"`
	ApplyOrderNo        string             `json:"apply_order_no"`
	BedLevelOrderNum    string             `json:"bed_level_order_num"`
	BureauCode          string             `json:"bureau_code"`
	CancelFlag          string             `json:"cancel_flag"`
	CardNum             string             `json:"card_num"`
	Channel             string             `json:"channel"`
	ChildNum            int                `json:"child_num"`
	ChooseSeat          string             `json:"choose_seat"`
	DisabilityNum       int                `json:"disability_num"`
	EndTime             oderReqeustDTOTime `json:"end_time"`
	FromStationName     string             `json:"from_station_name"`
	FromStationTelecode string             `json:"from_station_telecode"`
	GetTicketPass       string             `json:"get_ticket_pass"`
	IDMode              string             `json:"id_mode"`
	IsShowPassCode      string             `json:"isShowPassCode"`
	LeftTicketGenTime   string             `json:"leftTicketGenTime"`
	OrderDate           string             `json:"order_date"`
	RealleftTicket      string             `json:"realleftTicket"`
	ReqIPAddress        string             `json:"reqIpAddress"`
	ReqTimeLeftStr      string             `json:"reqTimeLeftStr"`
	ReserveFlag         string             `json:"reserve_flag"`
	SeatDetailTypeCode  string             `json:"seat_detail_type_code"`
	SeatTypeCode        string             `json:"seat_type_code"`
	SequenceNo          string             `json:"sequence_no"`
	StartTime           oderReqeustDTOTime `json:"start_time"`
	StartTimeStr        string             `json:"start_time_str"`
	StationTrainCode    string             `json:"station_train_code"`
	StudentNum          int                `json:"student_num"`
	TicketNum           int                `json:"ticket_num"`
	TicketTypeOrderNum  string             `json:"ticket_type_order_num"`
	ToStationName       string             `json:"to_station_name"`
	ToStationTelecode   string             `json:"to_station_telecode"`
	TourFlag            string             `json:"tour_flag"`
	TrainCodeText       string             `json:"trainCodeText"`
	TrainDate           oderReqeustDTOTime `json:"train_date"`
	TrainDateStr        string             `json:"train_date_str"`
	TrainLocation       string             `json:"train_location"`
	TrainNo             string             `json:"train_no"`
	TrainOrder          string             `json:"train_order"`
	VarStr              string             `json:"varStr"`
}

//oderReqeustDTOTime
type oderReqeustDTOTime struct {
	Date           int   `json:"date"`
	Day            int   `json:"day"`
	Hours          int   `json:"hours"`
	Minutes        int   `json:"minutes"`
	Month          int   `json:"month"`
	Seconds        int   `json:"seconds"`
	Time           int64 `json:"time"`
	TimezoneOffset int   `json:"timezoneOffset"`
	Year           int   `json:"year"`
}

//queryLeftTicketRequest ...
type queryLeftTicketRequestDTO struct {
	YpInfoDetail string `json:"ypInfoDetail"`
}
