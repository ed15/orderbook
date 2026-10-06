package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	OrdersTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "vibranium_orders_total",
		Help: "Total de ordenes procesadas",
	})
	TradesTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "vibranium_trades_total",
		Help: "Total de trades ejecutados",
	})
)
