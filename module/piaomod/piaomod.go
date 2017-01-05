package piaomod

import "traintickets/base/contract"
import "fmt"
import "traintickets/base/piaohttputil"
import "encoding/json"
import "errors"

//PIAO ...
type PIAO struct {
}

//ticketResult ...
type ticketResult struct {
	ValidateMessagesShowID string      `json:"validateMessagesShowId"`
	Status                 bool        `json:"status"`
	Httpstatus             int         `json:"httpstatus"`
	Data                   []queryData `json:"data"`
	Messages               []string    `json:"messages"`
	ValidateMessages       interface{} `json:"validateMessages"`
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
	RwNum                  string `json:"rw_num"`
	GrNum                  string `json:"gr_num"`
	ZyNum                  string `json:"zy_num"`
	ZeNum                  string `json:"ze_num"`
	TzNum                  string `json:"tz_num"`
	GgNum                  string `json:"gg_num"`
	YbNum                  string `json:"yb_num"`
	WzNum                  string `json:"wz_num"`
	QtNum                  string `json:"qt_num"`
}

//QueryATicket ...
func (piao *PIAO) QueryATicket(query contract.TicketQuery) (*string, error) {
	r, err := ticketLog(query)
	if err != nil {
		return nil, err
	}
	if !r {
		return nil, errors.New("ticketLog false")
	}
	_, err = queryTicket(query)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

//ResolveQueryAResult ...
func resolveQueryAResult(data *string) (*ticketResult, error) {

	return nil, nil
}

func ticketLog(query contract.TicketQuery) (bool, error) {

	//https://kyfw.12306.cn/otn/leftTicket/log?leftTicketDTO.train_date=2017-01-26&leftTicketDTO.from_station=SHH&leftTicketDTO.to_station=BJP&purpose_codes=ADULT

	formatStr := "https://kyfw.12306.cn/otn/leftTicket/log?leftTicketDTO.train_date=%s&leftTicketDTO.from_station=%s&leftTicketDTO.to_station=%s&purpose_codes=%s"
	date := fmt.Sprintf("%d-%02d-%02d", query.TrainDate.Year(), query.TrainDate.Month(), query.TrainDate.Day())
	url := fmt.Sprintf(formatStr, date, query.FromStation, query.ToStation, query.PurposeCodes)
	resp, err := piaohttputil.Get(url)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	buf, err := piaohttputil.ReadRespBody(resp.Body)
	if err != nil {
		return false, err
	}
	fmt.Println("log:", buf.String())
	result := &ticketResult{}
	err = json.Unmarshal(buf.Bytes(), result)
	if err != nil {
		return false, err
	}

	return result.Status, nil
}

func queryTicket(query contract.TicketQuery) (*ticketResult, error) {
	formatStr := "https://kyfw.12306.cn/otn/leftTicket/queryA?leftTicketDTO.train_date=%s&leftTicketDTO.from_station=%s&leftTicketDTO.to_station=%s&purpose_codes=%s"
	date := fmt.Sprintf("%d-%02d-%02d", query.TrainDate.Year(), query.TrainDate.Month(), query.TrainDate.Day())
	url := fmt.Sprintf(formatStr, date, query.FromStation, query.ToStation, query.PurposeCodes)

	resp, err := piaohttputil.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	buf, err := piaohttputil.ReadRespBody(resp.Body)
	if err != nil {
		return nil, err
	}

	fmt.Println("data:", buf.String())
	result := &ticketResult{}
	err = json.Unmarshal(buf.Bytes(), result)
	if err != nil {
		return nil, err
	}
	fmt.Println("reslut:=", result)
	return result, nil
}
