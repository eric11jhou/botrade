package botrade

import (
	"fmt"
	"github.com/adshao/go-binance/v2"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetReportCaller(true)
	log.SetFormatter(&log.JSONFormatter{})
}

// Strategy 策略需實作的方法
type Strategy interface {
	SetAdvisor(*Advisor)
	OnInit()
	OnDeinit()
	OnTick()
}

// Bot 交易機器人
type Bot struct {
	advisor *Advisor
}

// NewBot 建立新Bot
func NewBot(apiKey, secretKey string) *Bot {
	return &Bot{
		advisor: &Advisor{
			apiKey: apiKey,
			secretKey: secretKey,
			tick: make(chan struct{}),
			nextTick: make(chan struct{}, 1),
			kline: make(map[string][]*binance.Kline),
			klineTemp: make(map[string][]*binance.Kline),
			openOrders: make([]*binance.Order, 0),
			orders: make([]*binance.Order, 0),
		},
	}
}

// Trading 開始交易
func (b *Bot) Trading(symbol string, s Strategy) {
	b.advisor.trade = true
	s.SetAdvisor(b.advisor)
	b.advisor.loadHistoryData(symbol)
	s.OnInit()
	b.advisor.startTick(symbol)
	for {
		<- b.advisor.tick
		b.advisor.mutex.Lock()
		s.OnTick()
		b.advisor.mutex.Unlock()
	}
}

// Testing 開始回測
func (b *Bot) Testing(balance float64, symbol string, s Strategy, startTime, endTime int64) {
	b.advisor.balance = balance
	b.advisor.trade = false
	s.SetAdvisor(b.advisor)
	b.advisor.loadHistoryDataTesting(symbol, startTime, endTime)
	s.OnInit()
	b.advisor.startTickTesting(symbol, startTime, endTime)
	b.advisor.nextTick <- struct{}{}
	for _ = range b.advisor.tick{
		b.advisor.orderCheckExecTesting()
		s.OnTick()
		b.advisor.nextTick <- struct{}{}
	}
	fmt.Println(b.advisor.equity())
}