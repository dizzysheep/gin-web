# gin-web

#### 介绍
gin+gorm 练习

#### 软件架构
软件架构说明


#### 安装教程

1. git clone xxxxx
2. go mod init gin-web
3. go mod vendor

#### 使用说明

1.  拉取代码
2.  初始化仓库
3.  拉取依赖库

### OpenAI Function Calling Demo

项目提供一个受 JWT 保护的后台 Agent 示例。模型通过 OpenAI Responses API 的 `tools` 字段看到 `list_articles` 工具；Go 服务收到 `function_call` 后按名称分发到已有文章查询服务，再把工具结果提交给模型生成回答。工具只读，只能查已发布文章。

Agent 会读取项目的 `config/app.toml` 中的 `[agent]` 配置：

```toml
[agent]
domain = "https://new.sharedchat.cc/codex"
key = "你的中转站 API Key"
model = "5.6-sol"
```

也可以用 `OPENAI_BASE_URL`、`OPENAI_API_KEY`、`OPENAI_MODEL` 环境变量覆盖对应配置。程序会向 `<domain>/v1/responses` 发送请求（若地址已以 `/v1` 结尾，则直接追加 `/responses`）。`config/app.toml.example` 已包含不带真实密钥的示例配置。

登录后携带项目签发的 JWT 调用：

```sh
curl -X POST http://localhost:8080/api/v1/admin/agent/chat \
  -H 'Content-Type: application/json' \
  -H 'Authorization: Bearer <JWT>' \
  -d '{"message":"列出最近的文章"}'
```

响应的 `data.answer` 是模型回答，`data.tool_calls` 展示本轮实际执行过的工具名称、参数和结果。可通过 `GET /api/v1/admin/agent/tools` 查看当前注册的工具定义（同样需要 JWT）。

#### 参与贡献

1.  Fork 本仓库
2.  新建 Feat_xxx 分支
3.  提交代码
4.  新建 Pull Request


#### 特技

1.  使用 Readme\_XXX.md 来支持不同的语言，例如 Readme\_en.md, Readme\_zh.md
2.  Gitee 官方博客 [blog.gitee.com](https://blog.gitee.com)
3.  你可以 [https://gitee.com/explore](https://gitee.com/explore) 这个地址来了解 Gitee 上的优秀开源项目
4.  [GVP](https://gitee.com/gvp) 全称是 Gitee 最有价值开源项目，是综合评定出的优秀开源项目
5.  Gitee 官方提供的使用手册 [https://gitee.com/help](https://gitee.com/help)
6.  Gitee 封面人物是一档用来展示 Gitee 会员风采的栏目 [https://gitee.com/gitee-stars/](https://gitee.com/gitee-stars/)
