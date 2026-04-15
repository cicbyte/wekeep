import React, { useState, useRef, useEffect, useMemo, useCallback } from 'react';
import { Search, BookOpen, Loader2, X, ChevronLeft, ChevronRight } from 'lucide-react';
import { Article } from '../types';
import ArticleCard from './ArticleCard';
import SearchableSelect from './SearchableSelect';

interface ListPageProps {
  articles: Article[];
  authorsList: Array<{ id: number; name: string; articleCount: number; avatar?: string }>;
  loading: boolean;
  searchArticles: (keyword: string, pageNum?: number, pageSize?: number) => Promise<void>;
  loadArticles: (params?: { pageNum?: number; pageSize?: number }) => Promise<void>;
  filterByAuthor: (author: string, pageNum?: number, pageSize?: number) => Promise<void>;
  selectedAuthor: string | null;
  onAuthorFilter: (author: string) => void;
  onDeleteArticle: (id: string) => void;
  onReparseArticle?: (article: Article) => void;
  isSidebarCollapsed: boolean;
  totalCount: number;
  currentPage: number;
}

const ROWS_PER_PAGE = 4; // 每页固定4行

// 根据屏幕宽度计算每行列数
const getColumnsPerRow = (width: number): number => {
  if (width >= 1536) return 6;  // 2xl
  if (width >= 1280) return 5;  // xl
  if (width >= 1024) return 4;  // lg
  if (width >= 768) return 3;   // md
  return 2;                     // 移动端
};

const ListPage: React.FC<ListPageProps> = ({
  articles,
  authorsList,
  loading,
  searchArticles,
  loadArticles,
  filterByAuthor,
  selectedAuthor,
  onAuthorFilter,
  onDeleteArticle,
  onReparseArticle,
  isSidebarCollapsed,
  totalCount,
  currentPage
}) => {
  const [searchQuery, setSearchQuery] = useState('');
  const [jumpPage, setJumpPage] = useState('');
  const [searchType, setSearchType] = useState<'title' | 'fulltext' | null>(null); // 跟踪搜索类型
  const [searchKeyword, setSearchKeyword] = useState(''); // 跟踪实际的搜索关键词

  // 响应式分页：检测屏幕宽度
  const [pageSize, setPageSize] = useState(12); // 默认 2列 × 6行 = 12
  const prevPageSizeRef = useRef(pageSize);
  const resizeTimeoutRef = useRef<NodeJS.Timeout | null>(null);

  useEffect(() => {
    const calculatePageSize = () => {
      const width = window.innerWidth;
      const columns = getColumnsPerRow(width);
      const newPageSize = columns * ROWS_PER_PAGE;
      setPageSize(newPageSize);
    };

    // 防抖处理窗口大小变化
    const handleResize = () => {
      if (resizeTimeoutRef.current) {
        clearTimeout(resizeTimeoutRef.current);
      }
      resizeTimeoutRef.current = setTimeout(calculatePageSize, 200);
    };

    calculatePageSize();
    window.addEventListener('resize', handleResize);
    return () => {
      window.removeEventListener('resize', handleResize);
      if (resizeTimeoutRef.current) {
        clearTimeout(resizeTimeoutRef.current);
      }
    };
  }, []);

  // 当 pageSize 变化时，重新加载第一页数据
  useEffect(() => {
    if (prevPageSizeRef.current !== pageSize) {
      prevPageSizeRef.current = pageSize;
      // 重新加载当前页数据
      if (searchType === 'title' && searchKeyword) {
        loadArticles({ pageNum: 1, pageSize, keyword: searchKeyword, titleOnly: true });
      } else if (searchType === 'fulltext' && searchKeyword) {
        searchArticles(searchKeyword, 1, pageSize);
      } else if (selectedAuthor) {
        filterByAuthor(selectedAuthor, 1, pageSize);
      } else {
        loadArticles({ pageNum: 1, pageSize });
      }
    }
  }, [pageSize, searchType, searchKeyword, selectedAuthor, searchArticles, filterByAuthor, loadArticles]);

  // 计算总页数
  const totalPages = Math.ceil(totalCount / pageSize);

  // 作者选项
  const authorOptions = useMemo(() =>
    authorsList.map(author => ({
      value: author.name,
      label: author.name,
      count: author.articleCount,
      avatar: author.avatar
    })),
    [authorsList]
  );

  // 解析搜索语法：检测 title: 前缀
  const parseSearchQuery = useCallback((query: string): { type: 'title' | 'fulltext'; keyword: string } => {
    const trimmedQuery = query.trim();
    // 检测 title: 前缀（不区分大小写）
    const titleMatch = trimmedQuery.match(/^title:(.+)$/i);
    if (titleMatch) {
      return { type: 'title', keyword: titleMatch[1].trim() };
    }
    return { type: 'fulltext', keyword: trimmedQuery };
  }, []);

  // 搜索输入变化（只更新输入框内容，不执行搜索）
  const handleSearchChange = useCallback((query: string) => {
    setSearchQuery(query);
  }, []);

  // 执行搜索
  const executeSearch = useCallback((query: string) => {
    if (!query.trim()) {
      setSearchType(null);
      setSearchKeyword('');
      loadArticles({ pageNum: 1, pageSize });
      return;
    }

    const { type, keyword } = parseSearchQuery(query);
    setSearchType(type);
    setSearchKeyword(keyword);

    if (type === 'title') {
      // 使用列表接口进行标题搜索（仅搜索标题）
      loadArticles({ pageNum: 1, pageSize, keyword, titleOnly: true });
    } else {
      // 使用全文检索
      searchArticles(keyword, 1, pageSize);
    }
  }, [searchArticles, loadArticles, pageSize, parseSearchQuery]);

  // 回车搜索
  const handleSearchKeyDown = useCallback((e: React.KeyboardEvent<HTMLInputElement>) => {
    if (e.key === 'Enter') {
      executeSearch(searchQuery);
    }
  }, [executeSearch, searchQuery]);

  // 清除搜索
  const handleClearSearch = useCallback(() => {
    setSearchQuery('');
    setSearchType(null);
    setSearchKeyword('');
    loadArticles({ pageNum: 1, pageSize });
  }, [loadArticles, pageSize]);

  // 分页导航
  const handlePageChange = useCallback((page: number) => {
    if (page < 1 || page > totalPages || loading) return;

    if (searchType === 'title' && searchKeyword) {
      // 标题搜索分页（仅搜索标题）
      loadArticles({ pageNum: page, pageSize, keyword: searchKeyword, titleOnly: true });
    } else if (searchType === 'fulltext' && searchKeyword) {
      // 全文搜索分页
      searchArticles(searchKeyword, page, pageSize);
    } else if (selectedAuthor) {
      filterByAuthor(selectedAuthor, page, pageSize);
    } else {
      loadArticles({ pageNum: page, pageSize });
    }
  }, [totalPages, loading, searchType, searchKeyword, selectedAuthor, searchArticles, filterByAuthor, loadArticles, pageSize]);

  // 跳转到指定页
  const handleJumpPage = useCallback(() => {
    if (!jumpPage) return;
    const page = parseInt(jumpPage, 10);
    if (isNaN(page) || page < 1 || page > totalPages) return;
    handlePageChange(page);
    setJumpPage('');
  }, [jumpPage, totalPages, handlePageChange]);

  // 跳转输入框键盘事件
  const handleJumpKeyDown = useCallback((e: React.KeyboardEvent<HTMLInputElement>) => {
    if (e.key === 'Enter') {
      handleJumpPage();
    }
  }, [handleJumpPage]);

  // 生成页码数组
  const getPageNumbers = useCallback(() => {
    const pages: (number | string)[] = [];
    const maxVisible = 5;

    if (totalPages <= maxVisible + 2) {
      for (let i = 1; i <= totalPages; i++) {
        pages.push(i);
      }
    } else {
      pages.push(1);

      if (currentPage > 3) {
        pages.push('...');
      }

      const start = Math.max(2, currentPage - 1);
      const end = Math.min(totalPages - 1, currentPage + 1);

      for (let i = start; i <= end; i++) {
        pages.push(i);
      }

      if (currentPage < totalPages - 2) {
        pages.push('...');
      }

      if (totalPages > 1) {
        pages.push(totalPages);
      }
    }

    return pages;
  }, [totalPages, currentPage]);

  return (
    <div className="space-y-4 animate-in fade-in duration-500 pb-36 md:pb-14">
      <div className="flex flex-wrap items-center gap-3">
        {/* 文章搜索框 */}
        <div className="relative flex-1 min-w-[200px] max-w-md">
          <Search size={18} className="absolute left-3 top-1/2 -translate-y-1/2 text-gray-400 dark:text-slate-500" />
          <input
            type="text"
            placeholder="搜索文章... (回车搜索，title:关键词 仅搜索标题)"
            value={searchQuery}
            onChange={(e) => handleSearchChange(e.target.value)}
            onKeyDown={handleSearchKeyDown}
            className="w-full bg-gray-100 dark:bg-slate-800 border-transparent focus:bg-white dark:focus:bg-slate-900 focus:ring-2 focus:ring-wechat/50 focus:border-wechat rounded-xl py-2 pl-10 pr-10 text-sm transition-all outline-none dark:text-slate-200 dark:placeholder-slate-500"
          />
          {/* 清除按钮或加载动画 */}
          {searchQuery && !loading ? (
            <button
              type="button"
              onClick={handleClearSearch}
              className="absolute right-3 top-1/2 -translate-y-1/2 text-gray-400 hover:text-gray-600 dark:hover:text-slate-300 transition-colors"
            >
              <X size={16} />
            </button>
          ) : loading ? (
            <Loader2 size={16} className="absolute right-3 top-1/2 -translate-y-1/2 text-wechat animate-spin" />
          ) : null}
        </div>

        {/* 作者筛选 */}
        <SearchableSelect
          options={authorOptions}
          value={selectedAuthor || ''}
          onChange={onAuthorFilter}
          placeholder="选择作者"
          emptyText="未找到作者"
        />
      </div>

      {/* 文章列表 */}
      <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5 2xl:grid-cols-6 gap-4">
        {articles.map(article => (
          <ArticleCard
            key={article.id}
            article={article}
            onDelete={onDeleteArticle}
            onReparse={onReparseArticle}
          />
        ))}
        {articles.length === 0 && !loading && (
          <div className="col-span-full py-12 text-center text-gray-400 dark:text-slate-500 bg-white dark:bg-slate-800 rounded-xl border border-dashed border-gray-200 dark:border-slate-700">
            <BookOpen size={48} className="mx-auto mb-3 opacity-20" />
            <p>{searchQuery || selectedAuthor ? '未找到匹配的文章' : '暂无文章'}</p>
          </div>
        )}
      </div>

      {/* 底部固定分页栏 */}
      {totalCount > 0 && (
        <div
          className={`fixed bottom-24 md:bottom-0 left-0 right-0 bg-white/95 dark:bg-slate-900/95 backdrop-blur-sm border-t border-gray-200 dark:border-slate-800 py-2 px-3 md:px-4 z-10 transition-all duration-300 ${
            isSidebarCollapsed ? 'md:left-20' : 'md:left-64'
          }`}
        >
          {/* 移动端简化布局 */}
          <div className="flex md:hidden items-center justify-between">
            {/* 上一页 */}
            <button
              onClick={() => handlePageChange(currentPage - 1)}
              disabled={currentPage <= 1 || loading}
              className={`p-2 rounded-lg transition-colors ${
                currentPage <= 1 || loading
                  ? 'text-gray-300 dark:text-slate-600'
                  : 'text-gray-500 dark:text-slate-400 active:bg-gray-100 dark:active:bg-slate-800'
              }`}
            >
              <ChevronLeft size={20} />
            </button>

            {/* 页码信息 */}
            <div className="flex items-center gap-2">
              <span className="text-sm font-medium text-slate-700 dark:text-slate-300">
                {currentPage} / {totalPages}
              </span>
              <span className="text-xs text-gray-400 dark:text-slate-500">
                共{totalCount}篇
              </span>
            </div>

            {/* 下一页 */}
            <button
              onClick={() => handlePageChange(currentPage + 1)}
              disabled={currentPage >= totalPages || loading}
              className={`p-2 rounded-lg transition-colors ${
                currentPage >= totalPages || loading
                  ? 'text-gray-300 dark:text-slate-600'
                  : 'text-gray-500 dark:text-slate-400 active:bg-gray-100 dark:active:bg-slate-800'
              }`}
            >
              <ChevronRight size={20} />
            </button>
          </div>

          {/* 桌面端完整布局 */}
          <div className="hidden md:flex items-center justify-center gap-2">
            {/* 上一页 */}
            <button
              onClick={() => handlePageChange(currentPage - 1)}
              disabled={currentPage <= 1 || loading}
              className={`p-1.5 rounded-lg transition-colors ${
                currentPage <= 1 || loading
                  ? 'text-gray-300 dark:text-slate-600 cursor-not-allowed'
                  : 'text-gray-500 dark:text-slate-400 hover:bg-gray-100 dark:hover:bg-slate-800'
              }`}
            >
              <ChevronLeft size={18} />
            </button>

            {/* 页码 */}
            <div className="flex items-center gap-1">
              {getPageNumbers().map((page, index) => (
                typeof page === 'number' ? (
                  <button
                    key={index}
                    onClick={() => handlePageChange(page)}
                    disabled={loading}
                    className={`min-w-[28px] h-7 px-1.5 rounded-lg text-sm font-medium transition-colors ${
                      currentPage === page
                        ? 'bg-wechat text-white'
                        : loading
                          ? 'text-gray-300 dark:text-slate-600 cursor-not-allowed'
                          : 'text-gray-500 dark:text-slate-400 hover:bg-gray-100 dark:hover:bg-slate-800'
                    }`}
                  >
                    {page}
                  </button>
                ) : (
                  <span key={index} className="px-0.5 text-gray-400 dark:text-slate-500 text-sm">
                    {page}
                  </span>
                )
              ))}
            </div>

            {/* 下一页 */}
            <button
              onClick={() => handlePageChange(currentPage + 1)}
              disabled={currentPage >= totalPages || loading}
              className={`p-1.5 rounded-lg transition-colors ${
                currentPage >= totalPages || loading
                  ? 'text-gray-300 dark:text-slate-600 cursor-not-allowed'
                  : 'text-gray-500 dark:text-slate-400 hover:bg-gray-100 dark:hover:bg-slate-800'
              }`}
            >
              <ChevronRight size={18} />
            </button>

            {/* 分隔线 */}
            <div className="w-px h-5 bg-gray-200 dark:bg-slate-700 mx-1" />

            {/* 跳转输入框 */}
            <div className="flex items-center gap-1">
              <span className="text-xs text-gray-400 dark:text-slate-500">跳至</span>
              <input
                type="text"
                inputMode="numeric"
                pattern="[0-9]*"
                value={jumpPage}
                onChange={(e) => setJumpPage(e.target.value.replace(/\D/g, ''))}
                onKeyDown={handleJumpKeyDown}
                className="w-10 h-7 text-center text-sm bg-gray-50 dark:bg-slate-800 border border-gray-200 dark:border-slate-700 rounded-lg outline-none focus:ring-1 focus:ring-wechat/50 focus:border-wechat dark:text-slate-200"
                placeholder=""
              />
              <span className="text-xs text-gray-400 dark:text-slate-500">页</span>
              <button
                type="button"
                onClick={handleJumpPage}
                disabled={!jumpPage || loading}
                className="px-2 h-7 text-xs bg-wechat/10 text-wechat rounded-lg hover:bg-wechat/20 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
              >
                Go
              </button>
            </div>

            {/* 总数信息 */}
            <span className="ml-2 text-xs text-gray-400 dark:text-slate-500 whitespace-nowrap">
              共{totalCount}篇
            </span>
          </div>
        </div>
      )}
    </div>
  );
};

export default ListPage;
