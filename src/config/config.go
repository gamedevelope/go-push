package config

import (
	"encoding/json"
	"github.com/gamedevelope/go-push/src/gateway"
	"github.com/gamedevelope/go-push/src/logic"
	"os"
)

type appConfig struct {
	GatewayConf gateway.Config `json:"gateway_conf"`
	LogicConf   logic.Config   `json:"logic_conf"`
}

var (
	AppConf appConfig
)

func Parse(path string) error {
	file, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	if err = json.Unmarshal(file, &AppConf); err != nil {
		return err
	}

	return nil
}
