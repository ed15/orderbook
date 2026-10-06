package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"orderbook/pkg/orderbook"
	"orderbook/pkg/wallet"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Application struct {
	walletMgr *wallet.Manager
	engine    *orderbook.Engine
}

type OrderRequest struct {
	UserID string              `json:"user_id"`
	Type   orderbook.OrderType `json:"type"` // "BUY" o "SELL"
	Price  float64             `json:"price"`
	Amount float64             `json:"amount"`
}

func main() {
	// 1. Inicializar Administrador de Billeteras y Motor
	wm := wallet.NewManager()
	engine := orderbook.NewEngine(wm)

	// Cargar usuarios de prueba con saldo inicial
	wm.GetOrCreateWallet("user_buyer", 1000000.0, 0.0) // 1,000,000 BRL
	wm.GetOrCreateWallet("user_seller", 0.0, 10000.0)  // 10,000 Vibranium

	// 2. Iniciar el Worker del Matching Engine en una Goroutine dedicada
	go engine.StartWorker()

	app := &Application{
		walletMgr: wm,
		engine:    engine,
	}

	// 3. Definir Rutas HTTP
	mux := http.NewServeMux()
	mux.HandleFunc("POST /orders", app.createOrderHandler)
	mux.HandleFunc("GET /wallets", app.getWalletsHandler)
	mux.HandleFunc("GET /trades", app.getTradesHandler)
	mux.HandleFunc("GET /health", app.healthHandler)
	mux.Handle("GET /metrics", promhttp.Handler())

	server := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	fmt.Println("Servidor del Order Book escuchando en http://localhost:8080")
	if err := server.ListenAndServe(); err != nil {
		fmt.Printf("Error al iniciar el servidor: %v\n", err)
	}
}

func (app *Application) healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"UP"}`))
}

func (app *Application) createOrderHandler(w http.ResponseWriter, r *http.Request) {
	var req OrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Payload JSON inválido", http.StatusBadRequest)
		return
	}

	// Validaciones básicas
	if req.Amount <= 0 || req.Price <= 0 {
		http.Error(w, "El monto y el precio deben ser mayores a cero", http.StatusBadRequest)
		return
	}

	// Intentar retener el saldo antes de enviar al Order Book
	if req.Type == orderbook.Buy {
		totalCost := req.Amount * req.Price
		if err := app.walletMgr.LockBRL(req.UserID, totalCost); err != nil {
			http.Error(w, fmt.Sprintf("Error en Billetera: %v", err), http.StatusBadRequest)
			return
		}
	} else if req.Type == orderbook.Sell {
		if err := app.walletMgr.LockVibranium(req.UserID, req.Amount); err != nil {
			http.Error(w, fmt.Sprintf("Error en Billetera: %v", err), http.StatusBadRequest)
			return
		}
	} else {
		http.Error(w, "Tipo de orden no válido (Usar 'BUY' o 'SELL')", http.StatusBadRequest)
		return
	}

	// Construir la orden
	order := orderbook.Order{
		ID:        fmt.Sprintf("ord_%d", time.Now().UnixNano()),
		UserID:    req.UserID,
		Type:      req.Type,
		Price:     req.Price,
		Amount:    req.Amount,
		Timestamp: time.Now(),
	}

	// Depositar en la cola en memoria (Channel)
	app.engine.OrderQueue <- order

	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]string{
		"status":   "ACCEPTED",
		"order_id": order.ID,
	})
}

func (app *Application) getWalletsHandler(w http.ResponseWriter, r *http.Request) {
	wallets := app.walletMgr.GetAllWallets()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(wallets)
}

func (app *Application) getTradesHandler(w http.ResponseWriter, r *http.Request) {
	trades := app.engine.GetTrades()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(trades)
}
