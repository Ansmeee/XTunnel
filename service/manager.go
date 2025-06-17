package service

import (
	"fmt"
	"log"
	"sync"
)

type TunnelManager struct {
	tunnels map[string]*Tunnel
	mutex   sync.Mutex
}

func NewTunnelManager() *TunnelManager {
	return &TunnelManager{
		tunnels: make(map[string]*Tunnel),
	}
}
func (tm *TunnelManager) AddTunnel(identifier string, config *TunnelConfig) (*Tunnel, error) {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()

	if _, ok := tm.tunnels[identifier]; ok {
		log.Printf("[%s] tunnel already exists", identifier)
		return nil, fmt.Errorf("[%s] tunnel already exists", identifier)
	}

	tunnel := NewTunnel(config)
	tunnel.identifier = identifier
	tm.tunnels[identifier] = tunnel
	log.Printf("[%s] tunnel added", identifier)
	return tunnel, nil
}

func (tm *TunnelManager) StartTunnel(identifier string) error {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()
	tunnel, ok := tm.tunnels[identifier]
	if !ok {
		return fmt.Errorf("[%s] tunnel not exists", identifier)
	}

	if tunnel.status != StatusStopped {
		return fmt.Errorf("[%s] tunnel already running", identifier)
	}

	go func() {
		if err := tunnel.Start(); err != nil {
			fmt.Printf("[%s] tunnel start error: %s", tunnel.identifier, err.Error())
		}
	}()

	return nil
}

func (tm *TunnelManager) StopTunnel(identifier string) error {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()
	tunnel, ok := tm.tunnels[identifier]
	if !ok {
		return fmt.Errorf("[%s] tunnel not exists", identifier)
	}
	tm.tunnels[identifier] = NewTunnel(&TunnelConfig{
		Username:   tunnel.config.Username,
		Password:   tunnel.config.Password,
		LocalAddr:  tunnel.config.LocalAddr,
		ServerAddr: tunnel.config.ServerAddr,
		RemoteAddr: tunnel.config.RemoteAddr,
	})

	tunnel.Stop()
	return nil
}

func (tm *TunnelManager) StatusTunnel(identifier string) (TunnelStatus, error) {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()

	tunnel, ok := tm.tunnels[identifier]
	if !ok {
		return 0, fmt.Errorf("[%s] tunnel not exists", identifier)
	}

	return tunnel.status, nil
}
