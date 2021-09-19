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
	Id        string
	Asset     string
	Price     float32
	Timestamp time.Time
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
