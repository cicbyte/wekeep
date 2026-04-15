package mcp

import (
	"github.com/mark3labs/mcp-go/server"
)

// NewMCPServer 创建并配置 MCP 服务器
// 该函数初始化一个新的 MCP 服务器实例，注册所有微信公众号相关工具
func NewMCPServer() *server.MCPServer {
	// 创建 MCP 服务器实例
	s := server.NewMCPServer(
		"WeKeep MCP Server",    // 服务器名称
		"1.0.0",                 // 版本号
		server.WithToolCapabilities(true), // 启用工具能力
	)

	// 注册微信公众号相关工具
	registerTools(s)

	return s
}

// NewStreamableHTTPServer 创建支持 HTTP 的 MCP 服务器
// 使用 StreamableHTTP 协议，支持通过 HTTP 请求调用 MCP 工具
func NewStreamableHTTPServer() *server.StreamableHTTPServer {
	mcpServer := NewMCPServer()
	return server.NewStreamableHTTPServer(mcpServer)
}
