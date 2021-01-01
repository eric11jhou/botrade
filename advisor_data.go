package botrade

import (
	"fmt"
	"strconv"
	"github.com/adshao/go-binance/v2"
	log "github.com/sirupsen/logrus"
)

// interval 1m, 3m, 5m, 15m, 30m, 1h, 2h, 4h, 6h, 8h, 12h, 1d, 3d, 1w, 1M
func (a *Advisor) loadHistoryData(symbol string) {

}

func (a *Advisor) startTick(symbol string) {
	if a.trade {
		a.startTick_(symbol)
	} else {
		a.startTickTesting(symbol)
	}
}

func (a *Advisor) startTick_(symbol string) {
	go func(){
		{
			wsKlineHandler := func(event *binance.WsKlineEvent) {
				fmt.Println(event)
			}
			errHandler := func(err error) {
				log.Error(err)
			}
			_, _, err := binance.WsKlineServe(symbol, "1m", wsKlineHandler, errHandler)
			if err != nil {
				log.Fatal(err)
			}
		}
		{
			wsMarketStatHandler := func(event *binance.WsMarketStatEvent) {
				if ask, err := strconv.ParseFloat(event.AskPrice, 64); err != nil {
					log.Error(err)
				} else {
					a.ask = ask
				}
				if bid, err := strconv.ParseFloat(event.BidPrice, 64); err != nil {
					log.Error(err)
				} else {
					a.bid = bid
				}
				a.tick <- struct{}{}
			}
			errHandler := func(err error) {
				log.Error(err)
			}
			_, _, err := binance.WsMarketStatServe(symbol, wsMarketStatHandler, errHandler)
			if err != nil {
				log.Fatal(err)
			}
		}
	}()
}

func (a *Advisor) startTickTesting(symbol string) {
	
}

// 取得K棒高點
// shift: 第幾根K棒
func (a *Advisor) High(shift int) float64 {
	return float64(shift)+1.1
}
