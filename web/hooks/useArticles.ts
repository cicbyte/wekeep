// 文章管理 Hook - 使用后端 API

import { useState, useEffect, useCallback } from 'react';
import { articlesApi, authorsApi, Article as ApiArticle, Author } from '../services/apiService';
import { Article } from '../types';

// 作者缓存（name -> id 和 id -> name）
const authorNameToIdCache = new Map<string, number>();
const authorIdToNameCache = new Map<number, string>();

// 根据作者名称获取或创建作者
const getOrCreateAuthorId = async (authorName: string): Promise<number> => {
  console.log('getOrCreateAuthorId 调用:', authorName);

  // 检查缓存
  if (authorNameToIdCache.has(authorName)) {
    console.log('从缓存获取作者ID:', authorNameToIdCache.get(authorName));
    return authorNameToIdCache.get(authorName)!;
  }

  try {
    // 先尝试查找现有作者
    const response = await authorsApi.list({ name: authorName, pageSize: 1 });
    console.log('查找作者响应:', response);

    if (response.authorsList && response.authorsList.length > 0) {
      const author = response.authorsList.find(a => a.name === authorName);
      console.log('找到的作者:', author);
      if (author) {
        authorNameToIdCache.set(authorName, author.id);
        authorIdToNameCache.set(author.id, author.name);
        console.log('返回现有作者ID:', author.id);
        return author.id;
      }
    }

    // 没找到则创建新作者
    console.log('未找到现有作者，尝试创建新作者');
    const createResult = await authorsApi.add({ name: authorName });
    console.log('创建作者结果:', createResult);
    const newId = createResult.id;
    authorNameToIdCache.set(authorName, newId);
    authorIdToNameCache.set(newId, authorName);
    console.log('返回新作者ID:', newId);
    return newId;
  } catch (err) {
    console.error('获取或创建作者失败:', err);
    throw err;
  }
};

// 根据作者ID获取名称
const getAuthorNameById = async (authorId: number): Promise<string> => {
  if (authorIdToNameCache.has(authorId)) {
    return authorIdToNameCache.get(authorId)!;
  }
  try {
    const author = await authorsApi.detail(authorId);
    authorIdToNameCache.set(authorId, author.name);
    authorNameToIdCache.set(author.name, authorId);
    return author.name;
  } catch {
    return `作者${authorId}`;
  }
};

// 转换后端数据格式到前端格式
const toFrontendArticle = (apiArticle: ApiArticle): Article => ({
  id: String(apiArticle.id),
  title: apiArticle.title,
  author: apiArticle.authorName || `作者${apiArticle.authorId || 0}`,
  authorId: apiArticle.authorId,
  url: apiArticle.url || '',
  summary: apiArticle.summary || '',
  content: apiArticle.content || '',
  tags: apiArticle.tags || [],
  dateAdded: apiArticle.dateAdded || Date.now(),
  // 搜索高亮字段
  formattedTitle: apiArticle.formattedTitle,
  formattedSummary: apiArticle.formattedSummary,
  contextSnippet: apiArticle.contextSnippet,
  matchFields: apiArticle.matchFields,
});

// 转换前端数据格式到后端格式
const toBackendArticle = (article: Article, authorId: number) => ({
  title: article.title,
  authorId: authorId,
  url: article.url,
  summary: article.summary,
  content: article.content,
  tags: article.tags,
  dateAdded: article.dateAdded,
});

export const useArticles = () => {
  const [articles, setArticles] = useState<Article[]>([]);
  const [authors, setAuthors] = useState<Author[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [totalCount, setTotalCount] = useState(0);
  const [currentPage, setCurrentPage] = useState(1);

  // 加载作者列表
  const loadAuthors = useCallback(async () => {
    try {
      const response = await authorsApi.list({ pageSize: 1000 });
      const authorsList = response?.authorsList || [];
      setAuthors(authorsList);
      // 更新缓存
      authorsList.forEach(author => {
        authorNameToIdCache.set(author.name, author.id);
        authorIdToNameCache.set(author.id, author.name);
      });
      return authorsList;
    } catch (err) {
      console.error('加载作者列表失败:', err);
      return [];
    }
  }, []);

  // 加载文章列表
  const loadArticles = useCallback(async (params?: {
    pageNum?: number;
    pageSize?: number;
    authorId?: number;
    tags?: string[];
    keyword?: string;
    titleOnly?: boolean;
  }) => {
    setLoading(true);
    setError(null);

    try {
      const response = await articlesApi.list(params || {
        pageNum: 1,
        pageSize: 20,
      });

      const articlesList = response?.articlesList || [];
      const frontendArticles = articlesList.map(toFrontendArticle);
      setArticles(frontendArticles);
      setTotalCount(response?.total || 0);
      setCurrentPage(response?.currentPage || 1);
      return frontendArticles;
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : '加载文章失败';
      setError(errorMessage);
      console.error('加载文章失败:', err);
      return [];
    } finally {
      setLoading(false);
    }
  }, []);

  // 添加文章（兼容原有接口：自动根据 author 名称处理 authorId）
  const addArticle = useCallback(async (article: Article) => {
    setLoading(true);
    setError(null);

    try {
      // 获取或创建作者ID
      let authorId = article.authorId;
      if (!authorId && article.author) {
        authorId = await getOrCreateAuthorId(article.author);
      }
      if (!authorId) {
        throw new Error('作者信息缺失');
      }

      await articlesApi.add(toBackendArticle(article, authorId));
      await loadArticles();
      return true;
    } catch (err: any) {
      const errorMessage = err?.message || '添加文章失败';
      setError(errorMessage);
      console.error('添加文章失败:', err);
      throw err; // 重新抛出异常，让调用者处理
    } finally {
      setLoading(false);
    }
  }, [loadArticles]);

  // 更新文章（兼容原有接口）
  const updateArticle = useCallback(async (article: Article) => {
    setLoading(true);
    setError(null);

    try {
      // 获取或创建作者ID
      let authorId = article.authorId;
      if (!authorId && article.author) {
        authorId = await getOrCreateAuthorId(article.author);
      }
      if (!authorId) {
        throw new Error('作者信息缺失');
      }

      const id = parseInt(article.id);
      await articlesApi.edit(id, {
        id,
        ...toBackendArticle(article, authorId),
      });

      setArticles(prev =>
        prev.map(a => (a.id === article.id ? article : a))
      );
      return true;
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : '更新文章失败';
      setError(errorMessage);
      console.error('更新文章失败:', err);
      return false;
    } finally {
      setLoading(false);
    }
  }, []);

  // 删除文章
  const deleteArticle = useCallback(async (id: string) => {
    setLoading(true);
    setError(null);

    try {
      const idNum = parseInt(id);
      await articlesApi.delete(idNum);
      setArticles(prev => prev.filter(a => a.id !== id));
      return true;
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : '删除文章失败';
      setError(errorMessage);
      console.error('删除文章失败:', err);
      return false;
    } finally {
      setLoading(false);
    }
  }, []);

  // 批量删除文章
  const batchDeleteArticles = useCallback(async (ids: string[]) => {
    setLoading(true);
    setError(null);

    try {
      const idNums = ids.map(id => parseInt(id));
      await articlesApi.batchDelete(idNums);
      setArticles(prev => prev.filter(a => !ids.includes(a.id)));
      return true;
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : '批量删除失败';
      setError(errorMessage);
      console.error('批量删除失败:', err);
      return false;
    } finally {
      setLoading(false);
    }
  }, []);

  // 搜索文章
  const searchArticles = useCallback(async (keyword: string, pageNum: number = 1, pageSize: number = 20) => {
    setLoading(true);
    setError(null);

    try {
      const response = await articlesApi.list({
        keyword,
        pageNum,
        pageSize,
      });

      const articlesList = response?.articlesList || [];
      const frontendArticles = articlesList.map(toFrontendArticle);
      setArticles(frontendArticles);
      setTotalCount(response?.total || 0);
      setCurrentPage(response?.currentPage || 1);
      return frontendArticles;
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : '搜索失败';
      setError(errorMessage);
      console.error('搜索失败:', err);
      return [];
    } finally {
      setLoading(false);
    }
  }, []);

  // 按作者筛选（支持 authorId 或 author 名称）
  const filterByAuthor = useCallback(async (authorOrId: string | number, pageNum: number = 1, pageSize: number = 20) => {
    setLoading(true);
    setError(null);

    try {
      let authorId: number | undefined;
      if (typeof authorOrId === 'number') {
        authorId = authorOrId;
      } else {
        // 根据名称查找 authorId
        const cached = authorNameToIdCache.get(authorOrId);
        if (cached) {
          authorId = cached;
        } else {
          const response = await authorsApi.list({ name: authorOrId });
          const found = response.authorsList?.find(a => a.name === authorOrId);
          if (found) {
            authorId = found.id;
          }
        }
      }

      const response = await articlesApi.list({
        authorId,
        pageNum,
        pageSize,
      });

      const articlesList = response?.articlesList || [];
      const frontendArticles = articlesList.map(toFrontendArticle);
      setArticles(frontendArticles);
      setTotalCount(response?.total || 0);
      setCurrentPage(response?.currentPage || 1);
      return frontendArticles;
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : '筛选失败';
      setError(errorMessage);
      console.error('筛选失败:', err);
      return [];
    } finally {
      setLoading(false);
    }
  }, []);

  // 按标签筛选
  const filterByTags = useCallback(async (tags: string[]) => {
    setLoading(true);
    setError(null);

    try {
      const response = await articlesApi.list({
        tags,
        pageNum: 1,
        pageSize: 1000,
      });

      const articlesList = response?.articlesList || [];
      const frontendArticles = articlesList.map(toFrontendArticle);
      setArticles(frontendArticles);
      return frontendArticles;
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : '筛选失败';
      setError(errorMessage);
      console.error('筛选失败:', err);
      return [];
    } finally {
      setLoading(false);
    }
  }, []);

  // 删除作者
  const deleteAuthor = useCallback(async (id: number): Promise<boolean> => {
    setLoading(true);
    setError(null);

    try {
      await authorsApi.delete(id);
      // 重新加载作者列表
      await loadAuthors();
      return true;
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : '删除作者失败';
      setError(errorMessage);
      console.error('删除作者失败:', err);
      return false;
    } finally {
      setLoading(false);
    }
  }, [loadAuthors]);

  // 初始化时加载文章和作者
  useEffect(() => {
    loadAuthors();
    loadArticles();
  }, [loadArticles, loadAuthors]);

  return {
    articles,
    authors,
    loading,
    error,
    loadArticles,
    loadAuthors,
    addArticle,
    updateArticle,
    deleteArticle,
    batchDeleteArticles,
    searchArticles,
    filterByAuthor,
    filterByTags,
    setArticles,
    deleteAuthor,
    totalCount,
    currentPage,
  };
};
