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
	"sync"
	"syscall"
	"time"
)

type TunnelConfig struct {
	Username   string
	Password   string
	LocalAddr  string
	ServerAddr string
	RemoteAddr string
}

type Tunnel struct {
	config    *TunnelConfig
	quit      chan struct{}
	sshClient *ssh.Client
	wg        sync.WaitGroup
}

func NewTunnel(config *TunnelConfig) *Tunnel {
	return &Tunnel{config: config, quit: make(chan struct{})}
}

func (t *Tunnel) Stop() {
	close(t.quit)
	t.wg.Wait()

	if t.sshClient != nil {
		if err := t.sshClient.Close(); err != nil {
			log.Printf("ssh client close error: %s", err.Error())
			return
		}
	}

	log.Println("tunnel stopped")
}

func (t *Tunnel) Start() error {
	t.wg.Add(1)
	defer t.wg.Done()

	if err := t.connectSSH(); err != nil {
		return err
	}

	listener, err := net.Listen("tcp", t.config.LocalAddr)
	if err != nil {
		log.Printf("failed to listen %s, err: %s", t.config.LocalAddr, err.Error())
		return err
	}
	defer listener.Close()

	log.Printf("ssh tunnel established successfully %s → %s via %s", t.config.LocalAddr, t.config.RemoteAddr, t.config.ServerAddr)

	go t.handleSignals()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sem := make(chan struct{}, 10)
	for {
		select {
		case <-t.quit:
			log.Println("ssh server stopped...")
			cancel()
			return nil
		default:
			conn, err := listener.Accept()
			if err != nil {
				log.Printf("ssh server accept err: %s", err.Error())
				continue
			}

			sem <- struct{}{}
			go func(c net.Conn) {
				defer func() { <-sem }()
				t.forward(ctx, conn)
			}(conn)
		}
	}
}

func (t *Tunnel) forward(ctx context.Context, localConn net.Conn) {
	t.wg.Add(1)
	defer t.wg.Done()
	defer localConn.Close()

	var err error
	var remoteConn net.Conn

	log.Printf("remote server connecting at %s...", t.config.RemoteAddr)

	remoteConn, err = t.sshClient.Dial("tcp", t.config.RemoteAddr)
	if err != nil {
		log.Printf("ssh client dial error: %s", err.Error())
		return
	}
	defer remoteConn.Close()

	log.Printf("成功连接到目标: %s", t.config.RemoteAddr)

	buf := make([]byte, 32*1024)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		n, err := io.CopyBuffer(remoteConn, localConn, buf)
		if err != nil {
			log.Printf("local to remote server forwarding err: %s", err.Error())
		}
		log.Printf("local to remote server forwarding completed, transmit %d bytes", n)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		n, err := io.CopyBuffer(localConn, remoteConn, buf)
		if err != nil {
			log.Printf("remote server to local forwarding err: %s", err.Error())
		}
		log.Printf("remote server to local forwarding completed, transmit %d bytes", n)
	}()

	wg.Wait()
}

func (t *Tunnel) handleSignals() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	t.Stop()
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
			break
		}
		log.Printf("ssh connect err: %s, [%d/3]rtetring...", err.Error(), i)
		time.Sleep(1 * time.Second)
	}

	if err != nil {
		log.Printf("ssh connect err: %s", err.Error())
		return err
	}

	// 增加保活和心跳检测机制
	go t.keepAlive()
	t.sshClient = sshClient
	return nil
}

func (t *Tunnel) keepAlive() {
	t.wg.Add(1)
	defer t.wg.Done()

	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	retryCount := 0
	maxRetries := 5

	for {
		select {
		case <-ticker.C:
			_, _, err := t.sshClient.SendRequest(fmt.Sprintf("http://%s", t.config.LocalAddr), true, nil)
			if err == nil {
				retryCount = 0
				continue
			}

			if retryCount > maxRetries {
				log.Printf("ssh keep alive err: exceed the maximum number of retries")
				t.Stop()
				return
			}

			retryCount++
			if err = t.connectSSH(); err != nil {
				log.Printf("ssh reconnect err: %s", err.Error())
				continue
			}

			retryCount = 0
			log.Printf("ssh reconnect success")
			return
		case <-t.quit:
			log.Println("ssh server stopped...")
			return
		}
	}
}
