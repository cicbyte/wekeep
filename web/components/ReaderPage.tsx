import React, { useState, useEffect, useRef } from 'react';
import { useParams, useNavigate, Link } from 'react-router-dom';
import { ArrowLeft, ExternalLink, Calendar, User, Edit3, Moon, Sun, Home, Palette, ChevronDown, Check } from 'lucide-react';
import { Article } from '../types';
import { MdPreview, MdCatalog } from 'md-editor-rt';
import 'md-editor-rt/lib/preview.css';
import { Loader2 } from 'lucide-react';
import { articlesApi } from '../services/apiService';

const IMAGE_BASE_URL = (import.meta.env.VITE_API_BASE_URL || '/api/v1') + '/images/file/';

// 处理 Markdown 内容中的相对图片路径
function processContent(content: string): string {
  return content.replace(
    /!\[([^\]]*)\]\(([^)]+)\)/g,
    (_, alt, src) => {
      if (src.startsWith('http://') || src.startsWith('https://') || src.startsWith('/')) return `![${alt}](${src})`;
      return `![${alt}](${IMAGE_BASE_URL}${src})`;
    }
  );
}

const themeOptions = [
  { value: 'github', label: 'GitHub' },
  { value: 'vuepress', label: 'VuePress' },
  { value: 'mk-cute', label: 'Cute' },
  { value: 'smart-blue', label: 'Smart Blue' },
  { value: 'cyanosis', label: 'Cyanosis' },
];

const ReaderPage: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const themeDropdownRef = useRef<HTMLDivElement>(null);

  const [article, setArticle] = useState<Article | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const [isDark, setIsDark] = useState(() => {
    const saved = localStorage.getItem('wekeep_dark_mode');
    if (saved !== null) return saved === 'true';
    return window.matchMedia('(prefers-color-scheme: dark)').matches;
  });
  const [previewTheme, setPreviewTheme] = useState<'default' | 'github' | 'vuepress' | 'mk-cute' | 'smart-blue' | 'cyanosis'>('github');
  const [showCatalog, setShowCatalog] = useState(false);
  const [showThemeDropdown, setShowThemeDropdown] = useState(false);
  const editorId = 'reader-preview';

  // 加载文章数据
  useEffect(() => {
    const loadArticle = async () => {
      if (!id) {
        setError('文章ID不存在');
        setLoading(false);
        return;
      }

      setLoading(true);
      setError(null);

      try {
        const response = await articlesApi.list({ pageNum: 1, pageSize: 1000 });
        const found = response?.articlesList?.find(a => String(a.id) === id);

        if (found) {
          setArticle({
            id: String(found.id),
            title: found.title,
            author: found.authorName || `作者${found.authorId || 0}`,
            authorId: found.authorId,
            url: found.url || '',
            summary: found.summary || '',
            content: found.content || '',
            tags: found.tags || [],
            dateAdded: found.dateAdded || Date.now(),
          });
        } else {
          setError('文章不存在');
        }
      } catch (err) {
        console.error('加载文章失败:', err);
        setError('加载文章失败');
      } finally {
        setLoading(false);
      }
    };

    loadArticle();
  }, [id]);

  // 滚动到顶部
  useEffect(() => {
    window.scrollTo(0, 0);
  }, [id]);

  const handleEdit = () => {
    if (article) {
      navigate(`/edit/${article.id}`);
    }
  };

  const handleBack = () => {
    // 如果有历史记录则返回，否则跳转到列表
    if (window.history.length > 1) {
      navigate(-1);
    } else {
      navigate('/list');
    }
  };

  // 点击外部关闭主题下拉菜单
  useEffect(() => {
    const handleClickOutside = (e: MouseEvent) => {
      if (themeDropdownRef.current && !themeDropdownRef.current.contains(e.target as Node)) {
        setShowThemeDropdown(false);
      }
    };
    if (showThemeDropdown) {
      document.addEventListener('mousedown', handleClickOutside);
      return () => document.removeEventListener('mousedown', handleClickOutside);
    }
  }, [showThemeDropdown]);

  const handleThemeSelect = (theme: typeof previewTheme) => {
    setPreviewTheme(theme);
    setShowThemeDropdown(false);
  };

  // 加载状态
  if (loading) {
    return (
      <div className={`min-h-screen flex items-center justify-center ${isDark ? 'bg-slate-900' : 'bg-white'}`}>
        <div className="flex flex-col items-center gap-3">
          <Loader2 size={48} className="text-wechat animate-spin" />
          <p className={`text-sm ${isDark ? 'text-slate-400' : 'text-gray-500'}`}>加载中...</p>
        </div>
      </div>
    );
  }

  // 错误状态
  if (error || !article) {
    return (
      <div className={`min-h-screen flex flex-col items-center justify-center ${isDark ? 'bg-slate-900' : 'bg-white'}`}>
        <p className={`text-lg mb-4 ${isDark ? 'text-slate-400' : 'text-gray-500'}`}>{error || '文章不存在'}</p>
        <Link
          to="/list"
          className="flex items-center gap-2 px-4 py-2 bg-wechat text-white rounded-lg hover:bg-wechat-dark transition-colors"
        >
          <ArrowLeft size={18} />
          返回列表
        </Link>
      </div>
    );
  }

  return (
    <div className={`fixed inset-0 z-50 overflow-y-auto ${isDark ? 'dark bg-slate-900' : 'bg-white'}`}>
      {/* 目录侧边栏 */}
      {showCatalog && article.content && (
        <div className="fixed left-0 top-0 h-full w-64 bg-gray-50 dark:bg-slate-800 border-r border-gray-200 dark:border-slate-700 overflow-y-auto z-20 hidden lg:block">
          <div className="p-4 pt-20">
            <h3 className="text-sm font-bold text-gray-500 dark:text-slate-400 uppercase mb-4">目录</h3>
            <MdCatalog
              editorId={editorId}
              theme={isDark ? 'dark' : 'light'}
              className="!bg-transparent"
            />
          </div>
        </div>
      )}

      {/* Header */}
      <header
        className={`sticky top-0 z-30 ${isDark ? 'bg-slate-900/95' : 'bg-white/95'} backdrop-blur-md border-b ${isDark ? 'border-slate-800' : 'border-gray-100'} px-4 py-3 ${showCatalog ? 'lg:ml-64' : ''} transition-all duration-300`}
      >
        <div className="flex justify-between items-center w-full">
          {/* 左侧按钮组 */}
          <div className="flex items-center gap-2 shrink-0">
            <button
              onClick={handleBack}
              className={`p-2 hover:${isDark ? 'bg-slate-800' : 'bg-gray-100'} rounded-full transition-colors ${isDark ? 'text-slate-400' : 'text-gray-500'}`}
              title="返回"
            >
              <ArrowLeft size={24} />
            </button>

            <Link
              to="/"
              className={`p-2 hover:${isDark ? 'bg-slate-800' : 'bg-gray-100'} rounded-full transition-colors ${isDark ? 'text-slate-400' : 'text-gray-500'}`}
              title="首页"
            >
              <Home size={20} />
            </Link>

            {/* 目录切换 */}
            <button
              onClick={() => setShowCatalog(!showCatalog)}
              className={`p-2 rounded-lg transition-colors ${showCatalog ? 'bg-wechat text-white' : `${isDark ? 'text-slate-400 hover:bg-slate-800' : 'text-gray-500 hover:bg-gray-100'}`}`}
              title="显示/隐藏目录"
            >
              <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                <line x1="8" y1="6" x2="21" y2="6"></line>
                <line x1="8" y1="12" x2="21" y2="12"></line>
                <line x1="8" y1="18" x2="21" y2="18"></line>
                <line x1="3" y1="6" x2="3.01" y2="6"></line>
                <line x1="3" y1="12" x2="3.01" y2="12"></line>
                <line x1="3" y1="18" x2="3.01" y2="18"></line>
              </svg>
            </button>

            {/* 主题选择 */}
            <div className="relative" ref={themeDropdownRef}>
              <button
                onClick={() => setShowThemeDropdown(!showThemeDropdown)}
                className={`flex items-center gap-1.5 px-2.5 py-1.5 rounded-lg text-xs font-medium transition-colors ${
                  showThemeDropdown
                    ? 'bg-wechat text-white'
                    : isDark
                      ? 'text-slate-400 hover:bg-slate-800'
                      : 'text-gray-500 hover:bg-gray-100'
                }`}
              >
                <Palette size={14} />
                <span className="hidden sm:inline">{themeOptions.find(t => t.value === previewTheme)?.label}</span>
                <ChevronDown size={12} className={`transition-transform ${showThemeDropdown ? 'rotate-180' : ''}`} />
              </button>

              {showThemeDropdown && (
                <div className={`absolute top-full left-0 mt-1 py-1 rounded-xl shadow-lg border z-50 min-w-[130px] animate-in fade-in zoom-in-95 duration-150 ${
                  isDark
                    ? 'bg-slate-800 border-slate-700'
                    : 'bg-white border-gray-200'
                }`}>
                  {themeOptions.map((option) => (
                    <button
                      key={option.value}
                      onClick={() => handleThemeSelect(option.value as typeof previewTheme)}
                      className={`w-full px-3 py-2 text-left text-xs flex items-center justify-between transition-colors ${
                        previewTheme === option.value
                          ? 'text-wechat font-medium'
                          : isDark
                            ? 'text-slate-300 hover:bg-slate-700'
                            : 'text-gray-600 hover:bg-gray-50'
                      }`}
                    >
                      {option.label}
                      {previewTheme === option.value && <Check size={12} />}
                    </button>
                  ))}
                </div>
              )}
            </div>

            {/* 深色模式切换 */}
            <button
              onClick={() => setIsDark(!isDark)}
              className={`p-2 rounded-lg hover:${isDark ? 'bg-slate-800' : 'bg-gray-100'} transition-colors ${isDark ? 'text-slate-400' : 'text-gray-500'}`}
              title={isDark ? '切换到浅色模式' : '切换到深色模式'}
            >
              {isDark ? <Sun size={18} /> : <Moon size={18} />}
            </button>
          </div>

          {/* 右侧按钮组 */}
          <div className="flex items-center gap-2 shrink-0">
            <button
              onClick={handleEdit}
              className={`flex items-center gap-2 px-4 py-2 text-sm font-medium ${isDark ? 'text-slate-300 hover:bg-slate-800' : 'text-gray-600 hover:bg-gray-100'} rounded-full transition-colors`}
            >
              <Edit3 size={18} />
              <span className="hidden sm:inline">编辑</span>
            </button>
            {article.url && (
              <a
                href={article.url}
                target="_blank"
                rel="noopener noreferrer"
                className="flex items-center gap-2 px-4 py-2 text-sm font-medium text-wechat hover:bg-wechat-light dark:hover:bg-wechat-dark/20 rounded-full transition-colors"
              >
                <ExternalLink size={18} />
                <span className="hidden sm:inline">原文</span>
              </a>
            )}
          </div>
        </div>
      </header>

      {/* Content */}
      <main className={`max-w-3xl mx-auto px-6 py-12 ${showCatalog ? 'lg:ml-64' : ''} transition-all duration-300`}>
        <header className="mb-10 text-center">
          <h1 className={`text-3xl md:text-4xl font-bold ${isDark ? 'text-white' : 'text-gray-900'} mb-6 leading-tight`}>
            {article.title}
          </h1>

          <div className={`flex flex-wrap items-center justify-center gap-4 ${isDark ? 'text-slate-400' : 'text-gray-500'} text-sm`}>
            <div className="flex items-center gap-1.5">
              <User size={16} />
              <span className={`font-medium ${isDark ? 'text-slate-300' : 'text-slate-700'}`}>{article.author}</span>
            </div>
            <div className={`w-1 h-1 ${isDark ? 'bg-slate-600' : 'bg-gray-300'} rounded-full`}></div>
            <div className="flex items-center gap-1.5">
              <Calendar size={16} />
              <span>{new Date(article.dateAdded).toLocaleDateString()}</span>
            </div>
          </div>

          {article.summary && (
            <div className={`mt-8 p-6 ${isDark ? 'bg-slate-800 text-slate-300 border-wechat-dark/50' : 'bg-gray-50 text-gray-600 border-wechat'} rounded-2xl italic border-l-4 text-left`}>
              {article.summary}
            </div>
          )}
        </header>

        {article.content ? (
          <div className={isDark ? 'dark' : ''}>
            <MdPreview
              id={editorId}
              modelValue={processContent(article.content)}
              theme={isDark ? 'dark' : 'light'}
              previewTheme={previewTheme}
              className="!bg-transparent"
            />
          </div>
        ) : (
          <div className={`text-center ${isDark ? 'text-slate-500' : 'text-gray-400'} py-20 flex flex-col items-center`}>
            <p className="mb-4">本文暂无内容。</p>
            <button
              onClick={handleEdit}
              className="text-wechat hover:underline"
            >
              手动添加内容
            </button>
          </div>
        )}

        {/* 底部导航 */}
        <div className={`mt-16 pt-8 border-t ${isDark ? 'border-slate-800' : 'border-gray-100'}`}>
          <div className="flex justify-between items-center">
            <button
              onClick={handleBack}
              className={`flex items-center gap-2 px-4 py-2 ${isDark ? 'text-slate-400 hover:text-slate-200' : 'text-gray-500 hover:text-gray-700'} transition-colors`}
            >
              <ArrowLeft size={18} />
              返回
            </button>

            <div className="flex gap-3">
              {article.tags && article.tags.length > 0 && (
                <div className="flex flex-wrap gap-2">
                  {article.tags.map((tag, index) => (
                    <span
                      key={index}
                      className={`px-3 py-1 text-xs rounded-full ${isDark ? 'bg-slate-800 text-slate-400' : 'bg-gray-100 text-gray-600'}`}
                    >
                      {tag}
                    </span>
                  ))}
                </div>
              )}
            </div>
          </div>
        </div>
      </main>
    </div>
  );
};

export default ReaderPage;
