package botrade

// Advisor 可取得各種訊息與交易功能
type Advisor struct {
	trade bool // true實倉交易, false策略測試
	apiKey string
	secretKey string
	tick chan float64 // 新報價觸發通道
}