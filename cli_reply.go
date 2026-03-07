package main

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/xpzouying/xiaohongshu-mcp/pkg/output"
)

var (
	commentID   string
	replyUserID string
)

var replyCmd = &cobra.Command{
	Use:     "reply",
	Short:   "回复评论",
	Example: `  xhs reply --feed-id "xxx" --comment-id "xxx" --user-id "xxx" -c "回复内容"`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		service := NewXiaohongshuService()

		result, err := service.ReplyCommentToFeed(ctx, feedID, xsecToken, commentID, replyUserID, commentContent)
		if err != nil {
			output.NewOutput(jsonOutput).Error(err)
			return
		}

		output.NewOutput(jsonOutput).Success("回复成功")
		output.NewOutput(jsonOutput).Print(result)
	},
}

func init() {
	replyCmd.Flags().StringVar(&feedID, "feed-id", "", "笔记 ID (必填)")
	replyCmd.Flags().StringVar(&xsecToken, "token", "", "xsec token")
	replyCmd.Flags().StringVar(&commentID, "comment-id", "", "评论 ID (必填)")
	replyCmd.Flags().StringVar(&replyUserID, "user-id", "", "被回复用户 ID (必填)")
	replyCmd.Flags().StringVarP(&commentContent, "content", "c", "", "回复内容 (必填)")

	rootCmd.AddCommand(replyCmd)
}
