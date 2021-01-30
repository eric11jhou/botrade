package botrade

import (
	"fmt"
	"context"
	"strconv"
	"github.com/adshao/go-binance/v2"
	log "github.com/sirupsen/logrus"
)

var intervals = []string{"1m", "3m", "5m", "15m", "30m", "1h", "2h", "4h", "6h", "8h", "12h", "1d", "3d", "1w", "1M"}
// loadHistoryData 載入歷史數據, K棒
func (a *Advisor) loadHistoryData(symbol string) {
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

// 抓取歷史所有K線至暫存
func (a *Advisor) loadHistoryDataTesting(symbol string, startTime, endTime int64) {
	client := binance.NewClient(a.apiKey, a.secretKey)
	for _, interval := range intervals {
		startTime_ := startTime - 1000*60*60*24*30*1 // 抓取回測起始時間多久之前的K棒
		total := (endTime - startTime_) / 60000 / a.getMin(interval) 
		status := fmt.Sprintf("(%d/%d)", 0, total)
		fmt.Printf("下載 %-4sK線: %20s", interval, status)
		klinesTemp := make([]*binance.Kline, 0)
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
			klinesTemp = append(klinesTemp, klines...)
			startTime_ = klines[len(klines)-1].CloseTime + 1
			current := (endTime - startTime_) / 60000 / a.getMin(interval) 
			status := fmt.Sprintf("(%d/%d)", total-current, total)
			fmt.Printf("\r下載 %-4sK線: %20s", interval, status)
		}
		a.klineTemp[interval] = klinesTemp
		fmt.Printf("\r下載 %-4sK線: 完成，共%d筆%-30s\n", interval, len(a.klineTemp[interval]), "")
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
				a.mutex.Lock()
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
				a.mutex.Unlock()
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
				a.time = event.Time
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

// 報價
type quote struct {
	ask float64
	bid float64
	volume float64
	time int64
}

func (a *Advisor) startTickTesting(symbol string, startTime, endTime int64) {
	quoteChan := make(chan *quote)
	// 報價源，可換為其他資料
	go func(){
		for _, klineTemp := range a.klineTemp["1m"] {
			if !(klineTemp.OpenTime >= startTime && klineTemp.OpenTime <= endTime) {
				continue
			}
			ask, err := strconv.ParseFloat(klineTemp.Close, 64)
			if err != nil {
				log.Panic("quote error")
			}
			bid, err := strconv.ParseFloat(klineTemp.Close, 64)
			if err != nil {
				log.Panic("quote error")
			}
			volume, err := strconv.ParseFloat(klineTemp.Volume, 64)
			if err != nil {
				log.Panic("quote error")
			}
			quoteChan <- &quote{
				ask: ask,
				bid: bid,
				volume: volume,
				time: klineTemp.OpenTime,
			}
		}
		close(quoteChan)
	}()
	// 接收新報價後，更新數據
	go func(){
		nextIndex := make(map[string]int)
		for quote := range quoteChan {
			<- a.nextTick
			for _, interval := range intervals {
				for i := nextIndex[interval]; i < len(a.klineTemp[interval]); i++ {
					klineTemp := a.klineTemp[interval][i]
					if klineTemp.CloseTime < quote.time { // 報價之前的K棒直接加入
						a.kline[interval] = append([]*binance.Kline{klineTemp}, a.kline[interval]...)
					} else if klineTemp.OpenTime <= quote.time && klineTemp.CloseTime > quote.time { // 當前K棒，未來應支援當前K棒變動
						a.kline[interval] = append([]*binance.Kline{klineTemp}, a.kline[interval]...)
					} else {
						nextIndex[interval] = i
						break
					}
				}
			}
			// 新報價
			a.ask = quote.ask
			a.bid = quote.bid
			a.time = quote.time
			a.tick <- struct{}{}
		}
	}()
	// 每個1m收盤價(支援每個報價? 時戳:價格) 觸發tick -> 更新數據: 
	// if tick.time > kline[0].CloseTime -> 更新此K線(從暫存載入)
	// else 更新kline[0]的數據(最高最低收盤價,量,額,筆數)
}
