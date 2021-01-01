package botrade

// Advisor 可取得各種訊息與交易功能
type Advisor struct {
	trade bool // true實倉交易, false策略測試
	apiKey string
	secretKey string
	tick chan struct{} // 新報價觸發通道
	ask float64
	bid float64
}

func (a *Advisor) Ask() float64 {
	return a.ask
}

func (a *Advisor) Bid() float64 {
	return a.bid
}