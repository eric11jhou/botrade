package botrade

import (
	"github.com/adshao/go-binance/v2"
)

type OrderCreateRequest struct {
	symbol string
	sideType binance.SideType
	orderType binance.OrderType
	timeInForce *binance.TimeInForceType
	quantity *string
	price *string
	stopPrice *string
}

// 下單
func (a *Advisor) OrderCreate(orderCreateRequest *OrderCreateRequest) (*binance.CreateOrderResponse, error) {
	if a.trade {
		return a.orderCreate(orderCreateRequest)
	} else {
		return a.orderCreateTesting(orderCreateRequest)
	}
}

func (a *Advisor) orderCreate(orderCreateRequest *OrderCreateRequest) (*binance.CreateOrderResponse, error) {
	return nil, nil
}

func (a *Advisor) orderCreateTesting(orderCreateRequest *OrderCreateRequest) (*binance.CreateOrderResponse, error) {
	return nil, nil
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
	return nil, nil
}

func (a *Advisor) orderListTesting(symbol string) ([]*binance.Order, error) {
	return nil, nil
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
	return nil, nil
}

func (a *Advisor) orderGetTesting(symbol string, orderId int64) (*binance.Order, error) {
	return nil, nil
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
	return nil, nil
}