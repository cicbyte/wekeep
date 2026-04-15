// API 服务 - 与后端 Go 框架通信

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || '/api/v1';

// 后端响应格式
interface ApiResponse<T> {
  code: number;
  message: string;
  data: T;
}

// 通用请求函数
async function request<T>(
  endpoint: string,
  options: RequestInit = {}
): Promise<T> {
  const url = `${API_BASE_URL}${endpoint}`;

  const config: RequestInit = {
    ...options,
    headers: {
      'Content-Type': 'application/json',
      ...options.headers,
    },
  };

  try {
    const response = await fetch(url, config);

    if (!response.ok) {
      const error = await response.json().catch(() => ({
        message: `HTTP ${response.status}: ${response.statusText}`,
      }));
      throw new Error(error.message || '请求失败');
    }

    const json: ApiResponse<T> = await response.json();

    // 检查业务状态码
    if (json.code !== 0) {
      throw new Error(json.message || '请求失败');
    }

    // 返回 data 字段
    return json.data;
  } catch (error) {
    console.error(`API 请求失败: ${endpoint}`, error);
    throw error;
  }
}

// 文章相关 API
export const articlesApi = {
  // 通过 URL 解析微信文章
  parseByUrl: (url: string) =>
    request<{
      title: string;
      author: string;
      content: string;
      baseUrl: string;
      publishTime: string;
    }>('/articles/parse-by-url', {
      method: 'POST',
      body: JSON.stringify({ url }),
    }).catch(err => {
      console.error('解析微信文章失败:', err);
      throw new Error(`解析失败: ${err.message || '未知错误'}`);
    }),

  // 解析微信文章
  parse: (htmlContent: string, baseUrl?: string) =>
    request<{
      title: string;
      author: string;
      content: string;
      baseUrl: string;
      publishTime: string;
    }>('/articles/parse', {
      method: 'POST',
      body: JSON.stringify({ htmlContent, baseUrl }),
    }).catch(err => {
      console.error('解析微信文章失败:', err);
      throw new Error(`解析失败: ${err.message || '未知错误'}`);
    }),

  // 获取文章列表
  list: (params?: {
    pageNum?: number;
    pageSize?: number;
    authorId?: number;       // 改为 authorId
    tags?: string[];
    keyword?: string;
    titleOnly?: boolean;     // 仅搜索标题
    orderBy?: string;
  }) => {
    const queryParams = new URLSearchParams();
    if (params?.pageNum) queryParams.set('pageNum', String(params.pageNum));
    if (params?.pageSize) queryParams.set('pageSize', String(params.pageSize));
    if (params?.authorId) queryParams.set('authorId', String(params.authorId));
    if (params?.keyword) queryParams.set('keyword', params.keyword);
    if (params?.titleOnly) queryParams.set('titleOnly', 'true');
    if (params?.orderBy) queryParams.set('orderBy', params.orderBy);
    if (params?.tags?.length) {
      params.tags.forEach(tag => queryParams.append('tags[]', tag));
    }
    const query = queryParams.toString();
    return request<{
      total: number;
      currentPage: number;
      articlesList: Article[];
    }>(`/articles/list${query ? `?${query}` : ''}`, {
      method: 'GET',
    });
  },

  // 获取文章详情
  detail: (id: number) =>
    request<Article>(`/articles/detail?id=${id}`, {
      method: 'GET',
    }),

  // 新增文章
  add: (data: {
    title: string;
    authorId: number;         // 改为 authorId
    url?: string;
    summary?: string;
    content?: string;
    tags?: string[];
    dateAdded?: number;
  }) =>
    request<{ id: number }>('/articles/add', {
      method: 'POST',
      body: JSON.stringify(data),
    }),

  // 编辑文章
  edit: (id: number, data: {
    id: number;
    title: string;
    authorId: number;         // 改为 authorId
    url?: string;
    summary?: string;
    content?: string;
    tags?: string[];
    dateAdded?: number;
  }) =>
    request<void>(`/articles/edit`, {
      method: 'PUT',
      body: JSON.stringify(data),
    }),

  // 删除文章
  delete: (id: number) =>
    request<void>(`/articles/del?id=${id}`, {
      method: 'DELETE',
    }),

  // 批量删除文章
  batchDelete: (ids: number[]) =>
    request<void>('/articles/batchdel', {
      method: 'DELETE',
      body: JSON.stringify({ ids }),
    }),

  // 重新解析文章
  reparse: (id: number) =>
    request<{
      title: string;
      author: string;
      content: string;
      imagesCount: number;
    }>('/articles/reparse', {
      method: 'POST',
      body: JSON.stringify({ id }),
    }),
};

// 分类相关 API
export const categoriesApi = {
  // 获取分类列表
  list: (params?: { pageNum?: number; pageSize?: number; name?: string }) =>
    request<{
      total: number;
      currentPage: number;
      categoriesList: Category[];
    }>('/categories/list', {
      method: 'GET',
    }),

  // 获取分类详情
  detail: (id: number) =>
    request<Category>(`/categories/detail?id=${id}`, {
      method: 'GET',
    }),

  // 新增分类
  add: (data: { name: string; description: string; icon?: string; sort?: number }) =>
    request<void>('/categories/add', {
      method: 'POST',
      body: JSON.stringify(data),
    }),

  // 编辑分类
  edit: (
    id: number,
    data: { name?: string; description?: string; icon?: string; sort?: number }
  ) =>
    request<void>(`/categories/edit?id=${id}`, {
      method: 'PUT',
      body: JSON.stringify(data),
    }),

  // 删除分类
  delete: (id: number) =>
    request<void>(`/categories/del?id=${id}`, {
      method: 'DELETE',
    }),

  // 批量删除分类
  batchDelete: (ids: number[]) =>
    request<void>('/categories/batchdel', {
      method: 'DELETE',
      body: JSON.stringify({ ids }),
    }),

  // 公开分类列表（无需认证）
  public: () =>
    request<Category[]>('/categories', {
      method: 'GET',
    }),
};

// 作者相关 API
export const authorsApi = {
  // 获取作者列表
  list: (params?: { pageNum?: number; pageSize?: number; name?: string }) => {
    const queryParams = new URLSearchParams();
    if (params?.pageNum) queryParams.set('pageNum', String(params.pageNum));
    if (params?.pageSize) queryParams.set('pageSize', String(params.pageSize));
    if (params?.name) queryParams.set('name', params.name);
    const query = queryParams.toString();
    return request<{
      total: number;
      currentPage: number;
      authorsList: Author[];
    }>(`/authors/list${query ? `?${query}` : ''}`, {
      method: 'GET',
    });
  },

  // 获取作者详情
  detail: (id: number) =>
    request<Author>(`/authors/detail?id=${id}`, {
      method: 'GET',
    }),

  // 新增作者
  add: (data: { name: string; avatar?: string; bio?: string; website?: string }) =>
    request<{ id: number }>('/authors/add', {
      method: 'POST',
      body: JSON.stringify(data),
    }),

  // 编辑作者
  edit: (id: number, data: { name: string; avatar?: string; bio?: string; website?: string }) =>
    request<void>(`/authors/edit`, {
      method: 'PUT',
      body: JSON.stringify({ id, ...data }),
    }),

  // 删除作者
  delete: (id: number) =>
    request<void>(`/authors/del?id=${id}`, {
      method: 'DELETE',
    }),

  // 获取作者选择选项（用于下拉框）
  select: () =>
    request<{ options: { id: number; name: string; articleCount: number }[] }>('/authors/select', {
      method: 'GET',
    }),
};

// 统计相关 API
export const statsApi = {
  // 获取文章总数
  totalArticles: () =>
    request<{ count: number }>('/articles/stats/total', {
      method: 'GET',
    }),

  // 获取作者统计
  authorStats: () =>
    request<{ authors: AuthorStats[] }>('/articles/stats/authors', {
      method: 'GET',
    }),

  // 获取标签统计
  tagStats: () =>
    request<{ tags: TagStats[] }>('/articles/stats/tags', {
      method: 'GET',
    }),

  // 获取时间趋势
  timeTrends: (days?: number) =>
    request<{ trends: TimeTrend[] }>(
      `/articles/stats/trends${days ? `?days=${days}` : ''}`,
      {
        method: 'GET',
      }
    ),
};

// 搜索相关 API
export const searchApi = {
  // 全文搜索
  search: (query: string, params?: { limit?: number; author?: string }) =>
    request<{
      total: number;
      hits: SearchHit[];
    }>(
      `/search?query=${encodeURIComponent(query)}${params?.limit ? `&limit=${params.limit}` : ''}${params?.author ? `&author=${encodeURIComponent(params.author)}` : ''}`,
      {
        method: 'GET',
      }
    ),

  // 获取搜索状态
  status: () =>
    request<{
      enabled: boolean;
      indexedCount: number;
    }>('/search/status', {
      method: 'GET',
    }),

  // 索引所有文章（管理接口）
  indexAll: () =>
    request<{
      msg: string;
      indexedCount: number;
    }>('/search/indexAll', {
      method: 'POST',
    }),
};

// 健康检查 API
export const healthApi = {
  // 简单健康检查
  check: () =>
    request<{ status: string; message: string }>('/health', {
      method: 'GET',
    }),

  // 详细健康检查
  detail: () =>
    request<{
      status: string;
      message: string;
      checks?: Array<{ name: string; status: string }>;
      uptime?: string;
    }>('/health/detail', {
      method: 'GET',
    }),

  // 版本信息
  version: () =>
    request<{ version: string }>('/health/version', {
      method: 'GET',
    }),
};

// 类型定义
export interface Article {
  id: number;
  title: string;
  authorId?: number;          // 作者ID
  authorName?: string;        // 作者名称（用于显示）
  url?: string;
  summary?: string;
  content?: string;
  tags: string[];
  dateAdded: number;
  createdAt?: number;
  updatedAt?: number;
  // 搜索高亮字段
  formattedTitle?: string;    // 高亮后的标题
  formattedSummary?: string;  // 高亮后的摘要
  contextSnippet?: string;    // 正文匹配上下文
  matchFields?: string[];     // 匹配的字段列表
}

export interface Category {
  id: number;
  name: string;
  description: string;
  icon?: string;
  sort: number;
  createdAt?: number;
  updatedAt?: number;
}

// 作者类型
export interface Author {
  id: number;
  name: string;
  avatar?: string;
  bio?: string;
  website?: string;
  articleCount: number;
  createdAt?: string;
}

export interface AuthorStats {
  authorId?: number;
  authorName: string;
  count: number;
}

export interface TagStats {
  tag: string;
  count: number;
}

export interface TimeTrend {
  date: string;
  count: number;
}

export interface SearchHit {
  id: number;
  title: string;
  author: string;
  summary?: string;
  url?: string;
  tags: string[];
  dateAdded: number;
  score: number;
}

// 存储相关类型
export interface MigrationStatus {
  running: boolean;
  source: string;
  target: string;
  total: number;
  completed: number;
  failed: number;
  currentFile: string;
  startTime: string;
  endTime: string;
  error: string;
}

export interface StorageStatus {
  currentStorage: string;
  initialized: boolean;
  migration: MigrationStatus;
}

export interface StorageStats {
  totalImages: number;
  totalSize: number;
  localSize: number;
  byStorageType: Record<string, number>;
}

export interface LocalConfig {
  basePath?: string;
  baseURL?: string;
}

export interface RustFSConfig {
  endpoint: string;
  bucket: string;
  username: string;
  password: string;
  timeout?: number;
}

export interface MigrationConfig {
  targetStorage: 'local' | 'rustfs';
  localConfig?: LocalConfig;
  rustfsConfig?: RustFSConfig;
  updateMarkdown?: boolean;
}

// 存储 API
export const storageApi = {
  // 获取存储状态
  status: () =>
    request<StorageStatus>('/storage/status', {
      method: 'GET',
    }),

  // 获取存储统计
  stats: () =>
    request<StorageStats>('/storage/stats', {
      method: 'GET',
    }),

  // 验证存储配置
  validate: (config: MigrationConfig) =>
    request<{ valid: boolean; error?: string }>('/storage/validate', {
      method: 'POST',
      body: JSON.stringify(config),
    }),

  // 执行迁移
  migrate: (config: MigrationConfig) =>
    request<{ message: string }>('/storage/migrate', {
      method: 'POST',
      body: JSON.stringify(config),
    }),

  // 切换存储
  switch: (config: MigrationConfig) =>
    request<{ message: string; requires?: string }>('/storage/switch', {
      method: 'POST',
      body: JSON.stringify(config),
    }),
};
