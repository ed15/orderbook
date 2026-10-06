package wallet

import (
	"errors"
	"sync"
)

type Wallet struct {
	UserID             string  `json:"user_id"`
	BRLAvailable       float64 `json:"brl_available"`
	BRLLocked          float64 `json:"brl_locked"`
	VibraniumAvailable float64 `json:"vibranium_available"`
	VibraniumLocked    float64 `json:"vibranium_locked"`
}

type Manager struct {
	mu      sync.RWMutex
	wallets map[string]*Wallet
}

func NewManager() *Manager {
	return &Manager{
		wallets: make(map[string]*Wallet),
	}
}

func (m *Manager) GetOrCreateWallet(userID string, initialBRL, initialVib float64) *Wallet {
	m.mu.Lock()
	defer m.mu.Unlock()

	w, exists := m.wallets[userID]
	if !exists {
		w = &Wallet{
			UserID:             userID,
			BRLAvailable:       initialBRL,
			VibraniumAvailable: initialVib,
		}
		m.wallets[userID] = w
	}
	return w
}

func (m *Manager) LockBRL(userID string, amount float64) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	w, exists := m.wallets[userID]
	if !exists || w.BRLAvailable < amount {
		return errors.New("saldo insuficiente en BRL")
	}

	w.BRLAvailable -= amount
	w.BRLLocked += amount
	return nil
}

func (m *Manager) LockVibranium(userID string, amount float64) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	w, exists := m.wallets[userID]
	if !exists || w.VibraniumAvailable < amount {
		return errors.New("saldo insuficiente en Vibranium")
	}

	w.VibraniumAvailable -= amount
	w.VibraniumLocked += amount
	return nil
}

func (m *Manager) SettleTrade(buyerID, sellerID string, amount, price float64) {
	m.mu.Lock()
	defer m.mu.Unlock()

	totalBRL := amount * price
	buyer := m.wallets[buyerID]
	seller := m.wallets[sellerID]

	// Transacción del Comprador
	buyer.BRLLocked -= totalBRL
	buyer.VibraniumAvailable += amount

	// Transacción del Vendedor
	seller.VibraniumLocked -= amount
	seller.BRLAvailable += totalBRL
}

func (m *Manager) GetAllWallets() map[string]Wallet {
	m.mu.RLock()
	defer m.mu.RUnlock()

	result := make(map[string]Wallet)
	for id, w := range m.wallets {
		result[id] = *w
	}
	return result
}
