package gateway

import (
	"encoding/json"
	"fmt"
	"github.com/gamedevelope/go-push/src/cli"
	"github.com/gamedevelope/go-push/src/internal/config"
	"github.com/gamedevelope/go-push/src/internal/gateway"
	"github.com/gamedevelope/go-push/src/pkg/util"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"net/http"
	"strconv"
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

type AuthDebug struct {
}

func (a AuthDebug) GetUid(uri, token string) (uint64, error) {
	// 将token转换为uid
	return strconv.ParseUint(token, 10, 64)
}

type Auth struct {
}

type UserInfoResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    uint64 `json:"data"`
}

func (a Auth) GetUid(uri, token string) (uint64, error) {
	// 从远程接口中获取认证信息
	logrus.Infof(`从远程接口中获取认证信息 %v`, uri)
	client := &http.Client{}

	// get user info from remote interface with header
	req, _ := http.NewRequest("GET", uri, nil)
	req.Header.Set("Authorization", `Bearer `+token)

	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf(`%v`, resp.StatusCode)
	}

	// parse user info
	var userResp UserInfoResponse
	if err = json.NewDecoder(resp.Body).Decode(&userResp); err != nil {
		return 0, err
	}

	return userResp.Data, nil
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

	gServer := gateway.NewServer(&config.AppConf.GatewayConf, upgrader, AuthDebug{})

	util.Daemon(func() {
		logrus.Infof(`gateway server quit`)
		gServer.Exit()
	})
}
