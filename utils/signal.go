package utils

import (
	"os"
	"os/signal"
	"syscall"
)

func SetupSignalHandler() (stopCh <-chan struct{}) {
	stop := make(chan struct{})
	c := make(chan os.Signal, 2)
	signal.Notify(
		c,
		syscall.SIGINT,
		syscall.SIGTERM,
		os.Interrupt,
	)

	go func() {
		<-c
		close(stop)
		<-c
		os.Exit(1)
	}()

	return stop
}
