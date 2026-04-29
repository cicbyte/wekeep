package controller

import (
	"context"

	api "github.com/cicbyte/wekeep/api/v1/health"
	service "github.com/cicbyte/wekeep/internal/service"
)

var Health = healthController{}

type healthController struct {
	BaseController
}

// Check 简单健康检查
func (c *healthController) Check(ctx context.Context, req *api.HealthReq) (res *api.HealthRes, err error) {
	res = &api.HealthRes{}
	res.Status, res.Message = service.Health().Check(ctx)
	return
}

// Detail 详细健康检查
func (c *healthController) Detail(ctx context.Context, req *api.HealthDetailReq) (res *api.HealthDetailRes, err error) {
	return service.Health().Detail(ctx)
}

// Version 版本信息
func (c *healthController) Version(ctx context.Context, req *api.VersionReq) (res *api.VersionRes, err error) {
	return service.Health().Version(ctx)
}