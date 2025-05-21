package domain

type Session struct {
	UserID      string
	IPAddress   string
	RefreshHash string
	PairID      string
	UserAgent   string
}

func (s Session) CheckSession(ip, refresh, pair, agent string) bool {
	return ip == s.IPAddress && refresh == s.RefreshHash && pair == s.PairID && agent == s.UserAgent
}
