# 代码生成（Code Generation）

## 概述

要的文件：api/openapi.yaml（入口） + paths/ responses/ schemas（内容） + ports/marketplaceapi/oapi-config.yaml（配置） + ports/marketplaceapi/generate.go（命令钩子）。
要的命令：go mod init（一次）→ go get（一次）→ go generate ./ports/marketplaceapi（每次改规范后）。
流程：写 spec → 配置生成 → 跑生成 → 实现接口 → 起服务测通。

坑点 oapi-codegen v2， oapi-config等要遵守V2规范

## 生成所需文件、命令与流程（简要）

要的文件：
- api/openapi.yaml（入口）
- api/paths/、api/responses/、api/schemas/（规范拆分的内容）
- ports/marketplaceapi/oapi-config.yaml（oapi-codegen 配置）
- ports/marketplaceapi/generate.go（命令钩子，包含 //go:generate 注释）

要的命令（简要步骤）：
1. go mod init （仅首次）
2. go get / go install 指定生成工具（仅首次或升级时）
3. 每次改 spec 后运行： 
    ```shell 
    go generate ./ports/marketplaceapi
    ```
典型流程（工程化建议）：
1. 写 OpenAPI spec（api/openapi.yaml + paths/ responses/ schemas/）
2. 配置生成器（ports/marketplaceapi/oapi-config.yaml）(注意用新版本)
3. 在 ports/marketplaceapi/generate.go 或源文件中添加 //go:generate，并运行 go generate
4. 实现生成的接口（实现 ServerInterface 等）
5. 启动服务并通过集成/端到端测试验证

说明：把生成步骤封装到 Makefile 或 generate.sh，有利于团队复现；关键生成产物（如 pb.go、marketplace.gen.go 等）可以考虑纳入版本控制或在 CI 中强制生成并校验差异。
