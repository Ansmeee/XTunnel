package main

import (
	"os"
	"os/signal"
	"syscall"
	"xtunnel/views"
)

func main() {
	views.NewWindow().Run()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}
