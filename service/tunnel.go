package service

import (
	"context"
	"fmt"
	"golang.org/x/crypto/ssh"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"
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

func (t *Tunnel) Start() error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.status != StatusStopped {
		return fmt.Errorf("[%s] tunnel already running or starting", t.identifier)
	}

	t.status = StatusStarting
	log.Printf("[%s] tunnel starting...", t.identifier)
	if err := t.connectSSH(); err != nil {
		t.status = StatusStopped
		return err
	}

	if err := t.listenNet(); err != nil {
		t.status = StatusStopped
		t.sshClient.Close()
		return err
	}

	t.status = StatusRunning
	t.startedAt = time.Now()
	log.Printf("[%s] ssh tunnel established %s → %s via %s", t.identifier, t.config.LocalAddr, t.config.RemoteAddr, t.config.ServerAddr)

	go t.handleSignals()
	go t.runTunnel()
	go t.monitorConnection()

	return nil
}

func (t *Tunnel) Stop() {
	t.mu.Lock()

	if t.status != StatusRunning {
		t.mu.Unlock()
		return
	}

	t.status = StatusStopping
	t.mu.Unlock()
	log.Printf("[%s] tunnel stopping...", t.identifier)

	t.cancel()

	if t.listener != nil {
		if err := t.listener.Close(); err != nil {
			log.Printf("[%s] tunnel listener close error: %s", t.identifier, err.Error())
		}
	}

	if t.sshClient != nil {
		if err := t.sshClient.Close(); err != nil {
			log.Printf("[%s] ssh client close error: %s", t.identifier, err.Error())
			return
		}
	}

	t.wg.Wait()
	t.mu.Lock()
	t.status = StatusStopped
	t.mu.Unlock()
	log.Printf("tunnel stopped: %s", t.identifier)
}

func (t *Tunnel) runTunnel() {
	sem := make(chan struct{}, 10)
	for {
		select {
		case <-t.ctx.Done():
			log.Printf("[%s] tunnel listener stopping...", t.identifier)
			return
		default:
			conn, err := t.listener.Accept()
			if err != nil {
				if !strings.Contains(err.Error(), "use of closed network connection") {
					log.Printf("[%s] tunnel accept err: %s", t.identifier, err.Error())
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
		log.Printf("[%s] remote addr dial error: %s", t.identifier, err.Error())
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
			log.Printf("[%s] local to remote server forwarding err: %s", t.identifier, err.Error())
		}
		log.Printf("[%s] local to remote server forwarding completed, transmit %d bytes", t.identifier, n)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		n, err := io.CopyBuffer(localConn, remoteConn, buf)
		if err != nil {
			log.Printf("[%s] remote server to local forwarding err: %s", t.identifier, err.Error())
		}
		log.Printf("[%s] remote server to local forwarding completed, transmit %d bytes", t.identifier, n)
	}()

	wg.Wait()
}

func (t *Tunnel) listenNet() error {
	listener, err := net.Listen("tcp", t.config.LocalAddr)
	if err != nil {
		log.Printf("[%s] failed to listen %s, err: %s", t.identifier, t.config.LocalAddr, err.Error())
		return err
	}

	t.listener = listener
	return nil
}

func (t *Tunnel) connectSSH() error {
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

		log.Printf("[%s] ssh connect err: %s, [%d/3]rtetring...", t.identifier, err.Error(), i)
		time.Sleep(1 * time.Second)
	}

	return fmt.Errorf("[%s] ssh connect after 3 attempts", t.identifier)
}

func (t *Tunnel) monitorConnection() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			_, _, err := t.sshClient.SendRequest(fmt.Sprintf("http://%s", t.config.LocalAddr), true, nil)
			if err != nil {
				log.Printf("[%s] keepalive err: %s", t.identifier, err.Error())
				t.Stop()
				return
			}
		case <-t.ctx.Done():
			return
		}
	}
}

func (t *Tunnel) handleSignals() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer signal.Stop(quit)

	select {
	case sig := <-quit:
		log.Printf("[%s] received signal: %v", t.identifier, sig)
		t.Stop()
	case <-t.ctx.Done():
	}
}
