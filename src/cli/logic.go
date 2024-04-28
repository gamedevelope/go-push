package cli

import (
	"fmt"
	logic2 "github.com/gamedevelope/go-push/src/logic"
	"github.com/spf13/cobra"
	"log/slog"
	"os"
	"time"
)

func registerLogic() {
	serveCmd := &cobra.Command{
		Use:   "logic",
		Short: "开启同步服务",
		Long:  "开启同步服务，同步链上事件并进行处理",
		Run:   logicRun,
	}

	rootCmd.AddCommand(serveCmd)
}

func logicRun(cmd *cobra.Command, args []string) {
	slog.Info(`logic run`)
	var (
		err error
	)

	if err = logic2.InitConfig(confFile); err != nil {
		goto ERR
	}

	if err = logic2.InitStats(); err != nil {
		goto ERR
	}

	if err = logic2.InitGateConnMgr(); err != nil {
		goto ERR
	}

	if err = logic2.InitService(); err != nil {
		goto ERR
	}

	for {
		time.Sleep(1 * time.Second)
	}

	os.Exit(0)

ERR:
	fmt.Fprintln(os.Stderr, err)
	os.Exit(-1)
}
