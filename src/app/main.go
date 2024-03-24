package main

import (
	"fmt"
	"github.com/gamedevelope/go-push/src/cli"
	_ "github.com/gamedevelope/go-push/src/cmd/gateway"
	"github.com/sirupsen/logrus"
	"runtime"
	"time"
)

var (
	buildPath string
)

func init() {
	buildPathLen := len(buildPath)

	logrus.SetReportCaller(true)
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:    true,
		TimestampFormat:  time.DateTime,
		DisableTimestamp: false,
		FieldMap:         nil,
		CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
			return ``, fmt.Sprintf(`%v:%v`, frame.File[buildPathLen:], frame.Line)
		},
	})
}

func main() {
	logrus.Infof(`main run`)
	cli.Exec()
}
