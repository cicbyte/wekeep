// Package storage 存储管理控制器
package storage

import (
	"context"
	"fmt"

	api "github.com/ciclebyte/wekeep/api/v1/storage"
	storageLogic "github.com/ciclebyte/wekeep/internal/storage"
	"github.com/gogf/gf/v2/frame/g"
)

// Storage 存储控制器实例
var Storage = controller{}

// controller 存储控制器
type controller struct{}

// Migrate 执行存储迁移
func (c *controller) Migrate(ctx context.Context, req *api.StorageMigrateReq) (res *api.StorageMigrateRes, err error) {
	res = &api.StorageMigrateRes{}

	if !storageLogic.IsMigrationEnabled() {
		return nil, fmt.Errorf("存储迁移功能未启用，请在配置文件中设置 storage.migration.enabled 为 true")
	}

	// 检查是否已有迁移任务在运行
	status := storageLogic.GetMigrationStatus()
	if status.Running {
		return nil, fmt.Errorf("迁移任务正在进行中，请等待完成")
	}

	// 转换配置
	migrationCfg := &storageLogic.MigrationConfig{
		SourceStorage:  req.SourceStorage,
		TargetStorage:  req.TargetStorage,
		UpdateMarkdown: req.UpdateMarkdown,
	}

	if req.LocalConfig != nil {
		migrationCfg.LocalConfig = &storageLogic.LocalConfig{
			BasePath: req.LocalConfig.BasePath,
			BaseURL:  req.LocalConfig.BaseURL,
		}
	}

	if req.RustFSConfig != nil {
		migrationCfg.RustFSConfig = &storageLogic.RustFSConfig{
			Endpoint: req.RustFSConfig.Endpoint,
			Bucket:   req.RustFSConfig.Bucket,
			Username: req.RustFSConfig.Username,
			Password: req.RustFSConfig.Password,
			Timeout:  req.RustFSConfig.Timeout,
		}
	}

	// 执行迁移
	err = storageLogic.MigrateStorage(ctx, migrationCfg)
	if err != nil {
		return nil, err
	}

	res.Message = "迁移任务已启动，请通过 /storage/status 接口查看进度"
	return res, nil
}

// Status 获取存储状态
func (c *controller) Status(ctx context.Context, req *api.StorageStatusReq) (res *api.StorageStatusRes, err error) {
	res = &api.StorageStatusRes{
		CurrentStorage: storageLogic.GetCurrentStorageType(),
		Initialized:    storageLogic.IsInitialized(),
		MigrationEnabled: storageLogic.IsMigrationEnabled(),
	}

	// 获取迁移状态
	migrationStatus := storageLogic.GetMigrationStatus()
	res.Migration = &api.MigrationStatus{
		Running:     migrationStatus.Running,
		Source:      migrationStatus.Source,
		Target:      migrationStatus.Target,
		Total:       migrationStatus.Total,
		Completed:   migrationStatus.Completed,
		Failed:      migrationStatus.Failed,
		CurrentFile: migrationStatus.CurrentFile,
		Error:       migrationStatus.Error,
	}

	if migrationStatus.StartTime != nil {
		res.Migration.StartTime = migrationStatus.StartTime.Format("Y-m-d H:i:s")
	}
	if migrationStatus.EndTime != nil {
		res.Migration.EndTime = migrationStatus.EndTime.Format("Y-m-d H:i:s")
	}

	// 获取当前配置（隐藏敏感信息）
	currentConfig := storageLogic.GetCurrentConfig()
	if currentConfig != nil {
		res.CurrentConfig = &api.StorageConfig{}
		if currentConfig.RustFS != nil {
			res.CurrentConfig.RustFS = &api.RustFSConfigDisplay{
				Endpoint: currentConfig.RustFS.Endpoint,
				Bucket:   currentConfig.RustFS.Bucket,
				Username: currentConfig.RustFS.Username,
				Timeout:  currentConfig.RustFS.Timeout,
			}
		}
		if currentConfig.Local != nil {
			res.CurrentConfig.Local = &api.LocalConfigDisplay{
				BasePath: currentConfig.Local.BasePath,
				BaseURL:  currentConfig.Local.BaseURL,
			}
		}
	}

	return res, nil
}

// Validate 验证存储配置
func (c *controller) Validate(ctx context.Context, req *api.StorageValidateReq) (res *api.StorageValidateRes, err error) {
	res = &api.StorageValidateRes{Valid: true}

	if !storageLogic.IsMigrationEnabled() {
		return nil, fmt.Errorf("存储迁移功能未启用")
	}

	// 转换配置
	migrationCfg := &storageLogic.MigrationConfig{
		TargetStorage: req.TargetStorage,
	}

	if req.LocalConfig != nil {
		migrationCfg.LocalConfig = &storageLogic.LocalConfig{
			BasePath: req.LocalConfig.BasePath,
			BaseURL:  req.LocalConfig.BaseURL,
		}
	}

	if req.RustFSConfig != nil {
		migrationCfg.RustFSConfig = &storageLogic.RustFSConfig{
			Endpoint: req.RustFSConfig.Endpoint,
			Bucket:   req.RustFSConfig.Bucket,
			Username: req.RustFSConfig.Username,
			Password: req.RustFSConfig.Password,
			Timeout:  req.RustFSConfig.Timeout,
		}
	}

	// 验证配置
	validateErr := storageLogic.ValidateStorageConfig(ctx, migrationCfg)
	if validateErr != nil {
		res.Valid = false
		res.Error = validateErr.Error()
	}

	return res, nil
}

// Stats 获取存储统计
func (c *controller) Stats(ctx context.Context, req *api.StorageStatsReq) (res *api.StorageStatsRes, err error) {
	res = &api.StorageStatsRes{
		ByStorageType: make(map[string]int64),
	}

	// 统计图片数量
	totalImages, err := g.DB().Model("images").Where("storage_path != ?", "").Count()
	if err != nil {
		return nil, err
	}
	res.TotalImages = totalImages

	// 统计各存储类型的大小（根据 URL 前缀判断）
	// 本地存储：/uploads 开头
	// RustFS：根据配置的 endpoint 判断

	var localCount int = 0

	// 统计本地存储
	localCount, err = g.DB().Model("images").
		Where("storage_url LIKE ?", "/uploads/%").
		Count()
	if err == nil {
		res.ByStorageType["local"] = int64(localCount)
	}

	// 统计 RustFS 存储（非本地存储的）
	rustfsCount := int64(totalImages) - int64(localCount)
	res.ByStorageType["rustfs"] = rustfsCount

	return res, nil
}

// Switch 切换存储（更新配置）
func (c *controller) Switch(ctx context.Context, req *api.StorageSwitchReq) (res *api.StorageSwitchRes, err error) {
	res = &api.StorageSwitchRes{}

	// 检查是否已完成迁移
	migrationStatus := storageLogic.GetMigrationStatus()
	if migrationStatus.Running {
		return nil, fmt.Errorf("迁移任务正在进行中，请等待完成后再切换")
	}

	// 这里只更新运行时配置，实际生产环境需要更新配置文件
	// 由于安全原因，不直接修改配置文件，而是提示用户手动更新

	res.Message = "请在配置文件中更新存储配置后重启服务"
	res.Requires = "需要手动修改 config.yaml 中的 storage.type 和相关配置，然后重启服务"

	return res, nil
}

// List 存储信息列表
func (c *controller) List(ctx context.Context, req *api.ListReq) (res *api.ListRes, err error) {
	res = &api.ListRes{}
	// 此接口预留，可用于返回存储使用详情
	return res, nil
}

// UpdateRefs 更新图片引用
func (c *controller) UpdateRefs(ctx context.Context, req *api.StorageUpdateRefsReq) (res *api.StorageUpdateRefsRes, err error) {
	res = &api.StorageUpdateRefsRes{}

	if !storageLogic.IsMigrationEnabled() {
		return nil, fmt.Errorf("存储迁移功能未启用")
	}

	// 构建图片 URL 映射（从数据库读取）
	urlMapping, err := storageLogic.BuildImageURLMapping(ctx)
	if err != nil {
		return nil, fmt.Errorf("构建图片映射失败: %w", err)
	}

	if len(urlMapping) == 0 {
		res.Message = "没有需要更新的图片引用"
		res.Updated = 0
		res.Total = 0
		return res, nil
	}

	// 更新 Markdown 内容
	updated, err := storageLogic.UpdateMarkdownContent(ctx, urlMapping)
	if err != nil {
		return nil, fmt.Errorf("更新引用失败: %w", err)
	}

	res.Message = "图片引用更新完成"
	res.Updated = updated
	res.Total = len(urlMapping)
	return res, nil
}
