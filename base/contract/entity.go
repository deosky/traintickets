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
	SecretStr        string
	StationTrainCode string
}
