package orderbook

import (
	"fmt"
	"orderbook/pkg/wallet"
	"time"
)

type OrderType string

const (
	Buy  OrderType = "BUY"
	Sell OrderType = "SELL"
)

type Order struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Type      OrderType `json:"type"`
	Price     float64   `json:"price"`
	Amount    float64   `json:"amount"`
	Timestamp time.Time `json:"timestamp"`
}

type Trade struct {
	ID          string    `json:"id"`
	BuyOrderID  string    `json:"buy_order_id"`
	SellOrderID string    `json:"sell_order_id"`
	Price       float64   `json:"price"`
	Amount      float64   `json:"amount"`
	Timestamp   time.Time `json:"timestamp"`
}

type Engine struct {
	bids          []*Order // Compras
	asks          []*Order // Ventas
	trades        []Trade
	OrderQueue    chan Order
	WalletManager *wallet.Manager
}

func NewEngine(wm *wallet.Manager) *Engine {
	return &Engine{
		bids:          make([]*Order, 0),
		asks:          make([]*Order, 0),
		OrderQueue:    make(chan Order, 200000),
		WalletManager: wm,
	}
}

// StartWorker corre en una única goroutine secuencial
func (e *Engine) StartWorker() {
	for order := range e.OrderQueue {
		e.processOrder(order)
	}
}

func (e *Engine) processOrder(incoming Order) {
	if incoming.Type == Buy {
		e.matchBuyOrder(&incoming)
	} else {
		e.matchSellOrder(&incoming)
	}
}

func (e *Engine) matchBuyOrder(buyOrder *Order) {
	for len(e.asks) > 0 && buyOrder.Amount > 0 {
		bestAsk := e.asks[0]

		if bestAsk.Price > buyOrder.Price {
			break
		}

		tradeAmount := buyOrder.Amount
		if bestAsk.Amount < tradeAmount {
			tradeAmount = bestAsk.Amount
		}

		// Ejecutar débito y crédito
		e.WalletManager.SettleTrade(buyOrder.UserID, bestAsk.UserID, tradeAmount, bestAsk.Price)

		// Registrar la transacción realizada
		e.trades = append(e.trades, Trade{
			ID:          fmt.Sprintf("trd_%d", time.Now().UnixNano()),
			BuyOrderID:  buyOrder.ID,
			SellOrderID: bestAsk.ID,
			Price:       bestAsk.Price,
			Amount:      tradeAmount,
			Timestamp:   time.Now(),
		})

		buyOrder.Amount -= tradeAmount
		bestAsk.Amount -= tradeAmount

		if bestAsk.Amount == 0 {
			e.asks = e.asks[1:]
		}
	}

	if buyOrder.Amount > 0 {
		e.bids = append(e.bids, buyOrder)
	}
}

func (e *Engine) matchSellOrder(sellOrder *Order) {
	for len(e.bids) > 0 && sellOrder.Amount > 0 {
		bestBid := e.bids[0]

		if bestBid.Price < sellOrder.Price {
			break
		}

		tradeAmount := sellOrder.Amount
		if bestBid.Amount < tradeAmount {
			tradeAmount = bestBid.Amount
		}

		// Ejecutar débito y crédito
		e.WalletManager.SettleTrade(bestBid.UserID, sellOrder.UserID, tradeAmount, sellOrder.Price)

		sellOrder.Amount -= tradeAmount
		bestBid.Amount -= tradeAmount

		if bestBid.Amount == 0 {
			e.bids = e.bids[1:]
		}
	}

	if sellOrder.Amount > 0 {
		e.asks = append(e.asks, sellOrder)
	}
}

func (e *Engine) GetTrades() []Trade {
	return e.trades
}
