// +build windows

package app

import (
	"os"
	"os/signal"
)

func (a *GcoreApp) loop() {
	signal.Notify(a.signal, os.Interrupt)

	for {
		sig := <-a.signal
		logger.Infof("接收到信号: %s", sig)

		switch sig {
		case os.Interrupt:
			return
		}
	}
}
