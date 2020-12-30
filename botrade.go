package botrade

type Strategy interface {
	SetAdvisor(*Advisor)
	OnInit()
	OnDeinit()
	OnTick()
}

type Bot struct {
	apiKey string
	secretKey string
}

func NewBot(apiKey, secretKey string) *Bot {
	return &Bot{
		apiKey: apiKey,
		secretKey: secretKey,
	}
}

func (b *Bot) Trading(s Strategy) {
	s.SetAdvisor(&Advisor{})
	s.OnInit()
	s.OnTick()
	s.OnDeinit()
}
