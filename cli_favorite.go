package main

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/xpzouying/xiaohongshu-mcp/pkg/output"
)

var unfavorite bool

var favoriteCmd = &cobra.Command{
	Use:   "favorite",
	Short: "收藏/取消收藏笔记",
	Example: `  xhs favorite --feed-id "xxx"
  xhs favorite --feed-id "xxx" --unfavorite`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		service := NewXiaohongshuService()

		var result interface{}
		var err error

		if unfavorite {
			result, err = service.UnfavoriteFeed(ctx, feedID, xsecToken)
		} else {
			result, err = service.FavoriteFeed(ctx, feedID, xsecToken)
		}

		if err != nil {
			output.NewOutput(jsonOutput).Error(err)
			return
		}

		if unfavorite {
			output.NewOutput(jsonOutput).Success("取消收藏成功")
		} else {
			output.NewOutput(jsonOutput).Success("收藏成功")
		}
		output.NewOutput(jsonOutput).Print(result)
	},
}

func init() {
	favoriteCmd.Flags().StringVar(&feedID, "feed-id", "", "笔记 ID (必填)")
	favoriteCmd.Flags().StringVar(&xsecToken, "token", "", "xsec token")
	favoriteCmd.Flags().BoolVar(&unfavorite, "unfavorite", false, "取消收藏")

	rootCmd.AddCommand(favoriteCmd)
}
