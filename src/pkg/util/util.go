package util

import (
	"os"
	"os/signal"
	"syscall"
)

// Daemon 通用的守护进程信号处理
func Daemon(handler func()) {
	//创建监听退出chan
	sig := make(chan os.Signal, 1)
	//监听指定信号 ctrl+c kill
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	for s := range sig {
		switch s {
		case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
			handler()
			os.Exit(0)
		}
	}
}
