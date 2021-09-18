package common

type Post struct {
	Id         string
	Username   string
	UserPicUrl string
	Body       string
}

type Ohlc struct {
	Id    string
	Open  string
	High  string
	Low   string
	Close string
}
