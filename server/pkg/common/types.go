package common

import "time"

type Post struct {
	Id         string
	Username   string
	UserPicUrl string
	Body       string
	Timestamp  time.Time
}

type Price struct {
	Asset      string
	TradePrice float32
	Timestamp  time.Time
}

type OHLC struct {
	Open      float32
	High      float32
	Low       float32
	Close     float32
	StartTime time.Time
}

type Wallet struct {
	Id       string
	Username string
	Asset    string
	Balance  float32
}
