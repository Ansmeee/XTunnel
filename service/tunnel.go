package service

import (
	"golang.org/x/crypto/ssh"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Tunnel struct {
	Username   string
	Password   string
	LocalAddr  string
	ServerAddr string
	RemoteAddr string
	sshClient  *ssh.Client
	Quit       chan struct{}
}

func (t *Tunnel) Stop() {
	t.Quit <- struct{}{}
}

func (t *Tunnel) Start() {
	if err := t.connectSSH(); err != nil {
		log.Println("ssh connect err:", err.Error())
		return
	}

	listener, err := net.Listen("tcp", t.LocalAddr)
	if err != nil {
		log.Printf("Failed to listen: %s", err.Error())
		return
	}
	defer listener.Close()

	log.Printf("SSH隧道已建立 %s → %s via %s", t.LocalAddr, t.RemoteAddr, t.ServerAddr)

	go t.handleSignals()

	for {
		select {
		case <-t.Quit:
			log.Println("ssh server stop")
			return
		default:
			log.Println("ssh server listening")
			conn, err := listener.Accept()
			if err != nil {
				log.Printf("ssh server accept err:", err.Error())
				continue
			}
			log.Printf("接受来自 %s 的新连接", conn.RemoteAddr())
			go t.forward(conn)
		}
	}
}

func (t *Tunnel) forward(localConn net.Conn) {
	defer localConn.Close()

	log.Printf("尝试连接到目标: %s", t.RemoteAddr)
	remoteConn, err := t.sshClient.Dial("tcp", t.RemoteAddr)
	if err != nil {
		log.Printf("remote server connect err:", err.Error())
		return
	}
	defer remoteConn.Close()

	log.Printf("成功连接到目标: %s", t.RemoteAddr)

	done := make(chan struct{}, 2)
	go func() {
		n, err := io.Copy(remoteConn, localConn)
		if err != nil {
			log.Printf("remote server copy err:", err.Error())
		}
		log.Printf("本地→远程转发完成: 传输 %d 字节", n)
		done <- struct{}{}
	}()

	go func() {
		n, err := io.Copy(localConn, remoteConn)
		if err != nil {
			log.Printf("remote server copy err:", err.Error())
		}
		log.Printf("远程→本地转发完成: 传输 %d 字节", n)
		done <- struct{}{}
	}()

	<-done

	log.Printf("连接 %s ↔ %s 已完成", localConn.RemoteAddr(), t.RemoteAddr)
}

func (t *Tunnel) handleSignals() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("收到终止信号，关闭隧道...")
	close(t.Quit)
	if t.sshClient != nil {
		t.sshClient.Close()
	}
	os.Exit(0)
}

func (t *Tunnel) connectSSH() error {
	config := &ssh.ClientConfig{
		User:            t.Username,
		Auth:            []ssh.AuthMethod{ssh.Password(t.Password)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         30 * time.Second,
	}

	sshClient, err := ssh.Dial("tcp", t.ServerAddr, config)
	if err != nil {
		return err
	}

	t.sshClient = sshClient
	return nil
}
