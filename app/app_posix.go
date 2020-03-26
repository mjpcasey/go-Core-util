// +build !windows

package app

import (
	"os/signal"
	"syscall"
)

func (a *GcoreApp) loop() {
	signal.Notify(a.signal, syscall.SIGUSR1, syscall.SIGUSR2, syscall.SIGINT, syscall.SIGTERM) // 监听信号

	for {
		sig := <-a.signal
		logger.Infof("接收到信号: %s", sig)

		switch sig {
		case syscall.SIGUSR1:
			a.share.Reload()
		case syscall.SIGUSR2:

		default:
			return
		}
	}
}
