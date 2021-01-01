package botrade

import (
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
		},
	}
}

// Trading 開始交易
func (b *Bot) Trading(symbol string, s Strategy) {
	b.advisor.trade = true
	s.SetAdvisor(b.advisor)
	s.OnInit()
	b.advisor.startTick(symbol)
	// for {
	// 	<- b.advisor.tick
	// 	s.OnTick()
	// }
	s.OnDeinit()
}

// Testing 開始回測
func (b *Bot) Testing(symbol string, s Strategy) {
	b.advisor.trade = false
	s.SetAdvisor(b.advisor)
	s.OnInit()
	go b.advisor.startTick(symbol)
	for {
		<- b.advisor.tick
		s.OnTick()
	}
}