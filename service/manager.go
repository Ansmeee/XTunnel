package service

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"sync"
	"xtunnel/logger"
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
func (tm *TunnelManager) AddTunnel(ctx context.Context, identifier string, config *TunnelConfig) (*Tunnel, error) {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()

	if _, ok := tm.tunnels[identifier]; ok {
		logger.Error(ctx, "tunnel already exists", g.Map{"identifier": identifier})
		return nil, fmt.Errorf("[%s] tunnel already exists", identifier)
	}

	tunnel := NewTunnel(config)
	tunnel.identifier = identifier
	tm.tunnels[identifier] = tunnel
	logger.Info(ctx, "tunnel added", g.Map{"identifier": identifier})
	return tunnel, nil
}

func (tm *TunnelManager) StartTunnel(ctx context.Context, identifier string) error {
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
		if err := tunnel.Start(ctx); err != nil {
			logger.Error(ctx, "tunnel start error", g.Map{"identifier": tunnel.identifier, "err": err.Error()})
		}
	}()

	return nil
}

func (tm *TunnelManager) StopTunnel(ctx context.Context, identifier string) error {
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

	tunnel.Stop(ctx)
	return nil
}

func (tm *TunnelManager) StopAll(ctx context.Context) {
	for _, tunnel := range tm.tunnels {
		tunnel.Stop(ctx)
	}
}

func (tm *TunnelManager) StatusTunnel(ctx context.Context, identifier string) (TunnelStatus, error) {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()

	tunnel, ok := tm.tunnels[identifier]
	if !ok {
		return 0, fmt.Errorf("[%s] tunnel not exists", identifier)
	}

	return tunnel.status, nil
}
