package botrade

import (
	"sync"
	"strconv"
	"github.com/adshao/go-binance/v2"
	log "github.com/sirupsen/logrus"
)

// Advisor 可取得各種訊息與交易功能
type Advisor struct {
	trade bool // true實倉交易, false策略測試
	balance float64 // **回測用 回測初始資金
	currencyVolume float64 // **回測用 幣持有數量
	lastEquityHigh float64 // **回測用 上次最高淨值
	drawdown float64 // **回測用 回撤率
	apiKey string
	secretKey string
	tick chan struct{} // 新報價觸發通道
	nextTick chan struct{} // **回測用 跑回測時可運算下一個tick的訊號
	mutex sync.Mutex // **實倉用 跑實倉的kline互斥鎖

	ask float64
	bid float64
	time int64
	// key: interval 1m, 3m, 5m, 15m, 30m, 1h, 2h, 4h, 6h, 8h, 12h, 1d, 3d, 1w, 1M
	kline map[string][]*binance.Kline // 目前K線
	klineTemp map[string][]*binance.Kline // **回測用 暫存所有K線

	orderIDCount int64 // **回測用訂單ID
	openOrders []*binance.Order // 掛單
	orders []*binance.Order // 訂單
}

func (a *Advisor) Ask() float64 {
	return a.ask
}

func (a *Advisor) Bid() float64 {
	return a.bid
}

// 算出回測的帳戶淨值
func (a *Advisor) equity() float64 {
	return a.balance + a.currencyVolume * (a.ask + a.bid) / 2
}

// 取得K棒開盤價
// shift: 第幾根K棒
func (a *Advisor) Open(interval string, shift int) float64 {
	if len(a.kline[interval]) <= shift {
		log.Error("out of range")
		return 0
	}
	if open, err := strconv.ParseFloat(a.kline[interval][shift].Open, 64); err != nil {
		log.Error(err)
		return 0
	} else {
		return open
	}
}

// 取得K棒最高點
// shift: 第幾根K棒
func (a *Advisor) High(interval string, shift int) float64 {
	if len(a.kline[interval]) <= shift {
		log.Error("out of range")
		return 0
	}
	if high, err := strconv.ParseFloat(a.kline[interval][shift].High, 64); err != nil {
		log.Error(err)
		return 0
	} else {
		return high
	}
}

// 取得K棒最低點
// shift: 第幾根K棒
func (a *Advisor) Low(interval string, shift int) float64 {
	if len(a.kline[interval]) <= shift {
		log.Error("out of range")
		return 0
	}
	if low, err := strconv.ParseFloat(a.kline[interval][shift].Low, 64); err != nil {
		log.Error(err)
		return 0
	} else {
		return low
	}
}

// 取得K棒收盤價
// shift: 第幾根K棒
func (a *Advisor) Close(interval string, shift int) float64 {
	if len(a.kline[interval]) <= shift {
		log.Error("out of range")
		return 0
	}
	if close, err := strconv.ParseFloat(a.kline[interval][shift].Close, 64); err != nil {
		log.Error(err)
		return 0
	} else {
		return close
	}
}

// 取得K棒成交量
// shift: 第幾根K棒
func (a *Advisor) Volume(interval string, shift int) float64 {
	if len(a.kline[interval]) <= shift {
		log.Error("out of range")
		return 0
	}
	if volume, err := strconv.ParseFloat(a.kline[interval][shift].Volume, 64); err != nil {
		log.Error(err)
		return 0
	} else {
		return volume
	}
}

// 取得K棒開盤時間
// shift: 第幾根K棒
func (a *Advisor) OpenTime(interval string, shift int) int64 {
	if len(a.kline[interval]) <= shift {
		log.Error("out of range")
		return 0
	}
	return a.kline[interval][shift].OpenTime
}

// 取得K棒收盤時間
// shift: 第幾根K棒
func (a *Advisor) CloseTime(interval string, shift int) int64 {
	if len(a.kline[interval]) <= shift {
		log.Error("out of range")
		return 0
	}
	return a.kline[interval][shift].CloseTime
}

// 取得K棒成交額
// shift: 第幾根K棒
func (a *Advisor) QuoteAssetVolume(interval string, shift int) float64 {
	if len(a.kline[interval]) <= shift {
		log.Error("out of range")
		return 0
	}
	if quoteAssetVolume, err := strconv.ParseFloat(a.kline[interval][shift].QuoteAssetVolume, 64); err != nil {
		log.Error(err)
		return 0
	} else {
		return quoteAssetVolume
	}
}

// 取得K棒成交筆數
// shift: 第幾根K棒
func (a *Advisor) TradeNum(interval string, shift int) int64 {
	if len(a.kline[interval]) <= shift {
		log.Error("out of range")
		return 0
	}
	return a.kline[interval][shift].TradeNum
}