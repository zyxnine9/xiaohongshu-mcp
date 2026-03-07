package main

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/xpzouying/xiaohongshu-mcp/configs"
)

var (
	serverPort      string
	serverHeadless  bool
	serverBin       string
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "启动 MCP 服务器",
	Example: `  xhs server
  xhs server --port 8080`,
	Run: func(cmd *cobra.Command, args []string) {
		configs.InitHeadless(serverHeadless)
		configs.SetBinPath(serverBin)

		xiaohongshuService := NewXiaohongshuService()
		appServer := NewAppServer(xiaohongshuService)
		if err := appServer.Start(serverPort); err != nil {
			logrus.Fatalf("failed to run server: %v", err)
		}
	},
}

func init() {
	serverCmd.Flags().StringVar(&serverPort, "port", ":18060", "服务端口")
	serverCmd.Flags().BoolVar(&serverHeadless, "headless", true, "是否无头模式")
	serverCmd.Flags().StringVar(&serverBin, "bin", "", "浏览器二进制路径")

	rootCmd.AddCommand(serverCmd)
}

// 避免编译错误 - 确保 serverHeadless 等变量在使用前被正确初始化
func init() {
	// 从环境变量读取默认值
	if serverBin == "" {
		serverBin = os.Getenv("ROD_BROWSER_BIN")
	}
}
