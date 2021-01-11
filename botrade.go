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
			nextTick: make(chan struct{}),
			kline: make(map[string][]*binance.Kline),
			klineTemp: make(map[string][]*binance.Kline),
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
		s.OnTick()
	}
}

// Testing 開始回測
func (b *Bot) Testing(symbol string, s Strategy, startTime, endTime int64) {
	b.advisor.trade = false
	s.SetAdvisor(b.advisor)
	b.advisor.loadHistoryDataTesting(symbol, startTime, endTime)
	s.OnInit()
	b.advisor.startTickTesting(symbol, startTime, endTime)
	for {
		b.advisor.nextTick <- struct{}{}
		<- b.advisor.tick
		fmt.Println("tick")
		s.OnTick()
	}
}