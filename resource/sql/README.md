# 数据库脚本

```
sql/
├── mysql/init.sql    # MySQL 初始化
├── sqlite/init.sql   # SQLite 初始化
└── README.md
```

## 快速开始

### MySQL

```bash
mysql -u root -p -e "CREATE DATABASE IF NOT EXISTS wekeep DEFAULT CHARSET utf8mb4 COLLATE utf8mb4_unicode_ci"
mysql -u root -p wekeep < resource/sql/mysql/init.sql
```

### SQLite

```bash
sqlite3 wekeep.db < resource/sql/sqlite/init.sql
```

## 数据表

| 表名 | 说明 |
|------|------|
| categories | 分类 |
| authors | 作者，normalized_name 唯一去重 |
| articles | 文章，author_id 外键关联 authors |
| images | 图片，original_url_hash 去重，ref_count 引用计数 |
| article_images | 文章-图片关联，CASCADE 删除 |

## 添加新表

1. 更新 `mysql/init.sql` 和 `sqlite/init.sql`
2. 执行 `make dao` 生成 DAO 代码
