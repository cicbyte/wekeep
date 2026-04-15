// API 测试页面 - 演示后端 API 集成

import React, { useState, useEffect } from 'react';
import { useArticles } from '../hooks/useArticles';
import { articlesApi, categoriesApi, statsApi, searchApi, healthApi } from '../services/apiService';
import { migrateToBackend, needsMigration } from '../services/migrationService';
import { Article } from '../types';
import { Category } from '../services/apiService';

const ApiTestPage: React.FC = () => {
  const [logs, setLogs] = useState<string[]>([]);
  const [testResults, setTestResults] = useState<Record<string, boolean>>({});

  const {
    articles,
    loading,
    error,
    loadArticles,
    addArticle,
    updateArticle,
    deleteArticle,
  } = useArticles();

  const addLog = (message: string) => {
    const timestamp = new Date().toLocaleTimeString();
    setLogs(prev => [`[${timestamp}] ${message}`, ...prev]);
  };

  // 测试健康检查
  const testHealthCheck = async () => {
    addLog('🔍 测试健康检查...');
    try {
      const result = await healthApi.detail();
      addLog(`✅ 健康检查成功: ${result.status}`);
      setTestResults(prev => ({ ...prev, health: true }));
    } catch (error) {
      addLog(`❌ 健康检查失败: ${error}`);
      setTestResults(prev => ({ ...prev, health: false }));
    }
  };

  // 测试加载文章
  const testLoadArticles = async () => {
    addLog('📚 测试加载文章...');
    try {
      await loadArticles();
      addLog(`✅ 加载成功: ${articles.length} 篇文章`);
      setTestResults(prev => ({ ...prev, loadArticles: true }));
    } catch (error) {
      addLog(`❌ 加载失败: ${error}`);
      setTestResults(prev => ({ ...prev, loadArticles: false }));
    }
  };

  // 测试添加文章
  const testAddArticle = async () => {
    addLog('➕ 测试添加文章...');
    const testArticle: Article = {
      id: String(Date.now()),
      title: '测试文章',
      author: '测试作者',
      summary: '这是一篇测试文章',
      content: '## 测试内容\n\n这是测试文章的内容。',
      tags: ['测试', 'API'],
      dateAdded: Date.now(),
    };

    try {
      const success = await addArticle(testArticle);
      if (success) {
        addLog('✅ 添加文章成功');
        setTestResults(prev => ({ ...prev, addArticle: true }));
      } else {
        addLog('❌ 添加文章失败');
        setTestResults(prev => ({ ...prev, addArticle: false }));
      }
    } catch (error) {
      addLog(`❌ 添加文章错误: ${error}`);
      setTestResults(prev => ({ ...prev, addArticle: false }));
    }
  };

  // 测试搜索
  const testSearch = async () => {
    addLog('🔍 测试搜索功能...');
    try {
      const result = await searchApi.search('Go');
      addLog(`✅ 搜索成功: 找到 ${result.total} 个结果`);
      setTestResults(prev => ({ ...prev, search: true }));
    } catch (error) {
      addLog(`❌ 搜索失败: ${error}`);
      setTestResults(prev => ({ ...prev, search: false }));
    }
  };

  // 测试统计
  const testStats = async () => {
    addLog('📊 测试统计功能...');
    try {
      const total = await statsApi.totalArticles();
      addLog(`✅ 文章总数: ${total.count}`);

      const authors = await statsApi.authorStats();
      addLog(`✅ 作者数量: ${authors.authors.length}`);

      const tags = await statsApi.tagStats();
      addLog(`✅ 标签数量: ${tags.tags.length}`);

      setTestResults(prev => ({ ...prev, stats: true }));
    } catch (error) {
      addLog(`❌ 统计失败: ${error}`);
      setTestResults(prev => ({ ...prev, stats: false }));
    }
  };

  // 测试数据迁移
  const testMigration = async () => {
    addLog('🔄 测试数据迁移...');
    if (!needsMigration()) {
      addLog('ℹ️  无需迁移（localStorage 为空）');
      return;
    }

    try {
      const result = await migrateToBackend();
      addLog(`✅ 迁移完成: 成功 ${result.success}, 失败 ${result.failed}`);
      if (result.errors.length > 0) {
        addLog(`⚠️  错误列表:`);
        result.errors.forEach(err => addLog(`   - ${err.title}: ${err.error}`));
      }
      setTestResults(prev => ({ ...prev, migration: true }));
    } catch (error) {
      addLog(`❌ 迁移失败: ${error}`);
      setTestResults(prev => ({ ...prev, migration: false }));
    }
  };

  // 运行所有测试
  const runAllTests = async () => {
    addLog('🚀 开始运行所有测试...');
    await testHealthCheck();
    await testLoadArticles();
    await testStats();
    await testSearch();
    addLog('✅ 所有测试完成');
  };

  // 清除日志
  const clearLogs = () => {
    setLogs([]);
    setTestResults({});
  };

  return (
    <div className="p-6 max-w-6xl mx-auto">
      <h1 className="text-3xl font-bold mb-6">API 测试页面</h1>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4 mb-6">
        <button
          onClick={testHealthCheck}
          className="bg-blue-500 hover:bg-blue-600 text-white px-4 py-2 rounded"
        >
          测试健康检查
        </button>

        <button
          onClick={testLoadArticles}
          className="bg-green-500 hover:bg-green-600 text-white px-4 py-2 rounded"
        >
          测试加载文章
        </button>

        <button
          onClick={testAddArticle}
          className="bg-yellow-500 hover:bg-yellow-600 text-white px-4 py-2 rounded"
        >
          测试添加文章
        </button>

        <button
          onClick={testSearch}
          className="bg-purple-500 hover:bg-purple-600 text-white px-4 py-2 rounded"
        >
          测试搜索
        </button>

        <button
          onClick={testStats}
          className="bg-pink-500 hover:bg-pink-600 text-white px-4 py-2 rounded"
        >
          测试统计
        </button>

        <button
          onClick={testMigration}
          className="bg-orange-500 hover:bg-orange-600 text-white px-4 py-2 rounded"
        >
          测试数据迁移
        </button>

        <button
          onClick={runAllTests}
          className="col-span-full bg-indigo-600 hover:bg-indigo-700 text-white px-4 py-2 rounded font-bold"
        >
          🚀 运行所有测试
        </button>

        <button
          onClick={clearLogs}
          className="col-span-full bg-gray-500 hover:bg-gray-600 text-white px-4 py-2 rounded"
        >
          清除日志
        </button>
      </div>

      {/* 测试结果 */}
      <div className="mb-6">
        <h2 className="text-xl font-bold mb-3">测试结果</h2>
        <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-6 gap-2">
          {Object.entries(testResults).map(([test, passed]) => (
            <div
              key={test}
              className={`px-3 py-2 rounded text-center ${
                passed ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'
              }`}
            >
              {passed ? '✅' : '❌'} {test}
            </div>
          ))}
        </div>
      </div>

      {/* 当前状态 */}
      <div className="mb-6 p-4 bg-gray-100 rounded">
        <h2 className="text-xl font-bold mb-2">当前状态</h2>
        <div className="space-y-1">
          <p>文章数量: <strong>{articles.length}</strong></p>
          <p>加载状态: <strong>{loading ? '加载中...' : '已就绪'}</strong></p>
          <p>错误信息: <strong>{error || '无'}</strong></p>
          <p>需要迁移: <strong>{needsMigration() ? '是' : '否'}</strong></p>
        </div>
      </div>

      {/* 日志 */}
      <div>
        <h2 className="text-xl font-bold mb-3">日志</h2>
        <div className="bg-gray-900 text-green-400 p-4 rounded font-mono text-sm h-96 overflow-y-auto">
          {logs.length === 0 ? (
            <p className="text-gray-500">等待测试...</p>
          ) : (
            logs.map((log, index) => (
              <div key={index} className="mb-1">
                {log}
              </div>
            ))
          )}
        </div>
      </div>

      {/* 文章列表 */}
      {articles.length > 0 && (
        <div className="mt-6">
          <h2 className="text-xl font-bold mb-3">当前文章列表</h2>
          <div className="grid gap-2">
            {articles.slice(0, 5).map(article => (
              <div key={article.id} className="p-3 bg-white border rounded">
                <h3 className="font-bold">{article.title}</h3>
                <p className="text-sm text-gray-600">作者: {article.author}</p>
                <p className="text-sm text-gray-600">
                  标签: {article.tags.join(', ')}
                </p>
              </div>
            ))}
            {articles.length > 5 && (
              <p className="text-center text-gray-500">
                还有 {articles.length - 5} 篇文章...
              </p>
            )}
          </div>
        </div>
      )}
    </div>
  );
};

export default ApiTestPage;
