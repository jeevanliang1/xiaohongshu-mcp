# Panic 恢复机制说明

## 概述

为了防止应用在遇到panic时崩溃，我们添加了全面的panic恢复机制。现在当发生panic时，应用会记录错误信息并继续运行，而不是直接退出。

## 实现的防御机制

### 1. 中间件级别的恢复

在 `middleware.go` 中添加了两个通用的panic恢复包装器：

- `panicRecoveryWrapper(fn func() error) error` - 用于包装返回error的函数
- `panicRecoveryWrapperWithResult[T any](fn func() (T, error)) (T, error)` - 用于包装返回结果和error的函数

### 2. MCP工具级别的恢复

在 `mcp_server.go` 中为所有MCP工具添加了panic恢复机制：

- 检查登录状态
- 获取登录二维码  
- 发布内容
- 搜索内容
- 获取Feed详情
- 用户主页
- 发表评论
- 发布视频
- 点赞/收藏功能

当工具执行过程中发生panic时，会返回友好的错误信息而不是让整个应用崩溃。

### 3. 服务级别的恢复

在 `service.go` 中为关键服务函数添加了panic恢复：

- `SearchFeeds` - 搜索功能
- `PublishContent` - 发布内容功能

### 4. HTTP服务器级别的恢复

在 `app_server.go` 中为HTTP服务器的goroutine添加了panic恢复：

- 当HTTP服务器goroutine发生panic时，会记录错误并尝试重启服务器
- 不会让整个应用退出

## 错误处理流程

1. **捕获panic**: 使用 `defer recover()` 捕获panic
2. **记录错误**: 使用logrus记录详细的错误信息
3. **返回友好错误**: 将panic转换为用户友好的错误信息
4. **继续运行**: 应用继续运行，不会崩溃

## 示例

### 之前的行为
```go
// 当发生panic时，整个应用会崩溃
panic("circular JSON structure")
// 应用退出，无法继续使用
```

### 现在的行为
```go
// 当发生panic时，会被捕获并转换为错误
err := panicRecoveryWrapper(func() error {
    panic("circular JSON structure")
})
// err 包含友好的错误信息，应用继续运行
```

## 日志输出

当发生panic时，会在日志中看到类似这样的信息：

```
ERROR[2024-01-01T12:00:00Z] 捕获到panic: circular JSON structure
ERROR[2024-01-01T12:00:00Z] 搜索内容失败: 操作失败，发生内部错误: circular JSON structure
```

## 优势

1. **提高稳定性**: 应用不会因为单个操作失败而崩溃
2. **更好的用户体验**: 用户会收到友好的错误信息而不是应用崩溃
3. **便于调试**: 详细的错误日志帮助定位问题
4. **自动恢复**: 某些情况下会自动重试或恢复

## 注意事项

- panic恢复机制会捕获所有panic，包括程序逻辑错误
- 建议在开发过程中仍然要修复导致panic的根本原因
- 恢复机制主要用于处理外部依赖（如浏览器、网络）的不可预测错误
