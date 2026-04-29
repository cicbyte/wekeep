package main

import (
	"os"
	"path/filepath"

	_ "github.com/cicbyte/wekeep/internal/packed"
	//重要 需要导入数据库驱动
	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"
	// SQLite数据库驱动,如果需要支持需要go get -u github.com/gogf/gf/contrib/drivers/sqlite/v2
	_ "github.com/gogf/gf/contrib/drivers/sqlite/v2"

	"github.com/gogf/gf/v2/os/gbuild"
	"github.com/gogf/gf/v2/os/gctx"

	"github.com/cicbyte/wekeep/internal/cmd"
)

// 默认配置（首次运行时自动创建，用户可修改）
const defaultConfig = `# WeKeep 配置文件（可修改）
server:
  address: ":8000"
  logPath: "log/server"
  logStdout: true
  errorStack: true
  errorLogEnabled: true
  errorLogPattern: "error-{Ymd}.log"
  accessLogEnabled: true
  accessLogPattern: "access-{Ymd}.log"

logger:
  path: "log/run"
  file: "{Y-m-d}.log"
  level: "all"
  stdout: true

database:
  default:
    link: "sqlite::@file(db/wekeep.db)"

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
	// 通过 gbuild 编译变量判断是否为生产构建（gf build 会注入，gf run 不会）
	if gbuild.Get(gbuild.BuiltVersion) != nil {
		if exe, err := os.Executable(); err == nil {
			if dir := filepath.Dir(exe); dir != "" {
				os.Chdir(dir)
			}
		}
		ensureDefaultConfig()
	}
}

func ensureDefaultConfig() {
	configPath := "manifest/config/config.yaml"
	if _, err := os.Stat(configPath); err == nil {
		return
	}
	os.MkdirAll(filepath.Dir(configPath), 0755)
	os.WriteFile(configPath, []byte(defaultConfig), 0644)
}

func main() {
	cmd.Main.Run(gctx.GetInitCtx())
}
