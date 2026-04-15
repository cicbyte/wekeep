package health

import (
	"github.com/gogf/gf/v2/frame/g"
)

// HealthReq 健康检查请求
type HealthReq struct {
	g.Meta `path:"/health" method:"get" tags:"系统" summary:"健康检查"`
}

// HealthRes 健康检查响应
type HealthRes struct {
	g.Meta `mime:"application/json"`
	Status  string `json:"status"`  // 状态: ok, error
	Message string `json:"message"` // 消息
}

// HealthDetailReq 详细健康检查请求
type HealthDetailReq struct {
	g.Meta `path:"/health/detail" method:"get" tags:"系统" summary:"详细健康检查"`
}

// VersionReq 版本请求
type VersionReq struct {
	g.Meta `path:"/health/version" method:"get" tags:"系统" summary:"版本信息"`
}

// VersionRes 版本响应
type VersionRes struct {
	g.Meta `mime:"application/json"`
	Version string `json:"version"` // 版本号
}

// CheckItem 检查项
type CheckItem struct {
	Name   string `json:"name"`   // 检查项名称
	Status string `json:"status"` // 状态: ok, error
	Error  string `json:"error,omitempty"` // 错误信息
}

// HealthDetailRes 详细健康检查响应
type HealthDetailRes struct {
	g.Meta   `mime:"application/json"`
	Status   string      `json:"status"`   // 总体状态: ok, error
	Message  string      `json:"message"`  // 消息
	Checks   []CheckItem `json:"checks"`   // 各项检查结果
	Uptime   string      `json:"uptime"`   // 运行时长
	Version  string      `json:"version"`  // 版本信息
}
