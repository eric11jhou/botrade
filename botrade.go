package botrade

// Strategy 策略需實作的方法
type Strategy interface {
	SetAdvisor(*Advisor)
	OnInit()
	OnDeinit()
	OnTick()
}

// Bot 交易機器人
type Bot struct {
	apiKey string
	secretKey string
}

// NewBot 建立新Bot
func NewBot(apiKey, secretKey string) *Bot {
	return &Bot{
		apiKey: apiKey,
		secretKey: secretKey,
	}
}

// Trading 開始交易
func (b *Bot) Trading(s Strategy) {
	s.SetAdvisor(&Advisor{})
	s.OnInit()
	s.OnTick()
	s.OnDeinit()
}

// Testing 開始回測
func (b *Bot) Testing(s Strategy) {

}