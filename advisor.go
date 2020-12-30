package botrade

type Advisor struct {}

func (b *Advisor) High(shift int) float64 {
	return float64(shift)+1.1
}