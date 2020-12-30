package botrade

// Advisor 可取得各種訊息與交易功能
type Advisor struct {}

// 取得K棒高點
// shift: 第幾根K棒
func (b *Advisor) High(shift int) float64 {
	return float64(shift)+1.1
}