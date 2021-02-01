package botrade

import (
	"fmt"
	"strconv"
	"github.com/adshao/go-binance/v2"
)

type OrderCreateRequest struct {
	Symbol string
	SideType binance.SideType
	OrderType binance.OrderType
	TimeInForce *binance.TimeInForceType
	Quantity *float64
	Price *float64
	StopPrice *float64
}

// 下單
func (a *Advisor) OrderCreate(o *OrderCreateRequest) (*binance.CreateOrderResponse, error) {
	if a.trade {
		return a.orderCreate(o)
	} else {
		return a.orderCreateTesting(o)
	}
}

func (a *Advisor) orderCreate(o *OrderCreateRequest) (*binance.CreateOrderResponse, error) {
	return nil, nil
}

func (a *Advisor) orderCreateTesting(o *OrderCreateRequest) (*binance.CreateOrderResponse, error) {
	switch o.OrderType {
	case binance.OrderTypeLimit:
		if o.Price == nil {
			return nil, fmt.Errorf("price required")
		}
		switch o.SideType {
		case binance.SideTypeBuy:
			if a.ask >= *o.Price {
				return nil, fmt.Errorf("invalid price")
			}
		case binance.SideTypeSell:
			if a.bid <= *o.Price {
				return nil, fmt.Errorf("invalid price")
			}
		}
	case binance.OrderTypeMarket:
	case binance.OrderTypeLimitMaker:
	case binance.OrderTypeStopLoss:
		if o.StopPrice == nil {
			return nil, fmt.Errorf("stop price required")
		}
		switch o.SideType {
		case binance.SideTypeBuy:
			if a.ask <= *o.StopPrice {
				return nil, fmt.Errorf("invalid stop price")
			}
		case binance.SideTypeSell:
			if a.bid >= *o.StopPrice {
				return nil, fmt.Errorf("invalid stop price")
			}
		}
	case binance.OrderTypeStopLossLimit:
		if o.Price == nil {
			return nil, fmt.Errorf("price required")
		}
		if o.StopPrice == nil {
			return nil, fmt.Errorf("stop price required")
		}
		switch o.SideType {
		case binance.SideTypeBuy:
			if a.ask <= *o.StopPrice {
				return nil, fmt.Errorf("invalid stop price")
			}
			if *o.Price <= *o.StopPrice {
				return nil, fmt.Errorf("invalid price")
			}
		case binance.SideTypeSell:
			if a.bid >= *o.StopPrice {
				return nil, fmt.Errorf("invalid stop price")
			}
			if *o.Price >= *o.StopPrice {
				return nil, fmt.Errorf("invalid price")
			}
		}
	case binance.OrderTypeTakeProfit:
		if o.StopPrice == nil {
			return nil, fmt.Errorf("stop price required")
		}
		switch o.SideType {
		case binance.SideTypeBuy:
			if a.ask >= *o.StopPrice {
				return nil, fmt.Errorf("invalid stop price")
			}
		case binance.SideTypeSell:
			if a.bid <= *o.StopPrice {
				return nil, fmt.Errorf("invalid stop price")
			}
		}
	case binance.OrderTypeTakeProfitLimit:
		if o.Price == nil {
			return nil, fmt.Errorf("price required")
		}
		if o.StopPrice == nil {
			return nil, fmt.Errorf("stop price required")
		}
		switch o.SideType {
		case binance.SideTypeBuy:
			if a.ask >= *o.StopPrice {
				return nil, fmt.Errorf("invalid stop price")
			}
			if *o.Price <= *o.StopPrice {
				return nil, fmt.Errorf("invalid price")
			}
		case binance.SideTypeSell:
			if a.bid <= *o.StopPrice {
				return nil, fmt.Errorf("invalid stop price")
			}
			if *o.Price >= *o.StopPrice {
				return nil, fmt.Errorf("invalid price")
			}
		}
	}
	a.orderIDCount++
	orderRes := &binance.CreateOrderResponse{
		Symbol: o.Symbol,
		OrderID: a.orderIDCount,
		Type: o.OrderType,
		Side: o.SideType,
	}
	order := &binance.Order{
		Symbol: o.Symbol,
		OrderID: a.orderIDCount,
		Status: binance.OrderStatusTypeNew,
		Type: o.OrderType,
		Side: o.SideType,
		Time: a.time,
	}
	if o.Quantity != nil {
		orderRes.OrigQuantity = strconv.FormatFloat(*o.Quantity, 'f', -1, 64)
		order.OrigQuantity = orderRes.OrigQuantity
	}
	if o.Price != nil {
		orderRes.Price = strconv.FormatFloat(*o.Price, 'f', -1, 64)
		order.Price = orderRes.Price
	}
	if o.TimeInForce != nil {
		orderRes.TimeInForce = *o.TimeInForce
		order.TimeInForce = *o.TimeInForce
	}
	if o.StopPrice != nil {
		order.StopPrice = strconv.FormatFloat(*o.StopPrice, 'f', -1, 64)
	}
	a.openOrders = append([]*binance.Order{order}, a.openOrders...)
	return orderRes, nil
}

// 取消訂單
func (a *Advisor) OrderCancel(symbol string, orderId int64) (*binance.CancelOrderResponse, error) {
	if a.trade {
		return a.orderCancel(symbol, orderId)
	} else {
		return a.orderCancelTesting(symbol, orderId)
	}
}

func (a *Advisor) orderCancel(symbol string, orderId int64) (*binance.CancelOrderResponse, error) {
	return nil, nil
}

func (a *Advisor) orderCancelTesting(symbol string, orderId int64) (*binance.CancelOrderResponse, error) {
	for i := 0; i < len(a.openOrders); i++ {
		order := a.openOrders[i]
		if order.Symbol == symbol && order.OrderID == orderId {
			orderRes := &binance.CancelOrderResponse{
				Symbol: order.Symbol,
				OrderID: order.OrderID,
				ClientOrderID: order.ClientOrderID,
				TransactTime: order.Time,
				Price: order.Price,
				OrigQuantity: order.OrigQuantity,
				ExecutedQuantity: order.ExecutedQuantity,
				CummulativeQuoteQuantity: order.CummulativeQuoteQuantity,
				Status: order.Status,
				TimeInForce: order.TimeInForce,
				Type: order.Type,
				Side: order.Side,
			}
			a.openOrders = append(a.openOrders[:i], a.openOrders[i+1:]...)
			i--
			return orderRes, nil
		}
	}
	return nil, fmt.Errorf("order not exists")
}

// 取得所有訂單
func (a *Advisor) OrderList(symbol string) ([]*binance.Order, error) {
	if a.trade {
		return a.orderList(symbol)
	} else {
		return a.orderListTesting(symbol)
	}
}

func (a *Advisor) orderList(symbol string) ([]*binance.Order, error) {
	return a.orders, nil
}

func (a *Advisor) orderListTesting(symbol string) ([]*binance.Order, error) {
	return a.orders, nil
}

// 取得所有掛單
func (a *Advisor) OpenOrderList(symbol string) ([]*binance.Order, error) {
	if a.trade {
		return a.openOrderList(symbol)
	} else {
		return a.openOrderListTesting(symbol)
	}
}

func (a *Advisor) openOrderList(symbol string) ([]*binance.Order, error) {
	return a.openOrders, nil
}

func (a *Advisor) openOrderListTesting(symbol string) ([]*binance.Order, error) {
	return a.openOrders, nil
}

// 取得訂單
func (a *Advisor) OrderGet(symbol string, orderId int64) (*binance.Order, error) {
	if a.trade {
		return a.orderGet(symbol, orderId)
	} else {
		return a.orderGetTesting(symbol, orderId)
	}
}

func (a *Advisor) orderGet(symbol string, orderId int64) (*binance.Order, error) {                                             
	for _, order := range a.orders {
		if order.Symbol == symbol && order.OrderID == orderId {
			return order, nil
		}
	}
	for _, order := range a.openOrders {
		if order.Symbol == symbol && order.OrderID == orderId {
			return order, nil
		}
	}
	return nil, fmt.Errorf("not exists")
}

func (a *Advisor) orderGetTesting(symbol string, orderId int64) (*binance.Order, error) {
	for _, order := range a.orders {
		if order.Symbol == symbol && order.OrderID == orderId {
			return order, nil
		}
	}
	for _, order := range a.openOrders {
		if order.Symbol == symbol && order.OrderID == orderId {
			return order, nil
		}
	}
	return nil, fmt.Errorf("not exists")
}

// **回測用 檢查訂單是否進場
func (a *Advisor) orderCheckExecTesting() {
	for i := 0; i < len(a.openOrders); i++ {
		o := a.openOrders[i]
		oPrice, _ := strconv.ParseFloat(o.Price, 64)
		oStopPrice, _ := strconv.ParseFloat(o.StopPrice, 64)
		switch o.Type {
		case binance.OrderTypeLimit:
			switch o.Side {
			case binance.SideTypeBuy:
				if a.ask <= oPrice {
					o.Price = strconv.FormatFloat(a.ask, 'f', -1, 64)
					o.ExecutedQuantity = o.OrigQuantity
					o.Status = binance.OrderStatusTypeFilled
					o.UpdateTime = a.time
					a.orders = append([]*binance.Order{o}, a.orders...)
					a.openOrders = append(a.openOrders[:i], a.openOrders[i+1:]...)
					i--
				}
			case binance.SideTypeSell:
				if a.bid >= oPrice {
					o.Price = strconv.FormatFloat(a.bid, 'f', -1, 64)
					o.ExecutedQuantity = o.OrigQuantity
					o.Status = binance.OrderStatusTypeFilled
					o.UpdateTime = a.time
					a.orders = append([]*binance.Order{o}, a.orders...)
					a.openOrders = append(a.openOrders[:i], a.openOrders[i+1:]...)
					i--
				}
			}
		case binance.OrderTypeMarket:
			switch  o.Side {
			case binance.SideTypeBuy:
				o.Price = strconv.FormatFloat(a.ask, 'f', -1, 64)
			case binance.SideTypeSell:
				o.Price = strconv.FormatFloat(a.bid, 'f', -1, 64)
			}
			o.ExecutedQuantity = o.OrigQuantity
			o.Status = binance.OrderStatusTypeFilled
			o.UpdateTime = a.time
			a.orders = append([]*binance.Order{o}, a.orders...)
			a.openOrders = append(a.openOrders[:i], a.openOrders[i+1:]...)
			i--
		case binance.OrderTypeLimitMaker:
		case binance.OrderTypeStopLoss:
			switch  o.Side {
			case binance.SideTypeBuy:
				if a.ask < oStopPrice {
					o.Type = binance.OrderTypeMarket
					o.UpdateTime = a.time
				}
			case binance.SideTypeSell:
				if a.bid > oStopPrice {
					o.Type = binance.OrderTypeMarket
					o.UpdateTime = a.time
				}
			}
		case binance.OrderTypeStopLossLimit:
			switch  o.Side {
			case binance.SideTypeBuy:
				if a.ask < oStopPrice {
					o.Type = binance.OrderTypeLimit
					o.UpdateTime = a.time
				}
			case binance.SideTypeSell:
				if a.bid > oStopPrice {
					o.Type = binance.OrderTypeLimit
					o.UpdateTime = a.time
				}
			}
		case binance.OrderTypeTakeProfit:
			switch  o.Side {
			case binance.SideTypeBuy:
				if a.ask >= oStopPrice {
					o.Type = binance.OrderTypeMarket
					o.UpdateTime = a.time
				}
			case binance.SideTypeSell:
				if a.bid <= oStopPrice {
					o.Type = binance.OrderTypeMarket
					o.UpdateTime = a.time
				}
			}
		case binance.OrderTypeTakeProfitLimit:
			switch  o.Side {
			case binance.SideTypeBuy:
				if a.ask >= oStopPrice {
					o.Type = binance.OrderTypeLimit
					o.UpdateTime = a.time
				}
			case binance.SideTypeSell:
				if a.bid <= oStopPrice {
					o.Type = binance.OrderTypeLimit
					o.UpdateTime = a.time
				}
			}
		}
	}
}












