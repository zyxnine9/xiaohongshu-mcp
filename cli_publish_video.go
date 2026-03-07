package main

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/xpzouying/xiaohongshu-mcp/pkg/output"
)

var (
	videoPath string
)

var publishVideoCmd = &cobra.Command{
	Use:   "publish-video",
	Short: "发布视频笔记",
	Example: `  xhs publish-video -t "标题" -c "正文" -v video.mp4
  xhs publish-video -t "标题" -c "正文" -v video.mp4 --json`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		service := NewXiaohongshuService()

		req := &PublishVideoRequest{
			Title:      title,
			Content:    content,
			Video:      videoPath,
			Tags:       tags,
			ScheduleAt: schedule,
		}

		resp, err := service.PublishVideo(ctx, req)
		if err != nil {
			output.NewOutput(jsonOutput).Error(err)
			return
		}

		output.NewOutput(jsonOutput).Success("发布成功")
		output.NewOutput(jsonOutput).Print(resp)
	},
}

func init() {
	publishVideoCmd.Flags().StringVarP(&title, "title", "t", "", "笔记标题 (必填)")
	publishVideoCmd.Flags().StringVarP(&content, "content", "c", "", "笔记正文 (必填)")
	publishVideoCmd.Flags().StringVarP(&videoPath, "video", "v", "", "视频文件路径 (必填)")
	publishVideoCmd.Flags().StringSliceVarP(&tags, "tag", "", []string{}, "标签，支持多个")
	publishVideoCmd.Flags().StringVarP(&schedule, "schedule", "s", "", "定时发布时间 (ISO8601格式)")

	rootCmd.AddCommand(publishVideoCmd)
}
