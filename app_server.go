package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/sirupsen/logrus"
)

// AppServer 应用服务器结构体，封装所有服务和处理器
type AppServer struct {
	xiaohongshuService *XiaohongshuService
	mcpServer          *mcp.Server
	router             *gin.Engine
	httpServer         *http.Server
}

// NewAppServer 创建新的应用服务器实例
func NewAppServer(xiaohongshuService *XiaohongshuService) *AppServer {
	appServer := &AppServer{
		xiaohongshuService: xiaohongshuService,
	}

	// 初始化 MCP Server（需要在创建 appServer 之后，因为工具注册需要访问 appServer）
	appServer.mcpServer = InitMCPServer(appServer)

	return appServer
}

// Start 启动服务器
func (s *AppServer) Start(port string) error {
	s.router = setupRoutes(s)

	s.httpServer = &http.Server{
		Addr:    port,
		Handler: s.router,
	}

	// 启动服务器的 goroutine
	go func() {
		defer func() {
			if r := recover(); r != nil {
				logrus.Errorf("HTTP服务器goroutine发生panic: %v", r)
				// 不退出程序，而是尝试重启服务器
				logrus.Info("尝试重启HTTP服务器...")
				time.Sleep(5 * time.Second)
				go func() {
					if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
						logrus.Errorf("服务器重启失败: %v", err)
					}
				}()
			}
		}()

		logrus.Infof("启动 HTTP 服务器: %s", port)
		logrus.Infof("MCP 服务地址: http://%s/mcp", port)

		// 显示局域网访问地址
		if port == "0.0.0.0:18060" || (len(port) > 0 && port[0] == ':') {
			// 如果绑定到所有接口或端口以冒号开头，显示局域网访问地址
			localIP := getLocalIP()
			if port[0] == ':' {
				logrus.Infof("局域网访问地址: http://%s%s/mcp", localIP, port)
			} else {
				// 提取端口号
				portOnly := port[strings.LastIndex(port, ":")+1:]
				logrus.Infof("局域网访问地址: http://%s:%s/mcp", localIP, portOnly)
			}
		}
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logrus.Errorf("服务器启动失败: %v", err)
			os.Exit(1)
		}
	}()

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logrus.Infof("正在关闭服务器...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.httpServer.Shutdown(ctx); err != nil {
		logrus.Warnf("等待连接关闭超时，强制退出: %v", err)
	} else {
		logrus.Infof("服务器已优雅关闭")
	}

	return nil
}
