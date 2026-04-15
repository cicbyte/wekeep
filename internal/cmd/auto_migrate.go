package cmd

import (
	"context"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/gres"
)

// autoMigrate 启动时自动检查并创建数据表
func autoMigrate(ctx context.Context) {
	// 获取数据库类型
	dbType := g.Cfg().MustGet(ctx, "database.default.link").String()
	if dbType == "" {
		return
	}

	var sqlFile string
	var checkSQL string

	switch {
	case strings.HasPrefix(dbType, "mysql"):
		sqlFile = "resource/sql/mysql/init.sql"
		checkSQL = "SHOW TABLES LIKE 'categories'"
	case strings.HasPrefix(dbType, "sqlite"):
		sqlFile = "resource/sql/sqlite/init.sql"
		checkSQL = "SELECT name FROM sqlite_master WHERE type='table' AND name='categories'"
	default:
		return
	}

	// 检查表是否已存在
	result, err := g.DB().Query(ctx, checkSQL)
	if err != nil {
		g.Log().Warningf(ctx, "Auto migrate check failed: %v", err)
		return
	}

	if len(result) > 0 {
		g.Log().Info(ctx, "Database tables already exist, skip auto migrate")
		return
	}

	// 表不存在，读取 init.sql（优先 gres，降级 gfile）
	var content string
	if file := gres.Get(sqlFile); file != nil {
		content = string(file.Content())
	} else if gfile.Exists(sqlFile) {
		content = gfile.GetContents(sqlFile)
	} else {
		g.Log().Errorf(ctx, "SQL file not found: %s", sqlFile)
		return
	}
	g.Log().Info(ctx, "Database tables not found, running auto migrate...")

	for _, stmt := range splitSQL(content) {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" || strings.HasPrefix(stmt, "--") {
			continue
		}
		if _, err := g.DB().Exec(ctx, stmt); err != nil {
			g.Log().Errorf(ctx, "Auto migrate failed: %v, SQL: %s", err, stmt)
			return
		}
	}

	g.Log().Info(ctx, "Auto migrate completed successfully")
}

// splitSQL 按分号分割 SQL 语句，跳过注释，正确处理 BEGIN...END 块
func splitSQL(sql string) []string {
	var statements []string
	var current strings.Builder
	inBlock := false

	for _, line := range strings.Split(sql, "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "--") {
			continue
		}
		current.WriteString(line)
		current.WriteString("\n")

		upper := strings.ToUpper(trimmed)
		if strings.Contains(upper, "BEGIN") {
			inBlock = true
		}
		if inBlock && strings.Contains(upper, "END") && strings.HasSuffix(upper, ";") {
			inBlock = false
			statements = append(statements, current.String())
			current.Reset()
			continue
		}
		if !inBlock && strings.HasSuffix(trimmed, ";") {
			statements = append(statements, current.String())
			current.Reset()
		}
	}

	if current.Len() > 0 {
		if remaining := strings.TrimSpace(current.String()); remaining != "" {
			statements = append(statements, remaining)
		}
	}

	return statements
}
