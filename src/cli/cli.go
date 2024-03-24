package cli

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	appPath string
	appMode string
)

var (
	confFile = ``

	rootCmd = &cobra.Command{
		Use:   "push",
		Short: ``,
		Long:  ``,
	}
)

func init() {

}

func Register(cmd *cobra.Command) {
	rootCmd.AddCommand(cmd)
}

func Exec() {
	//registerGateway()
	registerLogic()

	if err := rootCmd.Execute(); err != nil {
		logrus.Panic(err)
	} else {
		logrus.Infof(`cli exec`)
	}
}
