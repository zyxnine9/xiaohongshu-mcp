# CLI 改造设计方案

## 概述

将现有的小红书 MCP 服务器改造为同时支持 CLI 工具和 MCP 服务器。

## 架构设计

```
┌─────────────────────────────────────────┐
│              main.go                    │
│  - 无参数运行: 启动 MCP 服务器          │
│  - 使用子命令: 执行 CLI 操作            │
└─────────────────────────────────────────┘
         ↓                              ↓
   ┌───────────┐                  ┌───────────┐
   │ CLI Mode  │                  │ MCP Mode  │
   │ (cobra)   │                  │ (原有逻辑) │
   └───────────┘                  └───────────┘
```

## 子命令设计

| 子命令 | 功能 |
|--------|------|
| `xhs publish` | 发布图文 |
| `xhs publish-video` | 发布视频 |
| `xhs search` | 搜索笔记 |
| `xhs get-feed` | 获取笔记详情 |
| `xhs get-user` | 获取用户主页 |
| `xhs comment` | 发表评论 |
| `xhs reply` | 回复评论 |
| `xhs like` | 点赞/取消点赞 |
| `xhs favorite` | 收藏/取消收藏 |
| `xhs login` | 获取登录二维码 |
| `xhs status` | 检查登录状态 |
| `xhs server` | 启动 MCP 服务器 |

## 输出格式

- **默认**: 人类可读的美化输出（带颜色、格式化）
- **JSON**: 添加 `--json` 或 `-o json` 标志输出 JSON

## 依赖库

使用 `github.com/spf13/cobra` 作为 CLI 框架。
