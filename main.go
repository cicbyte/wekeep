package main

import (
	"os"
	"path/filepath"

	_ "github.com/ciclebyte/wekeep/internal/packed"
	//重要 需要导入数据库驱动
	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"
	// SQLite数据库驱动,如果需要支持需要go get -u github.com/gogf/gf/contrib/drivers/sqlite/v2
	_ "github.com/gogf/gf/contrib/drivers/sqlite/v2"

	"github.com/gogf/gf/v2/os/gctx"

	"github.com/ciclebyte/wekeep/internal/cmd"
)

// 默认配置（首次运行时自动创建，用户可修改）
const defaultConfig = `# WeKeep 配置文件（可修改）
server:
  address: ":8000"
  logPath: "log"
  logStdout: true
  errorStack: true
  errorLogEnabled: true
  errorLogPattern: "error-{Ymd}.log"
  accessLogEnabled: true
  accessLogPattern: "access-{Ymd}.log"

logger:
  path: "log"
  file: "{Y-m-d}.log"
  level: "all"
  stdout: true

database:
  default:
    link: "sqlite::@file(wekeep.db)"

storage:
  type: "local"
  local:
    basePath: "uploads"
    baseURL: "/uploads"
  image:
    maxFileSize: 10485760
    pathPrefix: "articles/images"
  migration:
    enabled: false

search:
  enabled: false
`

func init() {
	// 将工作目录切换到二进制文件所在目录
	if exe, err := os.Executable(); err == nil {
		if dir := filepath.Dir(exe); dir != "" {
			os.Chdir(dir)
		}
	}

	// 确保默认配置文件存在（首次运行自动创建）
	ensureDefaultConfig()
}

func ensureDefaultConfig() {
	configPath := "manifest/config/config.yaml"
	if _, err := os.Stat(configPath); err == nil {
		return // 配置已存在，不覆盖
	}
	os.MkdirAll(filepath.Dir(configPath), 0755)
	os.WriteFile(configPath, []byte(defaultConfig), 0644)
}

func main() {
	cmd.Main.Run(gctx.GetInitCtx())
}
