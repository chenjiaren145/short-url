# 项目启动与运行支持计划 (修订版)

根据您的反馈，我将更新计划，不仅包括启动项目，还包括详细解释每个命令的含义以及如何使用 `curl` 验证请求。

## 目标

1. 确保项目可以在本地顺利启动，并支持 Docker 环境。
2. 提供详尽的命令解释文档，包括启动命令和验证命令。

## 执行步骤

### 1. 创建 Dockerfile

* 在项目根目录创建 `Dockerfile`，用于构建应用的 Docker 镜像。

* 使用多阶段构建：

  * **构建阶段**: 使用 `golang:alpine` 编译 `cmd/short-url/main.go`。

  * **运行阶段**: 使用轻量级 `alpine` 镜像，复制构建好的二进制文件。

### 2. 更新 docker-compose.yml

* 修改现有的 `docker-compose.yml`，添加 `app` 服务。

* 配置 `app` 服务依赖 `redis`，并通过命令参数连接 Redis：`-store redis -redis-addr redis:6379`。

* 映射端口 `8081:8081`。

### 3. 更新文档 (README.md)

* 添加 **"如何启动"** 章节，并对每个命令进行解释：

  * **方式一：本地内存模式** (无需依赖，适合快速开发):

    ```bash
    go run cmd/short-url/main.go
    ```

    * **解释**: 直接编译并运行 Go 程序，使用内存作为临时存储。

  * **方式二：Docker 完整模式** (适合生产环境模拟):

    ```bash
    docker-compose up --build
    ```

    * **解释**: 构建并启动所有定义在 `docker-compose.yml` 中的服务（App + Redis）。

### 4. 添加 API 验证指南

* 在 README 中添加 **"API 验证"** 章节，包含具体的 `curl` 命令及其含义：

  * **创建短链接**:

    ```bash
    curl -X POST -H "Content-Type: application/json" -d '{"original_url": "https://www.google.com"}' http://localhost:8081/shorten
    ```

    * **解释**:

      * `-X POST`: 指定 HTTP 方法为 POST。

      * `-H "Content-Type: ..."`: 设置请求头，表明发送的是 JSON 数据。

      * `-d '{...}'`: 发送的 JSON 数据体，包含原始 URL。

  * **访问短链接**:

    ```bash
    curl -I http://localhost:8081/<short_code>
    ```

    * **解释**:

      * `-I` (或 `--head`): 只获取 HTTP 响应头，不下载内容。用于检查重定向状态码 (302) 和 `Location` 头。

### 5. 验证

* 构建 Docker 镜像并启动。

* 执行 `curl` 命令验证 API 功能。

