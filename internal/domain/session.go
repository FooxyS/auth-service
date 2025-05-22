package domain

type Session struct {
	UserID      string
	IPAddress   string
	RefreshHash string
	PairID      string
	UserAgent   string
}
