// 数据迁移服务 - 将 localStorage 数据迁移到后端 API

import { articlesApi } from './apiService';
import { Article } from '../types';

// 从 localStorage 读取文章数据
export const loadFromLocalStorage = (): Article[] => {
  try {
    const saved = localStorage.getItem('wekeep_articles');
    if (saved) {
      return JSON.parse(saved);
    }
  } catch (error) {
    console.error('读取 localStorage 失败:', error);
  }
  return [];
};

// 将文章保存到 localStorage
export const saveToLocalStorage = (articles: Article[]) => {
  try {
    localStorage.setItem('wekeep_articles', JSON.stringify(articles));
  } catch (error) {
    console.error('保存到 localStorage 失败:', error);
  }
};

// 清除 localStorage 中的文章数据
export const clearLocalStorage = () => {
  try {
    localStorage.removeItem('wekeep_articles');
  } catch (error) {
    console.error('清除 localStorage 失败:', error);
  }
};

// 迁移数据到后端
export const migrateToBackend = async (): Promise<{
  success: number;
  failed: number;
  errors: Array<{ title: string; error: string }>;
}> => {
  const localArticles = loadFromLocalStorage();
  const result = {
    success: 0,
    failed: 0,
    errors: [] as Array<{ title: string; error: string }>,
  };

  console.log(`开始迁移 ${localArticles.length} 篇文章...`);

  for (const article of localArticles) {
    try {
      // 转换数据格式 - 使用 author 字符串通过 API 自动创建/查找作者
      await articlesApi.add({
        title: article.title,
        authorId: 0, // 将由后端根据 author 自动创建
        url: article.url,
        summary: article.summary,
        content: article.content,
        tags: article.tags,
        dateAdded: article.dateAdded,
      });
      result.success++;
      console.log(`✓ 迁移成功: ${article.title}`);
    } catch (error) {
      result.failed++;
      result.errors.push({
        title: article.title,
        error: error instanceof Error ? error.message : '未知错误',
      });
      console.error(`✗ 迁移失败: ${article.title}`, error);
    }

    // 添加延迟避免请求过快
    await new Promise(resolve => setTimeout(resolve, 100));
  }

  console.log(`迁移完成: 成功 ${result.success}, 失败 ${result.failed}`);

  return result;
};

// 检查是否需要迁移
export const needsMigration = (): boolean => {
  const localArticles = loadFromLocalStorage();
  return localArticles.length > 0;
};

// 从后端加载文章
export const loadFromBackend = async (): Promise<Article[]> => {
  try {
    const response = await articlesApi.list({
      pageNum: 1,
      pageSize: 1000,
    });

    // 转换后端数据格式到前端格式
    return response.articlesList.map(article => ({
      id: String(article.id),
      title: article.title,
      author: article.authorName || `作者${article.authorId || 0}`,
      authorId: article.authorId,
      url: article.url,
      summary: article.summary,
      content: article.content,
      tags: article.tags || [],
      dateAdded: article.dateAdded || Date.now(),
    }));
  } catch (error) {
    console.error('从后端加载文章失败:', error);
    throw error;
  }
};

// 同步策略：优先使用后端数据，如果后端为空则提示迁移
export const syncArticles = async (): Promise<Article[]> => {
  try {
    // 尝试从后端加载
    const backendArticles = await loadFromBackend();

    // 如果后端有数据，直接使用
    if (backendArticles.length > 0) {
      console.log(`从后端加载了 ${backendArticles.length} 篇文章`);
      return backendArticles;
    }

    // 如果后端为空但本地有数据，提示迁移
    if (needsMigration()) {
      console.log('后端为空，检测到本地数据，建议迁移');
      // 这里可以自动迁移或提示用户
      // return await migrateAndLoad();
    }

    return [];
  } catch (error) {
    console.error('同步失败，回退到 localStorage:', error);

    // 如果后端请求失败，回退到 localStorage
    return loadFromLocalStorage();
  }
};

// 迁移并加载数据
export const migrateAndLoad = async (): Promise<Article[]> => {
  const result = await migrateToBackend();

  if (result.success > 0) {
    // 迁移成功后，从后端重新加载数据
    const articles = await loadFromBackend();

    // 可选：清除本地数据
    // clearLocalStorage();

    return articles;
  }

  throw new Error('迁移失败');
};
