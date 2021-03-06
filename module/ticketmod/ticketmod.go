package ticketmod

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
	"traintickets/base/appconfig"
	"traintickets/base/contract"
	"traintickets/base/piaohttputil"
)

const (
	//STOP 停止状态
	STOP = -1
	//RUN 运行状态
	RUN = 0
	//WAIT 等待状态
	WAIT = 1
)

//PIAO ...
type PIAO struct {
	state int
}

var (
	clientid         int
	mux              sync.Mutex
	chanStopSlice    = make([](chan<- bool), 1)
	chanWaitSlice    = make([](chan<- bool), 1)
	chanRestartSlice = make([](chan<- bool), 1)

	pat1 = "globalRepeatSubmitToken\\s+=\\s+'([a-z0-9]+)'"
	reg1 = regexp.MustCompile(pat1)
	pat2 = "var\\s*ticketInfoForPassengerForm\\s*=\\s*(\\{.*\\})\\s*;"
	reg2 = regexp.MustCompile(pat2)
	pat3 = "'"
	reg3 = regexp.MustCompile(pat3)

	appconf        = appconfig.GetAppConfig()
	cLeftTicketURL = "leftTicket/queryA"
)

//ReStart 重新启动所有的查票线程
func (piao *PIAO) ReStart() error {
	if piao.state == STOP {
		return errors.New("STATE STOP")
	}
	if piao.state == WAIT {
		for _, c := range chanRestartSlice {
			c <- true
		}
	}
	return nil
}

//Wait 阻塞所有的查票线程
func (piao *PIAO) Wait() error {
	if piao.state == STOP {
		return errors.New("STATE STOP")
	}
	if piao.state != WAIT {
		for _, c := range chanWaitSlice {
			c <- true
		}
	}
	return nil
}

//Stop 停止所有所有的查票线程
func (piao *PIAO) Stop() {
	if piao.state != STOP {
		piao.ReStart()
		for _, c := range chanStopSlice {
			c <- true
		}
	}
}

//QueryTicket ...
func (piao *PIAO) QueryTicket(query *contract.TicketQuery) <-chan *contract.TicketResult {
	clientid = safeAddValue(clientid, 1)
	res := make(chan *contract.TicketResult, 1)

	for i := 0; i < 1; i++ {
		stopSign := make(chan bool, 1)
		waitSign := make(chan bool, 1)
		restartSign := make(chan bool, 1)
		chanStopSlice = append(chanStopSlice, stopSign)
		chanWaitSlice = append(chanWaitSlice, waitSign)
		chanRestartSlice = append(chanRestartSlice, restartSign)
		go func(s <-chan bool, w <-chan bool, r <-chan bool) {
			ticker := time.NewTicker(query.IntervalTime)
			defer ticker.Stop()
		stop:
			for {
				select {
				case <-s:
					break stop
				case <-w:
					<-r
				case <-ticker.C:
					log.Println("开始查询剩余车票")
					f, ticket, err := queryATicket(clientid, query)
					log.Println("结束查询剩余车票")
					if err != nil {
						fmt.Println(err.Error())
						if err.Error() == "IS_TIME_NOT_BUY" {
							break stop
						}
					}
					//抢到票了
					if f {
						log.Println("存在指定的票,开始买票")
						log.Println(ticket)
						res <- ticket
						//break stop
						piao.Wait()
					}

				}
			}
			log.Println("结束")

		}(stopSign, waitSign, restartSign)
	}
	return res
}

func safeAddValue(v1 int, v2 int) int {
	mux.Lock()
	defer mux.Unlock()
	return v1 + v2
}

//QueryATicket ...
func queryATicket(clientID int, query *contract.TicketQuery) (bool, *contract.TicketResult, error) {

	r, err := ticketLog(clientID, query)
	if err != nil {
		return false, nil, err
	}
	if !r {
		return false, nil, errors.New("ticketLog false")
	}
	res, err := queryTicket(clientID, query, cLeftTicketURL)
	if err != nil {
		return false, nil, err
	}

	result, err := resolveQueryAResult(res)
	if err != nil {
		return false, nil, err
	}
	if !result.Status {
		if result.CNAME == "CLeftTicketUrl" {
			cLeftTicketURL = result.CURL
		}
		return false, nil, fmt.Errorf("查票失败:%v", result)
	}
	//判断是否存在指定的票
	for _, p := range result.Data {
		if p.Dto.CanWebBuy == "IS_TIME_NOT_BUY" {
			return false, nil, errors.New("IS_TIME_NOT_BUY")
		}
		//如果指定了车次则只判断指定的车次
		if len(query.StationTrainCode) > 0 {
			_, ok := query.StationTrainCode[strings.ToLower(p.Dto.StationTrainCode)]
			if !ok {
				continue
			}
		}
		tc := &contract.TicketResult{
			StationTrainCode:     p.Dto.StationTrainCode,
			TrainDate:            query.TrainDate.Format("2006-01-02"),
			BackTrainDate:        query.TrainDate.Format("2006-01-02"),
			TourFlag:             "dc",
			PurposeCodes:         query.PurposeCodes,
			QueryFromStationName: p.Dto.FromStationName,
			QueryToStationName:   p.Dto.ToStationName,
			FromStationTelecode:  p.Dto.FromStationTelecode,
			ToStationTelecode:    p.Dto.ToStationTelecode,
		}
		//判断指定的座位类型是否有票
		for _, seat := range query.SeatTypes {
			switch seat {
			case contract.TZ: //特等
				tzStr := p.Dto.TzNum
				tzNum, _ := strconv.Atoi(p.Dto.TzNum)
				if tzStr == "有" || tzNum > 0 {
					tc.SecretStr = p.SecretStr
					tc.StationTrainCode = p.Dto.StationTrainCode
					return true, tc, nil
				}
			case contract.ZY: //一等
				zyStr := p.Dto.ZyNum
				zyNum, _ := strconv.Atoi(p.Dto.ZyNum)
				if zyStr == "有" || zyNum > 0 {
					tc.SecretStr = p.SecretStr
					tc.StationTrainCode = p.Dto.StationTrainCode
					return true, tc, nil
				}
			case contract.ZE: //二等
				zeStr := p.Dto.ZeNum
				zeNum, _ := strconv.Atoi(p.Dto.ZeNum)
				if zeStr == "有" || zeNum > 0 {
					tc.SecretStr = p.SecretStr
					tc.StationTrainCode = p.Dto.StationTrainCode
					return true, tc, nil
				}
			case contract.RW: //软卧
				rwStr := p.Dto.RwNum
				rwNum, _ := strconv.Atoi(p.Dto.RwNum)
				if rwStr == "有" || rwNum > 0 {
					tc.SecretStr = p.SecretStr
					tc.StationTrainCode = p.Dto.StationTrainCode
					return true, tc, nil
				}

			case contract.YW: //硬卧
				log.Println(p.Dto.YwNum)
				ywStr := p.Dto.YwNum
				ywNum, _ := strconv.Atoi(p.Dto.YwNum)
				log.Println(ywStr)
				log.Println(ywNum)
				if ywStr == "有" || ywNum > 0 {
					tc.SecretStr = p.SecretStr
					tc.StationTrainCode = p.Dto.StationTrainCode
					tc.SeatTypes = contract.YW
					return true, tc, nil
				}
			case contract.SRRB: //动卧
				// srrbStr := p.Dto.num
				// ywNum, _ := strconv.Atoi(p.Dto.YwNum)
				// if ywStr == "有" || ywNum > 0 {
				// 	//有票
				// }
			case contract.YYRW: //高级动卧
			case contract.RZ: //软座
				rzStr := p.Dto.RzNum
				rzNum, _ := strconv.Atoi(p.Dto.RzNum)
				if rzStr == "有" || rzNum > 0 {
					tc.SecretStr = p.SecretStr
					tc.StationTrainCode = p.Dto.StationTrainCode
					return true, tc, nil
				}
			case contract.YZ: //硬座
				ywStr := p.Dto.YwNum
				ywNum, _ := strconv.Atoi(p.Dto.YwNum)
				if ywStr == "有" || ywNum > 0 {
					tc.SecretStr = p.SecretStr
					tc.StationTrainCode = p.Dto.StationTrainCode
					tc.SeatTypes = contract.YZ
					return true, tc, nil
				}
			case contract.WZ: //无座
				wzStr := p.Dto.WzNum
				wzNum, _ := strconv.Atoi(p.Dto.WzNum)
				if wzStr == "有" || wzNum > 0 {
					tc.SecretStr = p.SecretStr
					tc.StationTrainCode = p.Dto.StationTrainCode
					return true, tc, nil
				}
			default:
				return false, nil, errors.New("无效的座位类型")
			}
		}
	}
	return false, nil, nil
}

// //TicketSResult ...
// func (piao *PIAO) TicketSResult() <-chan (*contract.TicketResult) {
// 	return chanTRes
// }

//ResolveQueryAResult ...
func resolveQueryAResult(data *bytes.Buffer) (*ticketResult, error) {

	//fmt.Println("data:", buf.String())
	result := &ticketResult{}
	err := json.Unmarshal(data.Bytes(), result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func ticketLog(clientID int, query *contract.TicketQuery) (bool, error) {

	//https://kyfw.12306.cn/otn/leftTicket/log?leftTicketDTO.train_date=2017-01-26&leftTicketDTO.from_station=SHH&leftTicketDTO.to_station=BJP&purpose_codes=ADULT
	logurlPrefix, _ := appconfig.Combine(appconf.MainURL, appconf.Ctx, "leftTicket/log")
	formatStr := logurlPrefix + "?leftTicketDTO.train_date=%s&leftTicketDTO.from_station=%s&leftTicketDTO.to_station=%s&purpose_codes=%s"
	date := fmt.Sprintf("%d-%02d-%02d", query.TrainDate.Year(), query.TrainDate.Month(), query.TrainDate.Day())
	urlStr := fmt.Sprintf(formatStr, date, query.FromStation, query.ToStation, query.PurposeCodes)
	referer, _ := appconfig.Combine(appconf.MainURL, appconf.Ctx, "leftTicket/init")
	resp, err := piaohttputil.GetV(clientID, urlStr, referer, true)
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

func queryTicket(clientID int, query *contract.TicketQuery, leftTicketURL string) (*bytes.Buffer, error) {
	urlPrefix, _ := appconfig.Combine(appconf.MainURL, appconf.Ctx, leftTicketURL)
	formatStr := "?leftTicketDTO.train_date=%s&leftTicketDTO.from_station=%s&leftTicketDTO.to_station=%s&purpose_codes=%s"
	date := fmt.Sprintf("%d-%02d-%02d", query.TrainDate.Year(), query.TrainDate.Month(), query.TrainDate.Day())
	urlStr := urlPrefix + fmt.Sprintf(formatStr, date, query.FromStation, query.ToStation, query.PurposeCodes)
	referer, _ := appconfig.Combine(appconf.MainURL, appconf.Ctx, "leftTicket/init")
	resp, err := piaohttputil.GetV(clientID, urlStr, referer, true)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	buf, err := piaohttputil.ReadRespBody(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("resp.StatusCode = %d , %s", resp.StatusCode, buf.String())
	}
	//fmt.Println("reslut:=", buf)
	return buf, nil
}

//CheckOutOrder ...
func (piao *PIAO) CheckOutOrder(clientID int, ckContext *contract.CheckOutOrderContext) (bool, error) {
	lgm := ckContext.LoginMod
	f, err := lgm.CheckUser(clientID)
	if err != nil {
		return false, err
	}
	if !f {
		_, err = lgm.Login(clientID, ckContext.UserName, ckContext.Pwd, ckContext.VCodeMod)
		if err != nil {
			return false, nil
		}
	}

	submitOrder := &submitOrderReqInfo{
		SecretStr:            ckContext.SecretStr,
		StationTrainCode:     ckContext.Train.StationTrainCode,
		TrainDate:            ckContext.Train.TrainDate,
		BackTrainDate:        ckContext.Train.BackTrainDate,
		TourFlag:             ckContext.Train.TourFlag,
		PurposeCodes:         ckContext.Train.PurposeCodes,
		QueryFromStationName: ckContext.Train.QueryFromStationName,
		QueryToStationName:   ckContext.Train.QueryToStationName,
	}
	err = submitOrderRequest(clientID, submitOrder)
	if err != nil {
		return false, err
	}
	//time.Sleep(1 * time.Second)
	initDcInfo, err := confirmPassengerInitDc(clientID)
	if err != nil {
		return false, err
	}
	pgdreq := &getPassengerDTOsReqInfo{
		JSONAtt:           "",
		RepeatSubmitToken: initDcInfo.GlobalRepeatSubmitToken,
	}
	passengers, err := getPassengerDTOs(clientID, pgdreq)
	if err != nil {
		return false, err
	}
	orderPassengers := []normalPassenger{}
	//取得订单人信息
nextPassenger:
	for _, id := range ckContext.PassengerIDCardNo {
		for _, p := range passengers.Data.NormalPassengers {
			if p.PassengerIDNo == id {
				orderPassengers = append(orderPassengers, p)
				continue nextPassenger
			}
		}
		return false, fmt.Errorf("未找到身份证号码为%s的乘客信息下单失败", id)
	}

	//获取验证码，这里还需要判断TODO
	_, err = ckContext.VCodeMod.CaptureVCode(clientID, "passenger", "randp")
	if err != nil {
		return false, nil
	}

	checkOrderReq := &checkOrderReqInfo{
		CancelFlag:         "2",
		BedLevelOrderNum:   "000000000000000000000000000000",
		PassengerTicketStr: getpassengerTickets(orderPassengers, ckContext.SeatType, ckContext.TicketType),
		OldPassengerStr:    getOldPassengers(orderPassengers),
		TourFlag:           initDcInfo.TicketInfoForPassengerForm.TourFlag,
		RandCode:           "", //是否为空？
		RepeatSubmitToken:  initDcInfo.GlobalRepeatSubmitToken,
		JSONAtt:            "",
	}
	checkOderResult, err := checkOrderInfo(clientID, checkOrderReq)
	if err != nil {
		return false, err
	}
	if !checkOderResult.Data.SubmitStatus {
		if checkOderResult.Data.IsRelogin {
			//重新登陆
		} else {
			if checkOderResult.Data.CheckSeatNum {
				return false, fmt.Errorf("很抱歉，无法提交您的订单!原因：%s ", checkOderResult.Data.ErrMsg)
			}
			return false, fmt.Errorf("出票失败!原因：%s ", checkOderResult.Data.ErrMsg)
		}
	}
	if checkOderResult.Data.IfShowPassCode == "Y" {
	}
	requestDtoTrainDate := initDcInfo.TicketInfoForPassengerForm.OrderRequestDTO.TrainDate
	orderRequestDTO := initDcInfo.TicketInfoForPassengerForm.OrderRequestDTO
	leftTicketRequestDTO := initDcInfo.TicketInfoForPassengerForm.QueryLeftTicketRequest
	ts := requestDtoTrainDate.Time / 1000
	t := time.Unix(ts, 000)
	if err != nil {
		return false, err
	}
	ypInfo := leftTicketRequestDTO.YpInfoDetail
	leftTicketStr, err := url.QueryUnescape(initDcInfo.TicketInfoForPassengerForm.LeftTicketStr)
	if err != nil {
		return false, fmt.Errorf("leftTicketStr解码失败:%s,raw:%s", err.Error(), initDcInfo.TicketInfoForPassengerForm.LeftTicketStr)
	}

	getQueueCountR := &getQueueCountReq{
		TrainDate:           t.Format("Mon Jan 02 2006 15:04:05 GMT+0800 (China Standard Time)"),
		TrainNo:             orderRequestDTO.TrainNo,
		StationTrainCode:    orderRequestDTO.StationTrainCode,
		SeatType:            ckContext.SeatType,
		FromStationTelecode: orderRequestDTO.FromStationTelecode,
		ToStationTelecode:   orderRequestDTO.ToStationTelecode,
		LeftTicket:          ypInfo,
		PurposeCodes:        initDcInfo.TicketInfoForPassengerForm.PurposeCodes,
		TrainLocation:       initDcInfo.TicketInfoForPassengerForm.TrainLocation,
		RepeatSubmitToken:   initDcInfo.GlobalRepeatSubmitToken,
		JSONAtt:             "",
	}
	getQueueCount(clientID, getQueueCountR)

	oderforQueue := &confirmSingleForQueueReq{
		PassengerTicketStr: getpassengerTickets(orderPassengers, ckContext.SeatType, ckContext.TicketType),
		OldPassengerStr:    getOldPassengers(orderPassengers),
		RandCode:           "",
		PurposeCodes:       initDcInfo.TicketInfoForPassengerForm.PurposeCodes,
		KeyCheckIsChange:   initDcInfo.TicketInfoForPassengerForm.KeyCheckIsChange,
		LeftTicketStr:      leftTicketStr,
		TrainLocation:      initDcInfo.TicketInfoForPassengerForm.TrainLocation,
		ChooseSeats:        "",
		SeatDetailType:     "000",
		RoomType:           "00",
		DwAll:              "N",
		JSONAtt:            "",
		RepeatSubmitToken:  initDcInfo.GlobalRepeatSubmitToken,
	}
	orderForQueueResult, err := confirmSingleForQueue(clientID, oderforQueue)
	if err != nil {
		return false, err
	}
	if !orderForQueueResult.Data.SubmitStatus {
		return false, fmt.Errorf("下单出票失败!原因：%s ", checkOderResult.Data.ErrMsg)
	}
	//下单成功开始等待 3秒钟调用一次查询接口
	return true, nil
}

//submitOrderRequest ...
func submitOrderRequest(clientID int, reqInfo *submitOrderReqInfo) error {
	vs := make(url.Values, 8)
	secretStr, _ := url.QueryUnescape(reqInfo.SecretStr)
	vs.Add("secretStr", secretStr)
	vs.Add("train_date", reqInfo.TrainDate)
	vs.Add("back_train_date", reqInfo.BackTrainDate)
	vs.Add("tour_flag", reqInfo.TourFlag)
	vs.Add("purpose_codes", reqInfo.PurposeCodes)
	vs.Add("query_from_station_name", reqInfo.QueryFromStationName)
	vs.Add("query_to_station_name", reqInfo.QueryToStationName)
	vsencode := vs.Encode()
	vsencode += "&undefined"
	urlStr, _ := appconfig.Combine(appconf.MainURL, appconf.Ctx, "leftTicket/submitOrderRequest")
	referer, _ := appconfig.Combine(appconf.MainURL, appconf.Ctx, "leftTicket/init")
	resp, err := piaohttputil.PostV(clientID, urlStr, "application/x-www-form-urlencoded; charset=UTF-8", referer, true, strings.NewReader(vsencode))
	defer resp.Body.Close()

	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("submitOrderRequest Error StatusCode:%d", resp.StatusCode)
	}
	buf, err := piaohttputil.ReadRespBody(resp.Body)
	if err != nil {
		return err
	}
	res := &submitOrderResp{}
	err = json.Unmarshal(buf.Bytes(), res)
	if err != nil {
		return err
	}
	if res.Status != true {
		return errors.New(fmt.Sprintln(res.Messages))
	}
	return nil
}

//confirmPassengerInitDc globalRepeatSubmitToken
func confirmPassengerInitDc(clientID int) (*confirmPassengerInitDcResp, error) {
	urlStr, _ := appconfig.Combine(appconf.MainURL, appconf.Ctx, "confirmPassenger/initDc")
	vs := make(url.Values, 1)
	vs.Add("_json_att", "")
	referer, _ := appconfig.Combine(appconf.MainURL, appconf.Ctx, "leftTicket/init")
	resp, err := piaohttputil.PostV(clientID, urlStr, "application/x-www-form-urlencoded", referer, false, strings.NewReader(vs.Encode()))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("confirmPassengerInitDc Error StatusCode:%d", resp.StatusCode)
	}
	buf, err := piaohttputil.ReadRespBody(resp.Body)
	if err != nil {
		return nil, err
	}
	respData := buf.Bytes()
	//fmt.Println(string(respData))
	res := &confirmPassengerInitDcResp{}

	groups1 := reg1.FindSubmatch(respData)
	if len(groups1) < 2 {
		return nil, fmt.Errorf("未匹配到globalRepeatSubmitToken的值！")
	}
	res.GlobalRepeatSubmitToken = string(groups1[1])

	groups2 := reg2.FindSubmatch(respData)
	if len(groups2) < 2 {
		return nil, fmt.Errorf("未匹配到ticketInfoForPassengerForm的值！")
	}

	jStr := reg3.ReplaceAll(groups2[1], []byte{'"'})

	ticketPassengerForm := &ticketInfoForPassengerForm{}
	err = json.Unmarshal(jStr, ticketPassengerForm)
	if err != nil {
		return nil, err
	}
	res.TicketInfoForPassengerForm = ticketPassengerForm

	return res, nil
}

//ownpwud ...
func ownpwud(clientID int) {
	urlStr := "https://kyfw.12306.cn/otn/dynamicJs/ownpwud"
	resp, _ := piaohttputil.Get(clientID, urlStr)
	defer resp.Body.Close()
}

//getPassengerDTOs ..._json_att=&REPEAT_SUBMIT_TOKEN=02b853c516d144427f39c393fd0fe159
func getPassengerDTOs(clientID int, reqInfo *getPassengerDTOsReqInfo) (*getPassengerDTOResp, error) {

	urlStr, _ := appconfig.Combine(appconf.MainURL, appconf.Ctx, "confirmPassenger/getPassengerDTOs") //"https://kyfw.12306.cn/otn/confirmPassenger/getPassengerDTOs"
	vs := make(url.Values, 2)
	vs.Add("_json_att", reqInfo.JSONAtt)
	vs.Add("REPEAT_SUBMIT_TOKEN", reqInfo.RepeatSubmitToken)
	referer, _ := appconfig.Combine(appconf.MainURL, appconf.Ctx, "confirmPassenger/initDc")
	resp, err := piaohttputil.PostV(clientID, urlStr, "application/x-www-form-urlencoded; charset=UTF-8", referer, true, strings.NewReader(vs.Encode()))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("getPassengerDTOs Error StatusCode:%d", resp.StatusCode)
	}
	buf, err := piaohttputil.ReadRespBody(resp.Body)
	if err != nil {
		return nil, err
	}

	passengerDtos := &getPassengerDTOResp{}
	err = json.Unmarshal(buf.Bytes(), passengerDtos)
	if err != nil {
		return nil, err
	}
	return passengerDtos, nil
}

//checkOrderInfo ...
func checkOrderInfo(clientID int, reqInfo *checkOrderReqInfo) (*checkOrderResp, error) {
	urlStr, _ := appconfig.Combine(appconf.MainURL, appconf.Ctx, "confirmPassenger/checkOrderInfo") //"https://kyfw.12306.cn/otn/confirmPassenger/checkOrderInfo"

	vs := make(url.Values, 8)

	vs.Add("cancel_flag", reqInfo.CancelFlag)
	vs.Add("bed_level_order_num", reqInfo.BedLevelOrderNum)
	vs.Add("passengerTicketStr", reqInfo.PassengerTicketStr)
	vs.Add("oldPassengerStr", reqInfo.OldPassengerStr)
	vs.Add("tour_flag", reqInfo.TourFlag)
	vs.Add("randCode", reqInfo.RandCode)
	vs.Add("_json_att", reqInfo.JSONAtt)
	vs.Add("REPEAT_SUBMIT_TOKEN", reqInfo.RepeatSubmitToken)
	referer, _ := appconfig.Combine(appconf.MainURL, appconf.Ctx, "confirmPassenger/initDc")
	resp, err := piaohttputil.PostV(clientID, urlStr, "application/x-www-form-urlencoded; charset=UTF-8", referer, true, strings.NewReader(vs.Encode()))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	buf, err := piaohttputil.ReadRespBody(resp.Body)
	if err != nil {
		return nil, err
	}
	cor := &checkOrderResp{}
	err = json.Unmarshal(buf.Bytes(), cor)
	if err != nil {
		return nil, err
	}
	return cor, nil
}

func getQueueCount(clientID int, reqInfo *getQueueCountReq) (*getQueueCountResp, error) {

	urlStr, _ := appconfig.Combine(appconf.MainURL, appconf.Ctx, "confirmPassenger/getQueueCount") //urlStr := "https://kyfw.12306.cn/otn/confirmPassenger/getQueueCount"
	vs := make(url.Values, 11)
	vs.Add("train_date", reqInfo.TrainDate)
	vs.Add("train_no", reqInfo.TrainNo)
	vs.Add("stationTrainCode", reqInfo.StationTrainCode)
	vs.Add("seatType", fmt.Sprintf("%c", reqInfo.SeatType))
	vs.Add("fromStationTelecode", reqInfo.FromStationTelecode)
	vs.Add("toStationTelecode", reqInfo.ToStationTelecode)
	vs.Add("leftTicket", reqInfo.LeftTicket)
	vs.Add("purpose_codes", reqInfo.PurposeCodes)
	vs.Add("train_location", reqInfo.TrainLocation)
	vs.Add("_json_att", reqInfo.JSONAtt)
	vs.Add("REPEAT_SUBMIT_TOKEN", reqInfo.RepeatSubmitToken)
	referer, _ := appconfig.Combine(appconf.MainURL, appconf.Ctx, "confirmPassenger/initDc")

	resp, err := piaohttputil.PostV(clientID, urlStr, "application/x-www-form-urlencoded; charset=UTF-8", referer, true, strings.NewReader(vs.Encode()))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	buf, err := piaohttputil.ReadRespBody(resp.Body)
	res := &getQueueCountResp{}
	err = json.Unmarshal(buf.Bytes(), res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func confirmSingleForQueue(clientID int, reqInfo *confirmSingleForQueueReq) (*confirmSingleForQueueResp, error) {

	//urlStr := "https://kyfw.12306.cn/otn/confirmPassenger/confirmSingleForQueue"
	urlStr, _ := appconfig.Combine(appconf.MainURL, appconf.Ctx, "confirmPassenger/confirmSingleForQueue")
	vs := make(url.Values, 13)
	vs.Add("passengerTicketStr", reqInfo.PassengerTicketStr)
	vs.Add("oldPassengerStr", reqInfo.OldPassengerStr)
	vs.Add("randCode", reqInfo.RandCode)
	vs.Add("purpose_codes", reqInfo.PurposeCodes)
	vs.Add("key_check_isChange", reqInfo.KeyCheckIsChange)
	vs.Add("leftTicketStr", reqInfo.LeftTicketStr)
	vs.Add("train_location", reqInfo.TrainLocation)
	vs.Add("choose_seats", reqInfo.ChooseSeats)
	vs.Add("seatDetailType", reqInfo.SeatDetailType)
	vs.Add("roomType", reqInfo.RoomType)
	vs.Add("dwAll", reqInfo.DwAll)
	vs.Add("_json_att", reqInfo.JSONAtt)
	vs.Add("REPEAT_SUBMIT_TOKEN", reqInfo.RepeatSubmitToken)
	referer, _ := appconfig.Combine(appconf.MainURL, appconf.Ctx, "confirmPassenger/initDc")

	resp, err := piaohttputil.PostV(clientID, urlStr, "application/x-www-form-urlencoded; charset=UTF-8", referer, true, strings.NewReader(vs.Encode()))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	buf, err := piaohttputil.ReadRespBody(resp.Body)
	if err != nil {
		return nil, err
	}
	res := &confirmSingleForQueueResp{}
	err = json.Unmarshal(buf.Bytes(), res)
	if err != nil {
		log.Println(string(buf.Bytes()))
		return nil, err
	}

	return res, nil
}

//queryOrderWaitTime ...
func queryOrderWaitTime(clientID int, reqInfo queryOrderWaitTimeReq) (*queryOrderWaitTimeData, error) {
	urlStr, _ := appconfig.Combine(appconf.MainURL, appconf.Ctx, "confirmPassenger/queryOrderWaitTime") //"https://kyfw.12306.cn/otn/confirmPassenger/queryOrderWaitTime?"
	urlStr = fmt.Sprintf("%s?random=%s&tourFlag=%s&_json_att=%s&REPEAT_SUBMIT_TOKEN=%s", urlStr, reqInfo.Random, reqInfo.TourFlag, reqInfo.JSONAtt, reqInfo.RepeatSubmitToken)
	referer, _ := appconfig.Combine(appconf.MainURL, appconf.Ctx, "confirmPassenger/initDc")

	resp, err := piaohttputil.GetV(clientID, urlStr, referer, true)
	if err != nil {
		return nil, err
	}
	buf, err := piaohttputil.ReadRespBody(resp.Body)
	if err != nil {
		return nil, err
	}
	res := &queryOrderWaitTimeData{}
	err = json.Unmarshal(buf.Bytes(), res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

//getpassengerTickets ...
func getpassengerTickets(passengers []normalPassenger, seatType byte, ticketType int) string {
	reslt := &bytes.Buffer{}
	for _, p := range passengers {
		b := fmt.Sprintf("%c,%s,%d,%s,%s,%s,%s,%s", seatType, "0", ticketType, p.PassengerName, p.PassengerIDTypeCode, p.PassengerIDNo, p.MobileNo, "N")
		reslt.WriteString(b + "_")
	}
	return string(reslt.Bytes()[:len(reslt.Bytes())-1])
}

func getOldPassengers(passengers []normalPassenger) string {
	reslt := &bytes.Buffer{}
	for _, p := range passengers {
		b := fmt.Sprintf("%s,%s,%s,%s", p.PassengerName, p.PassengerIDTypeCode, p.PassengerIDNo, p.PassengerType)
		reslt.WriteString(b + "_")
	}
	return reslt.String()
}
