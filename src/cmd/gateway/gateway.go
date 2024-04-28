package gateway

import (
	"github.com/gamedevelope/go-push/src/cli"
	"github.com/gamedevelope/go-push/src/internal/config"
	"github.com/gamedevelope/go-push/src/internal/gateway"
	"github.com/gamedevelope/go-push/src/pkg/util"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"net/http"
)

var (
	appPath string
	appMode string
)

func init() {
	serveCmd := &cobra.Command{
		Use:   "gateway",
		Short: "开启同步服务",
		Long:  "开启同步服务，同步链上事件并进行处理",
		Run:   gatewayRun,
	}

	logrus.Infof(`cli init`)
	serveCmd.Flags().StringVarP(&appPath, "app_path", `P`, ``, "application path")
	serveCmd.Flags().StringVarP(&appMode, "app_mode", `M`, ``, "application mode: local/testnet/prod")

	cli.Register(serveCmd)
}

func gatewayRun(cmd *cobra.Command, args []string) {
	cf := appPath + `/env/app.` + appMode + `.json`

	if err := config.Parse(cf); err != nil {
		logrus.Panicf(`error on parse file %v`, err)
	} else {
		logrus.Infof(`成功解析配置文件 %v`, cf)
	}

	upgrader := &websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	gServer := gateway.NewServer(&config.AppConf.GatewayConf, upgrader)

	util.Daemon(func() {
		logrus.Infof(`gateway server quit`)
		gServer.Exit()
	})
}
