package main

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/xpzouying/xiaohongshu-mcp/pkg/output"
)

var (
	keyword string
	limit   int
)

var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "搜索笔记",
	Example: `  xhs search -k "美食"
  xhs search -k "旅游" -l 20`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		service := NewXiaohongshuService()

		result, err := service.SearchFeeds(ctx, keyword)
		if err != nil {
			output.NewOutput(jsonOutput).Error(err)
			return
		}

		if limit > 0 && len(result.Feeds) > limit {
			result.Feeds = result.Feeds[:limit]
		}

		output.NewOutput(jsonOutput).Success("搜索完成")
		output.NewOutput(jsonOutput).Print(result)
	},
}

func init() {
	searchCmd.Flags().StringVarP(&keyword, "keyword", "k", "", "搜索关键词 (必填)")
	searchCmd.Flags().IntVarP(&limit, "limit", "l", 0, "返回结果数量限制")

	rootCmd.AddCommand(searchCmd)
}
