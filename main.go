package main

import (
	"flag"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/xpzouying/xiaohongshu-mcp/configs"
)

func main() {
	// 检查是否有子命令参数，有则运行 CLI
	if len(os.Args) > 1 {
		Execute()
		return
	}

	// 无参数，运行 MCP 服务器
	runServer()
}

var (
	headless bool
	binPath  string
	port     string
)

func init() {
	flag.BoolVar(&headless, "headless", true, "是否无头模式")
	flag.StringVar(&binPath, "bin", "", "浏览器二进制文件路径")
	flag.StringVar(&port, "port", ":18060", "端口")
}

func runServer() {
	flag.Parse()

	if len(binPath) == 0 {
		binPath = os.Getenv("ROD_BROWSER_BIN")
	}

	configs.InitHeadless(headless)
	configs.SetBinPath(binPath)

	xiaohongshuService := NewXiaohongshuService()
	appServer := NewAppServer(xiaohongshuService)
	if err := appServer.Start(port); err != nil {
		logrus.Fatalf("failed to run server: %v", err)
	}
}
