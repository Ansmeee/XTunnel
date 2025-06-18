package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"xtunnel/logger"
	"xtunnel/views"
)

func main() {
	logger.Init()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	win := views.NewWindow(ctx, cancel)
	go win.Run()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer signal.Stop(quit)

	select {
	case <-quit:
		win.Destroy(ctx)
	case <-ctx.Done():
		win.Destroy(ctx)
	}
}
