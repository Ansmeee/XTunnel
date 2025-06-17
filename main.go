package main

import (
	"os"
	"os/signal"
	"syscall"
	"xtunnel/logger"
	"xtunnel/views"
)

func main() {
	logger.Init()

	views.NewWindow().Run()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer signal.Stop(quit)
	
	<-quit
}
