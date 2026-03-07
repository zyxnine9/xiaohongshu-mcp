package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var jsonOutput bool

var RootCmd = &cobra.Command{
	Use:   "xhs",
	Short: "小红书 CLI 工具",
	Long:  `小红书 MCP 服务器兼 CLI 工具`,
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	RootCmd.PersistentFlags().BoolVarP(&jsonOutput, "json", "o", false, "输出 JSON 格式")
}
