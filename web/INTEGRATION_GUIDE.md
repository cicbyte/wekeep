# 前端与后端 API 集成指南

本文档说明如何将现有的 React 前端应用与 Go 后端 API 集成。

## 📦 已创建的服务

### 1. API 服务 (`services/apiService.ts`)

封装了所有后端 API 调用，包括：

- **文章 API**: `articlesApi.list()`, `add()`, `edit()`, `delete()`, `batchDelete()`
- **分类 API**: `categoriesApi.list()`, `add()`, `edit()`, `delete()`, `batchDelete()`
- **统计 API**: `statsApi.totalArticles()`, `authorStats()`, `tagStats()`, `timeTrends()`
- **搜索 API**: `searchApi.search()` (Meilisearch 全文搜索)
- **健康检查**: `healthApi.check()`, `detail()`

### 2. 数据迁移服务 (`services/migrationService.ts`)

提供数据迁移功能：

- `loadFromLocalStorage()` - 从 localStorage 读取数据
- `migrateToBackend()` - 迁移数据到后端
- `loadFromBackend()` - 从后端加载数据
- `syncArticles()` - 智能同步策略

### 3. 自定义 Hook (`hooks/useArticles.ts`)

React Hook 管理文章状态：

```typescript
const {
  articles,           // 文章列表
  loading,            // 加载状态
  error,              // 错误信息
  loadArticles,       // 加载文章
  addArticle,         // 添加文章
  updateArticle,      // 更新文章
  deleteArticle,      // 删除文章
  batchDeleteArticles,// 批量删除
  searchArticles,     // 搜索
  filterByAuthor,     // 按作者筛选
  filterByTags,       // 按标签筛选
} = useArticles();
```

## 🔧 集成步骤

### 步骤 1: 在 App.tsx 中引入 Hook

```typescript
import { useArticles } from './hooks/useArticles';

const App: React.FC = () => {
  // 替换原有的 useState
  const [articles, setArticles] = useState<Article[]>(() => {
    const saved = localStorage.getItem('wekeep_articles');
    return saved ? JSON.parse(saved) : SAMPLE_ARTICLES;
  });

  // 改为使用自定义 Hook
  const {
    articles,
    loading,
    error,
    addArticle,
    updateArticle,
    deleteArticle,
    searchArticles,
  } = useArticles();
```

### 步骤 2: 修改数据持久化逻辑

**原有代码 (localStorage):**
```typescript
useEffect(() => {
  localStorage.setItem('wekeep_articles', JSON.stringify(articles));
}, [articles]);
```

**新代码 (后端 API):**
```typescript
// 不再需要 localStorage 的 useEffect
// 数据由 useArticles Hook 自动管理
```

### 步骤 3: 更新事件处理函数

**添加文章:**
```typescript
// 原有代码
const handleAddArticle = (article: Article) => {
  setArticles(prev => [article, ...prev]);
};

// 新代码
const handleAddArticle = async (article: Article) => {
  const success = await addArticle(article);
  if (success) {
    // 文章已自动添加到列表
    console.log('文章添加成功');
  }
};
```

**更新文章:**
```typescript
// 原有代码
const handleUpdateArticle = (updatedArticle: Article) => {
  setArticles(prev =>
    prev.map(article =>
      article.id === updatedArticle.id ? updatedArticle : article
    )
  );
};

// 新代码
const handleUpdateArticle = async (updatedArticle: Article) => {
  const success = await updateArticle(updatedArticle);
  if (success) {
    console.log('文章更新成功');
  }
};
```

**删除文章:**
```typescript
// 原有代码
const confirmDelete = () => {
  if (articleToDelete) {
    setArticles(prev => prev.filter(a => a.id !== articleToDelete));
    setArticleToDelete(null);
  }
};

// 新代码
const confirmDelete = async () => {
  if (articleToDelete) {
    const success = await deleteArticle(articleToDelete);
    if (success) {
      setArticleToDelete(null);
      console.log('文章删除成功');
    }
  }
};
```

**搜索文章:**
```typescript
// 原有代码 (客户端过滤)
const filteredArticles = useMemo(() => {
  return articles.filter(article =>
    article.title.toLowerCase().includes(searchQuery.toLowerCase()) ||
    article.author.toLowerCase().includes(searchQuery.toLowerCase())
  );
}, [articles, searchQuery]);

// 新代码 (后端搜索)
const handleSearch = async (query: string) => {
  if (query.trim()) {
    await searchArticles(query);
  } else {
    await loadArticles();
  }
};
```

### 步骤 4: 数据迁移

**选项 A: 手动迁移按钮**

在设置页面添加迁移功能：

```typescript
import { migrateToBackend, needsMigration } from './services/migrationService';

const SettingsPanel = () => {
  const handleMigrate = async () => {
    const result = await migrateToBackend();
    console.log(`迁移完成: 成功 ${result.success}, 失败 ${result.failed}`);
  };

  return (
    <div>
      {needsMigration() && (
        <button onClick={handleMigrate}>
          迁移本地数据到服务器
        </button>
      )}
    </div>
  );
};
```

**选项 B: 自动迁移**

在应用启动时自动迁移：

```typescript
useEffect(() => {
  const initApp = async () => {
    if (needsMigration()) {
      console.log('检测到本地数据，开始迁移...');
      await migrateToBackend();
    }
    await loadArticles();
  };

  initApp();
}, []);
```

## 🎯 完整示例

```typescript
import React, { useState, useEffect } from 'react';
import { useArticles } from './hooks/useArticles';
import { migrateToBackend, needsMigration } from './services/migrationService';

const App: React.FC = () => {
  const {
    articles,
    loading,
    error,
    loadArticles,
    addArticle,
    updateArticle,
    deleteArticle,
  } = useArticles();

  const [searchQuery, setSearchQuery] = useState('');

  // 初始化应用
  useEffect(() => {
    const init = async () => {
      // 自动迁移数据
      if (needsMigration()) {
        console.log('正在迁移本地数据...');
        await migrateToBackend();
      }
      // 加载文章
      await loadArticles();
    };

    init();
  }, []);

  // 搜索处理
  const handleSearch = async (query: string) => {
    setSearchQuery(query);
    if (query.trim()) {
      // 使用后端搜索
      await searchArticles(query);
    } else {
      await loadArticles();
    }
  };

  // 错误处理
  if (error) {
    return <div>错误: {error}</div>;
  }

  // 加载状态
  if (loading) {
    return <div>加载中...</div>;
  }

  return (
    <div>
      <h1>WeKeep 文章管理</h1>
      <input
        type="text"
        value={searchQuery}
        onChange={(e) => handleSearch(e.target.value)}
        placeholder="搜索文章..."
      />
      <ArticleList
        articles={articles}
        onAdd={addArticle}
        onUpdate={updateArticle}
        onDelete={deleteArticle}
      />
    </div>
  );
};
```

## 📝 注意事项

1. **数据格式转换**: Hook 会自动处理前后端数据格式的转换
2. **错误处理**: 所有 API 调用都有错误处理，通过 `error` 状态获取
3. **加载状态**: 使用 `loading` 状态显示加载指示器
4. **向后兼容**: 可以保留 localStorage 作为备用存储
5. **渐进式迁移**: 可以先迁移部分功能，逐步替换

## 🚀 测试步骤

1. 确保后端服务器运行在 `http://localhost:8000`
2. 测试健康检查: `curl http://localhost:8000/api/v1/health`
3. 在浏览器中打开前端应用
4. 尝试添加、编辑、删除文章
5. 检查浏览器控制台的 API 请求日志
6. 验证数据是否正确保存到数据库

## 📚 相关文件

- `web/services/apiService.ts` - API 服务定义
- `web/services/migrationService.ts` - 数据迁移服务
- `web/hooks/useArticles.ts` - 文章管理 Hook
- `web/types.ts` - 类型定义

## 🔄 从 localStorage 迁移到后端

如果您已有 localStorage 数据，可以：

1. 在设置页面添加"导出数据"按钮（保留原有功能）
2. 添加"迁移到服务器"按钮
3. 迁移后可以选择清除本地数据
4. 保留导出功能作为备份
