package gateway

import (
	"github.com/gamedevelope/go-push/src/cli"
	"github.com/gamedevelope/go-push/src/internal/config"
	"github.com/gamedevelope/go-push/src/internal/gateway"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
	"time"
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

	var (
		err error
	)

	gateway.InitConfig(&config.AppConf.GatewayConf)

	// 统计
	if err = gateway.InitStats(); err != nil {
		logrus.Panicf(`init stats %v`, err)
	}

	// 初始化连接管理器
	if err = gateway.InitConnMgr(); err != nil {
		logrus.Panicf(`init conn mgr %v`, err)
	}

	// 初始化websocket服务器
	if err = gateway.InitWSServer(); err != nil {
		logrus.Panicf(`init ws server %v`, err)
	}

	// 初始化merger合并层
	if err = gateway.InitMerger(); err != nil {
		logrus.Panicf(`init merger %v`, err)
	}

	// 初始化service接口
	if err = gateway.InitService(); err != nil {
		logrus.Panicf(`init service %v`, err)
	}

	for {
		time.Sleep(1 * time.Second)
	}

	os.Exit(0)
}
