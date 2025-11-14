# 固定链接配置指南

## 问题描述

默认情况下，程序绑定到 `0.0.0.0:18060`，但文档中显示的都是 `localhost:18060`。如果你想要一个固定的链接来访问运行中的程序，而不需要每次都查找局域网中的 IP 地址，可以使用以下解决方案。

## 解决方案

### 方案一：使用 localhost（推荐用于本地开发）

如果你只在本地使用，可以绑定到 `127.0.0.1`：

```bash
# 使用源码运行
go run . -host=127.0.0.1 -port=:18060

# 使用二进制文件运行
./xiaohongshu-mcp-darwin-arm64 -host=127.0.0.1 -port=:18060
```

这样你就可以使用固定链接：`http://localhost:18060/mcp`

### 方案二：使用固定端口和所有接口（推荐用于局域网访问）

如果你需要在局域网内访问，可以绑定到所有接口：

```bash
# 使用源码运行
go run . -host=0.0.0.0 -port=:18060

# 使用二进制文件运行
./xiaohongshu-mcp-darwin-arm64 -host=0.0.0.0 -port=:18060
```

然后你可以使用以下固定链接：
- 本地访问：`http://localhost:18060/mcp`
- 局域网访问：`http://[你的电脑IP]:18060/mcp`

### 方案三：使用自定义端口

如果你想要一个更容易记住的端口：

```bash
# 使用源码运行
go run . -host=0.0.0.0 -port=:8080

# 使用二进制文件运行
./xiaohongshu-mcp-darwin-arm64 -host=0.0.0.0 -port=:8080
```

这样你就可以使用：`http://localhost:8080/mcp`

### 方案四：使用环境变量

你也可以通过环境变量来设置：

```bash
# 设置环境变量
export XIAOHONGSHU_HOST=127.0.0.1
export XIAOHONGSHU_PORT=:18060

# 然后运行程序
go run .
```

## 获取你的电脑 IP 地址

如果你需要在局域网内访问，可以使用以下命令获取你的 IP 地址：

### macOS/Linux:
```bash
# 获取主要网络接口的 IP
ifconfig | grep "inet " | grep -v 127.0.0.1 | awk '{print $2}' | head -1

# 或者使用 ip 命令（Linux）
ip route get 1 | awk '{print $7}' | head -1
```

### Windows:
```cmd
ipconfig | findstr "IPv4"
```

## 配置示例

### 1. 本地开发配置

```bash
# 启动命令
go run . -host=127.0.0.1 -port=:18060

# MCP 客户端配置
{
  "mcpServers": {
    "xiaohongshu-mcp": {
      "url": "http://localhost:18060/mcp",
      "description": "小红书内容发布服务"
    }
  }
}
```

### 2. 局域网访问配置

```bash
# 启动命令（假设你的 IP 是 192.168.1.100）
go run . -host=0.0.0.0 -port=:18060

# MCP 客户端配置
{
  "mcpServers": {
    "xiaohongshu-mcp": {
      "url": "http://192.168.1.100:18060/mcp",
      "description": "小红书内容发布服务"
    }
  }
}
```

### 3. 自定义端口配置

```bash
# 启动命令
go run . -host=0.0.0.0 -port=:8080

# MCP 客户端配置
{
  "mcpServers": {
    "xiaohongshu-mcp": {
      "url": "http://localhost:8080/mcp",
      "description": "小红书内容发布服务"
    }
  }
}
```

## 自动显示局域网IP

程序启动时会自动显示局域网访问地址，无需手动查找IP：

```bash
$ go run . -host=0.0.0.0 -port=:18060
time="2025-10-26T14:28:08+08:00" level=info msg="启动 HTTP 服务器: 0.0.0.0:18060"
time="2025-10-26T14:28:08+08:00" level=info msg="MCP 服务地址: http://0.0.0.0:18060/mcp"
time="2025-10-26T14:28:08+08:00" level=info msg="局域网访问地址: http://192.168.0.113:18060/mcp"
```

这样你就可以直接复制显示的局域网地址来配置MCP客户端了！

**注意**：程序会自动处理图片URL的生成，确保文生图功能返回的图片链接使用正确的局域网IP地址，而不是错误的拼接地址。

## 验证配置

启动服务后，你可以通过以下方式验证配置是否正确：

### 1. 检查服务状态
```bash
curl http://localhost:18060/health
```

### 2. 测试 MCP 连接
```bash
curl -X POST http://localhost:18060/mcp \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","method":"initialize","params":{},"id":1}'
```

### 3. 使用 MCP Inspector
```bash
npx @modelcontextprotocol/inspector
```
然后在浏览器中输入你的 MCP 服务地址。

## 注意事项

1. **防火墙设置**：如果使用局域网访问，确保防火墙允许相应端口的访问
2. **端口冲突**：确保选择的端口没有被其他程序占用
3. **安全性**：如果绑定到 `0.0.0.0`，确保你的网络环境是安全的
4. **持久化配置**：建议将启动命令保存为脚本或使用 systemd/docker 等工具进行管理

## 推荐配置

对于大多数用户，推荐使用以下配置：

```bash
# 启动命令
go run . -host=127.0.0.1 -port=:18060

# 这样你就可以始终使用固定链接：
# http://localhost:18060/mcp
```

这个配置既简单又安全，适合本地开发和测试使用。
