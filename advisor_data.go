package botrade

import (
	"fmt"
	"time"
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

func (a *Advisor) loadHistoryDataTesting(symbol string, startTime, endTime int64) {
	intervals := []string{"1m", "3m", "5m", "15m", "30m", "1h", "2h", "4h", "6h", "8h", "12h", "1d", "3d", "1w", "1M"}
	client := binance.NewClient(a.apiKey, a.secretKey)
	
	for _, interval := range intervals {
		total := (endTime - startTime) / 60000 / a.getMin(interval) 
		startTime_ := startTime
		status := fmt.Sprintf("(%d/%d)", 0, total)
		fmt.Printf("下載 %s K線: %20s", interval, status)
		klinesTemp := make([]*binance.Kline, 0)
		ticker := time.Tick(time.Second)
		for {
			klines, err := client.NewKlinesService().
			Symbol(symbol).
			Interval(interval).
			StartTime(startTime_).
			EndTime(endTime).
			Limit(1000).
			Do(context.Background())
			if err != nil {
				log.Fatal(err)
			}
			if len(klines) == 0 {
				break
			}
			for _ , kline := range klines {
				klinesTemp = append([]*binance.Kline{kline}, klinesTemp...)
			}
			startTime_ = klines[len(klines)-1].CloseTime + 1
			current := (endTime - startTime_) / 60000 / a.getMin(interval) 
			status := fmt.Sprintf("(%d/%d)", total-current, total)
			fmt.Printf("\r下載 %s K線: %20s", interval, status)
			<- ticker
		}
		a.kline[interval] = klinesTemp
		fmt.Printf("\r下載 %s K線: 完成，共%d筆%-30s\n", interval, len(a.kline[interval]), "")
	}
}

func (a *Advisor) getMin(interval string) int64 {
	switch interval {
	case "1m":
		return 1;
	case "3m":
		return 3;
	case "5m":
		return 5;
	case "15m":
		return 15;
	case "30m":
		return 30;
	case "1h":
		return 60;
	case "2h":
		return 120;
	case "4h":
		return 240;
	case "6h":
		return 360;
	case "8h":
		return 480;
	case "12h":
		return 720;
	case "1d":
		return 1440;
	case "3d":
		return 1440*3;
	case "1w":
		return 1440*7;
	case "1M":
		return 1440*30;
	default:
		return 1;
	}
}

func (a *Advisor) startTick(symbol string) {
	go func(){
		{
			wsKlineHandler := func(event *binance.WsKlineEvent) {
				for k, v := range a.kline {
					if event.Kline.StartTime > v[0].CloseTime { // 此interval已收盤
						client := binance.NewClient(a.apiKey, a.secretKey)
						klines, err := client.NewKlinesService().
						Symbol(symbol).
						Interval(k).
						Limit(5).
						Do(context.Background())
						if err != nil {
							log.Error(err)
							continue
						}
						for _, newKline := range klines {
							if newKline.OpenTime > v[0].CloseTime {
								a.kline[k] = append([]*binance.Kline{newKline}, v...)
							}
							for i, kline := range v {
								if kline.OpenTime == newKline.OpenTime {
									v[i] = newKline
								}
							}
						}
					} else if event.Kline.StartTime >= v[0].OpenTime { // 此interval目前K棒尚未收盤
						// 更新收盤價
						v[0].Close = event.Kline.Close
						// 更新最高價
						if newHigh, err := strconv.ParseFloat(event.Kline.High, 64); err == nil {
							if high, err := strconv.ParseFloat(v[0].High, 64); err == nil && high < newHigh {
								v[0].High = event.Kline.High
							}
						}
						// 更新最低價
						if newLow, err := strconv.ParseFloat(event.Kline.Low, 64); err == nil {
							if low, err := strconv.ParseFloat(v[0].Low, 64); err == nil && low > newLow {
								v[0].Low = event.Kline.Low
							}
						}
						// 更新交易量
						if event.Kline.IsFinal {
							if event.Kline.EndTime != v[0].CloseTime {
								// 成交量
								if newVolume, err := strconv.ParseFloat(event.Kline.Volume, 64); err == nil {
									if volume, err := strconv.ParseFloat(v[0].Volume, 64); err == nil {
										v[0].Volume = strconv.FormatFloat(volume + newVolume, 'f', -1, 64)
									}
								}
								// 成交額
								if newQuoteVolume, err := strconv.ParseFloat(event.Kline.QuoteVolume, 64); err == nil {
									if quoteVolume, err := strconv.ParseFloat(v[0].QuoteAssetVolume, 64); err == nil {
										v[0].QuoteAssetVolume = strconv.FormatFloat(quoteVolume + newQuoteVolume, 'f', -1, 64)
									}
								}
								// 成交筆數
								v[0].TradeNum += event.Kline.TradeNum
							}
						}
					}
				}
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

func (a *Advisor) startTickTesting(symbol string, startTime, endTime int64) {
	// 抓取歷史所有K線(先載入起始報價之前數據,其餘暫存)
	// 每個1m收盤價(支援每個報價? 時戳:價格) 觸發tick -> 更新數據: 
	// if tick.time > kline[0].CloseTime -> 更新此K線(從暫存載入)
	// else 更新kline[0]的數據(最高最低收盤價,量,額,筆數)
}
