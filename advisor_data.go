package botrade

import (
	"fmt"
	"context"
	"strconv"
	"github.com/adshao/go-binance/v2"
	log "github.com/sirupsen/logrus"
)

// loadHistoryData 載入歷史數據, K棒
// interval 1m, 3m, 5m, 15m, 30m, 1h, 2h, 4h, 6h, 8h, 12h, 1d, 3d, 1w, 1M
func (a *Advisor) loadHistoryData(symbol string) {
	intervals := []string{"1m", "3m", "5m", "15m", "30m", "1h", "2h", "4h", "6h", "8h", "12h", "1d", "3d", "1w", "1M"}
	client := binance.NewClient(a.apiKey, a.secretKey)
	for _, interval := range intervals {
		klines, err := client.NewKlinesService().
		Symbol(symbol).
		Interval(interval).
		Limit(1000).
		Do(context.Background())
		if err != nil {
			log.Fatal(err)
		}
		for i := len(klines) - 1; i >= 0; i-- {
			a.kline[interval] = append(a.kline[interval] , klines[i])
		}
	}
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
