package common

import "time"

type Post struct {
	Id         string    `json:"id"`
	Username   string    `json:"username"`
	UserPicUrl string    `json:"userpicurl"`
	Body       string    `json:"body"`
	Timestamp  time.Time `json:"timestamp"`
}

type Price struct {
	Asset      string
	TradePrice float32
	Timestamp  time.Time
}

type Ohlc struct {
	Id        string
	Open      string
	High      string
	Low       string
	Close     string
	StartTime time.Time
	EndTime   time.Time
}
