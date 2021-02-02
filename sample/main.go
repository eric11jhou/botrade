package main

import(
	"fmt"
	bt "github.com/eric11jhou/botrade"
)

const (
	API_KEY = ""
	SECRET_KEY = ""
)

type CustomStrategy struct {
	*bt.Advisor
}

func (c *CustomStrategy) SetAdvisor(a *bt.Advisor) {
	c.Advisor = a
}

func (c *CustomStrategy) OnInit() {
	
}

func (c *CustomStrategy) OnDeinit() {

}

var isOpen = false

func (c *CustomStrategy) OnTick() {
	fmt.Printf("Price: %f\n", c.Close("1h", 1))
	if c.Close("1h", 1) >= 40904 && !isOpen {
		r, _ := c.OrderCreate(&bt.OrderCreateRequest{
			Symbol: "BTCUSDT",
			SideType: "SELL",
			OrderType: "MARKET",
			TimeInForce: "GTC",
			Quantity: 1.1,
		})
		fmt.Println(r)
	}
}

func main() {
	bot := bt.NewBot(API_KEY, SECRET_KEY)

	bot.Testing(100000, "BTCUSDT", &CustomStrategy{}, 1610172871000, 1610272871000)
	//for{}
}