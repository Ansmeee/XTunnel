package service

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"golang.org/x/crypto/ssh"
	"io"
	"net"
	"strings"
	"sync"
	"time"
	"xtunnel/logger"
)

type TunnelStatus int

const (
	StatusStopped TunnelStatus = iota
	StatusStarting
	StatusRunning
	StatusStopping
)

type TunnelConfig struct {
	Username   string
	Password   string
	LocalAddr  string
	ServerAddr string
	RemoteAddr string
}

type Tunnel struct {
	identifier string
	status     TunnelStatus
	config     *TunnelConfig
	quit       chan struct{}
	sshClient  *ssh.Client
	listener   net.Listener
	wg         sync.WaitGroup
	mu         sync.RWMutex
	ctx        context.Context
	cancel     context.CancelFunc
	startedAt  time.Time
}

func NewTunnel(config *TunnelConfig) *Tunnel {
	ctx, cancel := context.WithCancel(context.Background())
	return &Tunnel{
		config: config,
		ctx:    ctx,
		cancel: cancel,
		status: StatusStopped,
	}
}

func (t *Tunnel) Start(ctx context.Context) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.status != StatusStopped {
		logger.Error(ctx, "tunnel already running or starting", g.Map{"identifier": t.identifier, "status": t.status})
		return fmt.Errorf("[%s] tunnel already running or starting", t.identifier)
	}

	t.status = StatusStarting
	logger.Info(ctx, "tunnel starting", g.Map{"identifier": t.identifier})
	if err := t.connectSSH(ctx); err != nil {
		t.status = StatusStopped
		return err
	}

	if err := t.listenNet(ctx); err != nil {
		t.status = StatusStopped
		t.sshClient.Close()
		return err
	}

	t.status = StatusRunning
	t.startedAt = time.Now()

	logger.Info(ctx, "ssh tunnel established", g.Map{
		"identifier": t.identifier,
		"localAddr":  t.config.LocalAddr,
		"remoteAddr": t.config.RemoteAddr,
		"serverAddr": t.config.ServerAddr,
	})

	go t.runTunnel(ctx)
	//go t.monitorConnection(ctx)

	return nil
}

func (t *Tunnel) Stop(ctx context.Context) {
	t.mu.Lock()

	if t.status != StatusRunning {
		t.mu.Unlock()
		return
	}

	logger.Info(ctx, "tunnel stopping", g.Map{"identifier": t.identifier})
	t.status = StatusStopping
	t.mu.Unlock()
	t.cancel()

	if t.listener != nil {
		if err := t.listener.Close(); err != nil {
			logger.Error(ctx, "tunnel listener close error", g.Map{"identifier": t.identifier, "err": err.Error()})
		}
	}

	if t.sshClient != nil {
		if err := t.sshClient.Close(); err != nil {
			logger.Error(ctx, "ssh client close error", g.Map{"identifier": t.identifier, "err": err.Error()})
			return
		}
	}

	t.wg.Wait()
	t.mu.Lock()
	t.status = StatusStopped
	t.mu.Unlock()

	logger.Info(ctx, "tunnel stopped", g.Map{"identifier": t.identifier})
}

func (t *Tunnel) runTunnel(ctx context.Context) {
	sem := make(chan struct{}, 20)
	for {
		select {
		case <-t.ctx.Done():
			return
		default:
			conn, err := t.listener.Accept()
			if err != nil {
				if !strings.Contains(err.Error(), "use of closed network connection") {
					logger.Error(ctx, "tunnel accept error", g.Map{"identifier": t.identifier, "err": err.Error()})
				}
				continue
			}

			select {
			case sem <- struct{}{}:
				t.wg.Add(1)
				go func(c net.Conn) {
					defer func() {
						<-sem
						t.wg.Done()
					}()
					t.forward(t.ctx, conn)
				}(conn)
			case <-t.ctx.Done():
				conn.Close()
				return
			}
		}
	}
}

func (t *Tunnel) forward(ctx context.Context, localConn net.Conn) {
	remoteConn, err := t.sshClient.Dial("tcp", t.config.RemoteAddr)
	if err != nil {
		logger.Error(ctx, "remote addr dial error", g.Map{"identifier": t.identifier, "err": err.Error()})
		return
	}
	defer remoteConn.Close()

	buf := make([]byte, 32*1024)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		n, err := io.CopyBuffer(remoteConn, localConn, buf)
		if err != nil {
			logger.Error(ctx, "local to remote server forwarding err", g.Map{"identifier": t.identifier, "err": err.Error()})
			return
		}
		logger.Info(ctx, "local to remote server forwarding completed", g.Map{"identifier": t.identifier, "transmit": n})
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		n, err := io.CopyBuffer(localConn, remoteConn, buf)
		if err != nil {
			logger.Error(ctx, "remote server to local forwarding err", g.Map{"identifier": t.identifier, "err": err.Error()})
			return
		}
		logger.Info(ctx, "local to remote server forwarding completed", g.Map{"identifier": t.identifier, "transmit": n})
	}()

	wg.Wait()
}

func (t *Tunnel) listenNet(ctx context.Context) error {
	listener, err := net.Listen("tcp", t.config.LocalAddr)
	if err != nil {
		logger.Error(ctx, "listen error", g.Map{
			"identifier": t.identifier,
			"localAddr":  t.config.LocalAddr,
			"err":        err.Error(),
		})
		return err
	}

	t.listener = listener
	return nil
}

func (t *Tunnel) connectSSH(ctx context.Context) error {
	config := &ssh.ClientConfig{
		User:            t.config.Username,
		Auth:            []ssh.AuthMethod{ssh.Password(t.config.Password)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         30 * time.Second,
	}

	var err error
	var sshClient *ssh.Client

	for i := 1; i <= 3; i++ {
		sshClient, err = ssh.Dial("tcp", t.config.ServerAddr, config)
		if err == nil {
			t.sshClient = sshClient
			return nil
		}

		logger.Error(ctx, "ssh connect error", g.Map{"identifier": t.identifier, "err": err.Error(), "retry": i})
		time.Sleep(1 * time.Second)
	}

	return fmt.Errorf("[%s] ssh connect after 3 attempts", t.identifier)
}

func (t *Tunnel) monitorConnection(ctx context.Context) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			_, _, err := t.sshClient.SendRequest(fmt.Sprintf("http://%s", t.config.LocalAddr), true, nil)
			if err != nil {
				logger.Error(ctx, "keepalive err", g.Map{"identifier": t.identifier, "err": err.Error()})
				t.Stop(ctx)
				return
			}
		case <-t.ctx.Done():
			return
		}
	}
}
