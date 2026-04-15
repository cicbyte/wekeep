
import React, { useState, useEffect, useMemo, useCallback, useRef } from 'react';
import {
  BrowserRouter,
  Routes,
  Route,
  Link,
  useLocation,
  Navigate,
  useNavigate
} from 'react-router-dom';
import {
  LayoutGrid,
  Users,
  Settings,
  Plus,
  Search,
  BookOpen,
  ChevronRight,
  Menu,
  ChevronLeft,
  Home,
  PanelLeft,
  Database,
  Info,
  Monitor,
  Download,
  Upload,
  Trash2,
  Check,
  Moon,
  Sun,
  Loader2,
  X,
  RefreshCw
} from 'lucide-react';
import { Article, ViewMode } from './types';
import ArticleCard from './components/ArticleCard';
import AddArticleModal from './components/AddArticleModal';
import SearchableSelect from './components/SearchableSelect';
import ListPage from './components/ListPage';
import Stats from './components/Stats';
import ConfirmDialog from './components/ConfirmDialog';
import ReaderModal from './components/ReaderModal';
import ReaderPage from './components/ReaderPage';
import EditArticleModal from './components/EditArticleModal';
import EditPage from './components/EditPage';
import { useToast } from './components/Toast';
import { useArticles } from './hooks/useArticles';
import { searchApi, articlesApi, healthApi } from './services/apiService';
import { migrateToBackend, needsMigration, loadFromLocalStorage } from './services/migrationService';

// App 内容组件：包含所有路由逻辑和 UI
const AppContent: React.FC = () => {
  const {
    articles,
    loading,
    error,
    loadArticles,
    addArticle,
    updateArticle,
    deleteArticle,
    batchDeleteArticles,
    searchArticles,
    filterByAuthor,
    authors: authorsList,
    loadAuthors,
    deleteAuthor,
    totalCount,
    currentPage,
  } = useArticles();

  const { showToast } = useToast();

  const [isAddModalOpen, setIsAddModalOpen] = useState(false);
  const [selectedAuthor, setSelectedAuthor] = useState<string | null>(null);
  const selectedAuthorRef = useRef<string | null>(null);  // 用于同步追踪
  const [isSidebarCollapsed, setIsSidebarCollapsed] = useState(false);
  const navigate = useNavigate();
  const location = useLocation();

  // Settings State
  const [defaultView, setDefaultView] = useState<ViewMode>(() => {
    return (localStorage.getItem('wekeep_default_view') as ViewMode) || 'dashboard';
  });

  const [isDarkMode, setIsDarkMode] = useState(() => {
    const saved = localStorage.getItem('wekeep_dark_mode');
    if (saved !== null) return saved === 'true';
    return window.matchMedia('(prefers-color-scheme: dark)').matches;
  });

  // Dialog States
  const [articleToDelete, setArticleToDelete] = useState<string | null>(null);
  const [readingArticle, setReadingArticle] = useState<Article | null>(null);
  const [editingArticle, setEditingArticle] = useState<Article | null>(null);

  // Context Menu State for Authors
  const [contextMenu, setContextMenu] = useState<{
    visible: boolean;
    x: number;
    y: number;
    author: { id: number; name: string; articleCount: number } | null;
  }>({ visible: false, x: 0, y: 0, author: null });

  // Migration State
  const [isMigrating, setIsMigrating] = useState(false);
  const [migrationStatus, setMigrationStatus] = useState<string>('');

  // Settings State
  const [settingsTab, setSettingsTab] = useState<'general' | 'data' | 'search' | 'storage' | 'about'>(() => {
    const hash = location.hash.slice(1);
    if (hash === 'data' || hash === 'search' || hash === 'about') return hash;
    if (hash === 'storage') return 'storage';
    return 'general';
  });

  // Search Engine State
  const [searchEnabled, setSearchEnabled] = useState<boolean | null>(null); // null = 未检测
  const [storageMigrationEnabled, setStorageMigrationEnabled] = useState(false);
  const [indexedCount, setIndexedCount] = useState<number>(0);
  const [isIndexing, setIsIndexing] = useState(false);
  const [searchStatusLoading, setSearchStatusLoading] = useState(true);
  const [appVersion, setAppVersion] = useState('');

  // 初始化应用 - 自动迁移数据
  useEffect(() => {
    const initApp = async () => {
      // 检查是否需要迁移 localStorage 数据
      if (needsMigration()) {
        setIsMigrating(true);
        setMigrationStatus('正在迁移本地数据到服务器...');

        try {
          const result = await migrateToBackend();
          setMigrationStatus(`迁移完成！成功 ${result.success} 篇，失败 ${result.failed} 篇`);

          // 迁移后刷新文章列表
          setTimeout(() => {
            loadArticles();
            setIsMigrating(false);
            setMigrationStatus('');
          }, 2000);
        } catch (error) {
          console.error('迁移失败:', error);
          setMigrationStatus('迁移失败，请检查网络连接');
          setTimeout(() => setIsMigrating(false), 3000);
        }
      }
    };

    initApp();
  }, [loadArticles]);

  // 加载搜索引擎状态
  useEffect(() => {
    const loadSearchStatus = async () => {
      // 使用 location.pathname 代替 viewMode（避免变量声明顺序问题）
      if (!location.pathname.startsWith('/settings')) return;

      setSearchStatusLoading(true);
      try {
        const status = await searchApi.status();
        setSearchEnabled(status.enabled);
        setIndexedCount(status.indexedCount);
      } catch (error) {
        console.error('获取搜索状态失败:', error);
        setSearchEnabled(false);
      } finally {
        setSearchStatusLoading(false);
        // 并行加载存储迁移开关状态
        fetch(`${import.meta.env.VITE_API_BASE_URL || '/api/v1'}/storage/status`)
          .then(r => r.json())
          .then(res => {
            if (res.code === 0 && res.data) {
              setStorageMigrationEnabled(!!res.data.migrationEnabled);
            }
          })
          .catch(() => {});
        // 加载版本号
        healthApi.version()
          .then(res => setAppVersion(res.version || ''))
          .catch(() => {});
      }
    };

    loadSearchStatus();
  }, [location.pathname]);

  // 深色模式应用
  useEffect(() => {
    localStorage.setItem('wekeep_dark_mode', String(isDarkMode));
    if (isDarkMode) {
      document.documentElement.classList.add('dark');
    } else {
      document.documentElement.classList.remove('dark');
    }
  }, [isDarkMode]);

  useEffect(() => {
    localStorage.setItem('wekeep_dark_mode', String(isDarkMode));
    if (isDarkMode) {
      document.documentElement.classList.add('dark');
    } else {
      document.documentElement.classList.remove('dark');
    }
  }, [isDarkMode]);

  const handleSaveDefaultView = (mode: ViewMode) => {
    setDefaultView(mode);
    localStorage.setItem('wekeep_default_view', mode);
    navigate(`/${mode}`);
  };

  const handleExportData = () => {
    const dataStr = JSON.stringify(articles, null, 2);
    const blob = new Blob([dataStr], { type: "application/json" });
    const url = URL.createObjectURL(blob);
    const link = document.createElement('a');
    link.href = url;
    link.download = `wekeep_backup_${new Date().toISOString().split('T')[0]}.json`;
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
  };

  // 同步文章到搜索引擎
  const handleIndexAllArticles = async () => {
    setIsIndexing(true);
    try {
      const result = await searchApi.indexAll();
      setIndexedCount(result.indexedCount);
      showToast(`索引完成，成功同步 ${result.indexedCount} 篇文章`, 'success');
    } catch (error: any) {
      console.error('索引失败:', error);
      showToast(error?.message || '索引失败', 'error');
    } finally {
      setIsIndexing(false);
    }
  };

  const handleImportData = async (event: React.ChangeEvent<HTMLInputElement>) => {
    const file = event.target.files?.[0];
    if (!file) return;

    try {
      const text = await file.text();
      const importedArticles = JSON.parse(text) as Article[];

      if (!Array.isArray(importedArticles)) {
        showToast('无效的备份文件格式', 'error');
        return;
      }

      let successCount = 0;
      let failCount = 0;

      for (const article of importedArticles) {
        const success = await addArticle({
          ...article,
          id: '', // 让后端生成新 ID
        });
        if (success) {
          successCount++;
        } else {
          failCount++;
        }
      }

      if (successCount > 0) {
        showToast(`导入成功 ${successCount} 篇文章${failCount > 0 ? `，失败 ${failCount} 篇` : ''}`, 'success');
        loadArticles();
      } else {
        showToast('导入失败，请检查文件格式', 'error');
      }
    } catch (error) {
      showToast('文件解析失败', 'error');
      console.error('Import error:', error);
    }

    // 清空文件输入
    event.target.value = '';
  };

  const handleAddArticle = async (article: Article) => {
    await addArticle(article);
  };

  const handleDeleteArticle = (id: string) => {
    setArticleToDelete(id);
  };

  // 重新解析文章
  const handleReparseArticle = async (article: Article) => {
    try {
      const result = await articlesApi.reparse(parseInt(article.id));
      showToast(`重新解析成功，处理了 ${result.imagesCount} 张图片`, 'success');
      // 刷新文章列表
      loadArticles();
    } catch (err: any) {
      showToast(err.message || '重新解析失败', 'error');
    }
  };

  // Author Context Menu Handlers
  const handleAuthorContextMenu = (e: React.MouseEvent, author: { id: number; name: string; articleCount: number }) => {
    e.preventDefault();
    e.stopPropagation();
    setContextMenu({
      visible: true,
      x: e.clientX,
      y: e.clientY,
      author
    });
  };

  const handleCloseContextMenu = () => {
    setContextMenu({ visible: false, x: 0, y: 0, author: null });
  };

  const handleDeleteAuthor = async () => {
    if (!contextMenu.author) return;

    const { id, name, articleCount } = contextMenu.author;

    // 只有文章数为0的作者可以删除
    if (articleCount > 0) {
      showToast('该作者还有文章，无法删除', 'error');
      handleCloseContextMenu();
      return;
    }

    handleCloseContextMenu();

    const success = await deleteAuthor(id);
    if (success) {
      showToast(`作者 "${name}" 已删除`, 'success');
    } else {
      showToast('删除作者失败', 'error');
    }
  };

  // 点击其他地方关闭右键菜单
  useEffect(() => {
    const handleClick = () => handleCloseContextMenu();
    const handleEscape = (e: KeyboardEvent) => {
      if (e.key === 'Escape') handleCloseContextMenu();
    };

    if (contextMenu.visible) {
      document.addEventListener('click', handleClick);
      document.addEventListener('keydown', handleEscape);
      return () => {
        document.removeEventListener('click', handleClick);
        document.removeEventListener('keydown', handleEscape);
      };
    }
  }, [contextMenu.visible]);

  const handleUpdateArticle = async (updatedArticle: Article) => {
    try {
      const success = await updateArticle(updatedArticle);
      if (success) {
        if (readingArticle?.id === updatedArticle.id) {
          setReadingArticle(updatedArticle);
        }
        setEditingArticle(null);
        showToast('文章更新成功！', 'success');
      } else {
        showToast('文章更新失败，请重试', 'error');
      }
    } catch (error) {
      console.error('更新文章错误:', error);
      showToast('文章更新失败，请重试', 'error');
    }
  };

  const confirmDeleteArticle = async () => {
    if (articleToDelete) {
      try {
        const success = await deleteArticle(articleToDelete);
        if (success) {
          showToast('文章已删除', 'success');
        } else {
          showToast('删除失败，请重试', 'error');
        }
      } catch (error) {
        console.error('删除文章错误:', error);
        showToast('删除失败，请重试', 'error');
      }
      setArticleToDelete(null);
    }
  };

  // 按作者筛选
  const handleAuthorFilter = useCallback(async (author: string) => {
    setSelectedAuthor(author);
    selectedAuthorRef.current = author || null;  // 同步更新 ref
    if (author) {
      try {
        await filterByAuthor(author);
      } catch (error) {
        console.error('筛选失败:', error);
      }
    } else {
      await loadArticles();
    }
  }, [filterByAuthor, loadArticles]);

  // 获取当前路由对应的 viewMode
  const viewMode = useMemo<ViewMode>(() => {
    const path = location.pathname;
    if (path.startsWith('/authors')) return 'authors';
    if (path.startsWith('/list')) return 'list';
    if (path.startsWith('/settings')) return 'settings';
    return 'dashboard';
  }, [location.pathname]);

  // 路由变化时重新加载数据
  useEffect(() => {
    const currentAuthor = selectedAuthorRef.current;
    // 进入列表页
    if (location.pathname === '/list') {
      if (currentAuthor) {
        // 有选中的作者，由 handleAuthorClick 处理筛选
      } else if (!selectedAuthor) {
        // 没有选中作者，加载全部
        loadArticles();
      }
    }
    // 离开列表页和作者页时清除作者筛选
    if (!location.pathname.startsWith('/list') && !location.pathname.startsWith('/authors')) {
      setSelectedAuthor(null);
      selectedAuthorRef.current = null;
      loadArticles();
    }
  }, [location.pathname]);

  // 现在搜索在后端完成，直接使用 articles
  // 但对于作者视图，仍然需要按作者分组
  const articlesByAuthor = useMemo(() => {
    const groups: Record<string, Article[]> = {};
    articles.forEach(a => {
      if (!groups[a.author]) groups[a.author] = [];
      groups[a.author].push(a);
    });
    return groups;
  }, [articles]);

  // 作者名称列表（用于某些场景）
  const authorNames = Object.keys(articlesByAuthor).sort();

  // 当前显示的文章列表
  const displayedArticles = useMemo(() => {
    if (viewMode === 'authors' && selectedAuthor) {
      return articlesByAuthor[selectedAuthor] || [];
    }
    return articles;
  }, [articles, articlesByAuthor, selectedAuthor, viewMode]);

  const renderBreadcrumbs = () => {
    let label = '';
    let icon = <Home size={14} />;

    switch (viewMode) {
      case 'dashboard':
        label = '首页';
        break;
      case 'authors':
        label = '作者';
        break;
      case 'list':
        label = '文章';
        icon = <BookOpen size={14} />;
        break;
      case 'settings':
        label = '设置';
        icon = <Settings size={14} />;
        break;
      default:
        label = '首页';
    }

    return (
      <nav className="flex items-center text-sm text-gray-500 dark:text-slate-400">
        <span className="flex items-center gap-1 font-semibold text-slate-800 dark:text-slate-100">
          {icon}
          {label}
        </span>
      </nav>
    );
  };

  // 存储设置组件
  const StorageSettings = () => {
    const [storageStatus, setStorageStatus] = useState<any>(null);
    const [storageStats, setStorageStats] = useState<any>(null);
    const [sourceStorage, setSourceStorage] = useState<'local' | 'rustfs'>('rustfs');
    const [targetStorage, setTargetStorage] = useState<'local' | 'rustfs'>('local');
    const [rustfsConfig, setRustfsConfig] = useState({
      endpoint: '',
      bucket: '',
      username: '',
      password: '',
      timeout: 60,
    });
    const [localConfig, setLocalConfig] = useState({
      basePath: './uploads',
      baseURL: '/uploads',
    });
    const [updateMarkdown, setUpdateMarkdown] = useState(true);
    const [isValidating, setIsValidating] = useState(false);
    const [isMigrating, setIsMigrating] = useState(false);
    const [isUpdatingRefs, setIsUpdatingRefs] = useState(false);
    const [validationResult, setValidationResult] = useState<{valid: boolean; error?: string} | null>(null);
    const [migrationProgress, setMigrationProgress] = useState<any>(null);
    const [migrationInterval, setMigrationInterval] = useState<NodeJS.Timeout | null>(null);

    // API 基础 URL
    const API_BASE = import.meta.env.VITE_API_BASE_URL || '/api/v1';

    // 加载存储状态
    const loadStorageStatus = useCallback(async () => {
      try {
        const [statusRes, statsRes] = await Promise.all([
          fetch(`${API_BASE}/storage/status`).then(r => r.json()),
          fetch(`${API_BASE}/storage/stats`).then(r => r.json()),
        ]);
        if (statusRes.code === 0) {
          setStorageStatus(statusRes.data);
          // 如果有当前配置，预填表单
          if (statusRes.data.currentConfig?.rustfs) {
            const cfg = statusRes.data.currentConfig.rustfs;
            setRustfsConfig(prev => ({
              ...prev,
              endpoint: cfg.endpoint || prev.endpoint,
              bucket: cfg.bucket || prev.bucket,
              username: cfg.username || prev.username,
              timeout: cfg.timeout || prev.timeout,
            }));
          }
          if (statusRes.data.currentConfig?.local) {
            const cfg = statusRes.data.currentConfig.local;
            setLocalConfig(prev => ({
              ...prev,
              basePath: cfg.basePath || prev.basePath,
              baseURL: cfg.baseURL || prev.baseURL,
            }));
          }
        }
        if (statsRes.code === 0) {
          setStorageStats(statsRes.data);
          // 自动检测源存储类型（基于图片分布）
          const byStorageType = statsRes.data.byStorageType || {};
          const rustfsCount = byStorageType.rustfs || 0;
          const localCount = byStorageType.local || 0;
          if (rustfsCount > localCount) {
            setSourceStorage('rustfs');
            setTargetStorage('local');
          } else if (localCount > 0) {
            setSourceStorage('local');
            setTargetStorage('rustfs');
          }
        }
      } catch (err) {
        console.error('Failed to load storage status:', err);
      }
    }, []);

    // 初始加载
    useEffect(() => {
      loadStorageStatus();
    }, [loadStorageStatus]);

    // 轮询迁移进度
    useEffect(() => {
      if (migrationProgress?.running) {
        const interval = setInterval(async () => {
          try {
            const res = await fetch(`${API_BASE}/storage/status`).then(r => r.json());
            if (res.code === 0) {
              setMigrationProgress(res.data.migration);
              if (!res.data.migration.running) {
                setIsMigrating(false);
                loadStorageStatus();
              }
            }
          } catch (err) {
            console.error('Failed to poll migration status:', err);
          }
        }, 1000);
        setMigrationInterval(interval);
        return () => clearInterval(interval);
      } else if (migrationInterval) {
        clearInterval(migrationInterval);
        setMigrationInterval(null);
      }
    }, [migrationProgress?.running, loadStorageStatus]);

    // 验证配置
    const handleValidate = async () => {
      setIsValidating(true);
      setValidationResult(null);
      try {
        // 验证源存储（如果不是本地存储且需要配置）
        if (sourceStorage === 'rustfs') {
          const sourceConfig: any = { targetStorage: 'rustfs', rustfsConfig };
          const sourceRes = await fetch(`${API_BASE}/storage/validate`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(sourceConfig),
          }).then(r => r.json());
          if (!sourceRes.data?.valid) {
            setValidationResult({ valid: false, error: `源存储验证失败: ${sourceRes.data?.error || '未知错误'}` });
            setIsValidating(false);
            return;
          }
        }

        // 验证目标存储
        const targetConfig: any = { targetStorage };
        if (targetStorage === 'local') {
          targetConfig.localConfig = localConfig;
        } else {
          targetConfig.rustfsConfig = rustfsConfig;
        }
        const targetRes = await fetch(`${API_BASE}/storage/validate`, {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify(targetConfig),
        }).then(r => r.json());
        setValidationResult(targetRes.data);
      } catch (err) {
        setValidationResult({ valid: false, error: String(err) });
      }
      setIsValidating(false);
    };

    // 执行迁移
    const handleMigrate = async () => {
      if (!validationResult?.valid) {
        alert('请先验证配置');
        return;
      }
      setIsMigrating(true);
      try {
        const config: any = { sourceStorage, targetStorage, updateMarkdown };
        // 根据源和目标存储类型添加配置
        if (sourceStorage === 'rustfs' || targetStorage === 'rustfs') {
          config.rustfsConfig = rustfsConfig;
        }
        if (sourceStorage === 'local' || targetStorage === 'local') {
          config.localConfig = localConfig;
        }
        const res = await fetch(`${API_BASE}/storage/migrate`, {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify(config),
        }).then(r => r.json());
        if (res.code !== 0) {
          alert(res.message);
          setIsMigrating(false);
        } else {
          // 开始轮询进度
          setMigrationProgress({ running: true });
        }
      } catch (err) {
        alert('迁移失败: ' + err);
        setIsMigrating(false);
      }
    };

    // 更新引用
    const handleUpdateRefs = async () => {
      setIsUpdatingRefs(true);
      try {
        const res = await fetch(`${API_BASE}/storage/update-refs`, {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
        }).then(r => r.json());
        if (res.code === 0) {
          alert(`更新完成！共 ${res.data.total} 个图片映射，更新了 ${res.data.updated} 篇文章`);
        } else {
          alert('更新失败: ' + res.message);
        }
      } catch (err) {
        alert('更新失败: ' + err);
      }
      setIsUpdatingRefs(false);
    };

    return (
      <div className="space-y-6 animate-in slide-in-from-right-2 duration-300">
        <div>
          <h3 className="text-xl font-bold text-slate-800 dark:text-slate-100 border-b border-gray-100 dark:border-slate-800 pb-4 mb-4">存储管理</h3>

          {/* 当前存储状态 */}
          <div className="mb-6 p-4 bg-gray-50 dark:bg-slate-800 rounded-lg">
            <div className="flex items-center justify-between mb-2">
              <span className="text-sm text-gray-500 dark:text-slate-400">当前存储类型</span>
              <span className="font-medium text-slate-800 dark:text-slate-100">
                {storageStatus?.currentStorage === 'local' ? '本地存储' : 'RustFS 对象存储'}
              </span>
            </div>
            <div className="flex items-center justify-between mb-2">
              <span className="text-sm text-gray-500 dark:text-slate-400">状态</span>
              <span className={`font-medium ${storageStatus?.initialized ? 'text-green-500' : 'text-red-500'}`}>
                {storageStatus?.initialized ? '已初始化' : '未初始化'}
              </span>
            </div>
            {storageStats && (
              <div className="flex items-center justify-between">
                <span className="text-sm text-gray-500 dark:text-slate-400">已存储图片</span>
                <span className="font-medium text-slate-800 dark:text-slate-100">{storageStats.totalImages} 张</span>
              </div>
            )}
          </div>

          {/* 迁移进度 */}
          {(migrationProgress?.running || migrationProgress?.completed > 0) && (
            <div className="mb-6 p-4 bg-blue-50 dark:bg-blue-900/20 rounded-lg border border-blue-200 dark:border-blue-800">
              <div className="flex items-center gap-2 mb-3">
                {migrationProgress.running && <Loader2 size={18} className="animate-spin text-blue-500" />}
                <span className="font-medium text-blue-700 dark:text-blue-300">
                  {migrationProgress.running ? '迁移进行中...' : '迁移完成'}
                </span>
              </div>
              <div className="space-y-2 text-sm">
                <div className="flex justify-between">
                  <span className="text-blue-600 dark:text-blue-400">进度</span>
                  <span className="text-blue-700 dark:text-blue-300">
                    {migrationProgress.completed} / {migrationProgress.total}
                    {migrationProgress.failed > 0 && ` (失败: ${migrationProgress.failed})`}
                  </span>
                </div>
                {migrationProgress.currentFile && (
                  <div className="text-xs text-blue-500 dark:text-blue-400 truncate">
                    当前: {migrationProgress.currentFile}
                  </div>
                )}
                {migrationProgress.error && (
                  <div className="text-xs text-red-500">错误: {migrationProgress.error}</div>
                )}
              </div>
              <div className="mt-3 h-2 bg-blue-200 dark:bg-blue-800 rounded-full overflow-hidden">
                <div
                  className="h-full bg-blue-500 transition-all duration-300"
                  style={{ width: `${migrationProgress.total > 0 ? (migrationProgress.completed / migrationProgress.total * 100) : 0}%` }}
                />
              </div>
            </div>
          )}

          {/* 存储迁移配置 */}
          {storageStatus?.migrationEnabled && (
          <div className="space-y-4">
            <h4 className="font-semibold text-slate-800 dark:text-slate-200">存储迁移</h4>

            {/* 源存储选择 */}
            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-slate-300 mb-2">
                源存储（从哪里迁移）
              </label>
              <div className="flex gap-3">
                <button
                  onClick={() => { setSourceStorage('rustfs'); if (targetStorage === 'rustfs') setTargetStorage('local'); }}
                  disabled={isMigrating}
                  className={`flex-1 p-3 rounded-lg border-2 transition-all ${
                    sourceStorage === 'rustfs'
                      ? 'border-blue-500 bg-blue-50 dark:bg-blue-900/20'
                      : 'border-gray-200 dark:border-slate-700 hover:border-gray-300'
                  }`}
                >
                  <Database size={20} className={sourceStorage === 'rustfs' ? 'text-blue-500' : 'text-gray-400'} />
                  <div className="font-medium mt-1 text-sm">RustFS</div>
                  <div className="text-xs text-gray-500">{storageStats?.byStorageType?.rustfs || 0} 张图片</div>
                </button>
                <button
                  onClick={() => { setSourceStorage('local'); if (targetStorage === 'local') setTargetStorage('rustfs'); }}
                  disabled={isMigrating}
                  className={`flex-1 p-3 rounded-lg border-2 transition-all ${
                    sourceStorage === 'local'
                      ? 'border-blue-500 bg-blue-50 dark:bg-blue-900/20'
                      : 'border-gray-200 dark:border-slate-700 hover:border-gray-300'
                  }`}
                >
                  <Monitor size={20} className={sourceStorage === 'local' ? 'text-blue-500' : 'text-gray-400'} />
                  <div className="font-medium mt-1 text-sm">本地存储</div>
                  <div className="text-xs text-gray-500">{storageStats?.byStorageType?.local || 0} 张图片</div>
                </button>
              </div>
            </div>

            {/* 目标存储选择 */}
            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-slate-300 mb-2">
                目标存储（迁移到哪里）
              </label>
              <div className="flex gap-3">
                <button
                  onClick={() => setTargetStorage('rustfs')}
                  disabled={sourceStorage === 'rustfs' || isMigrating}
                  className={`flex-1 p-3 rounded-lg border-2 transition-all ${
                    targetStorage === 'rustfs'
                      ? 'border-wechat bg-wechat/5'
                      : 'border-gray-200 dark:border-slate-700 hover:border-gray-300'
                  } ${sourceStorage === 'rustfs' ? 'opacity-50 cursor-not-allowed' : ''}`}
                >
                  <Database size={20} className={targetStorage === 'rustfs' ? 'text-wechat' : 'text-gray-400'} />
                  <div className="font-medium mt-1 text-sm">RustFS</div>
                  <div className="text-xs text-gray-500">对象存储</div>
                </button>
                <button
                  onClick={() => setTargetStorage('local')}
                  disabled={sourceStorage === 'local' || isMigrating}
                  className={`flex-1 p-3 rounded-lg border-2 transition-all ${
                    targetStorage === 'local'
                      ? 'border-wechat bg-wechat/5'
                      : 'border-gray-200 dark:border-slate-700 hover:border-gray-300'
                  } ${sourceStorage === 'local' ? 'opacity-50 cursor-not-allowed' : ''}`}
                >
                  <Monitor size={20} className={targetStorage === 'local' ? 'text-wechat' : 'text-gray-400'} />
                  <div className="font-medium mt-1 text-sm">本地存储</div>
                  <div className="text-xs text-gray-500">零配置</div>
                </button>
              </div>
            </div>

            {/* RustFS 配置表单 */}
            {(sourceStorage === 'rustfs' || targetStorage === 'rustfs') && (
              <div className="space-y-3 p-4 bg-gray-50 dark:bg-slate-800 rounded-lg">
                <div>
                  <label className="block text-sm font-medium text-gray-700 dark:text-slate-300 mb-1">Endpoint</label>
                  <input
                    type="text"
                    value={rustfsConfig.endpoint}
                    onChange={(e) => setRustfsConfig({ ...rustfsConfig, endpoint: e.target.value })}
                    placeholder="http://localhost:9000"
                    className="w-full px-3 py-2 border border-gray-300 dark:border-slate-600 rounded-lg bg-white dark:bg-slate-700 text-slate-800 dark:text-slate-100"
                  />
                </div>
                <div className="grid grid-cols-2 gap-3">
                  <div>
                    <label className="block text-sm font-medium text-gray-700 dark:text-slate-300 mb-1">Bucket</label>
                    <input
                      type="text"
                      value={rustfsConfig.bucket}
                      onChange={(e) => setRustfsConfig({ ...rustfsConfig, bucket: e.target.value })}
                      placeholder="wekeep"
                      className="w-full px-3 py-2 border border-gray-300 dark:border-slate-600 rounded-lg bg-white dark:bg-slate-700 text-slate-800 dark:text-slate-100"
                    />
                  </div>
                  <div>
                    <label className="block text-sm font-medium text-gray-700 dark:text-slate-300 mb-1">Timeout</label>
                    <input
                      type="number"
                      value={rustfsConfig.timeout}
                      onChange={(e) => setRustfsConfig({ ...rustfsConfig, timeout: parseInt(e.target.value) || 60 })}
                      className="w-full px-3 py-2 border border-gray-300 dark:border-slate-600 rounded-lg bg-white dark:bg-slate-700 text-slate-800 dark:text-slate-100"
                    />
                  </div>
                </div>
                <div className="grid grid-cols-2 gap-3">
                  <div>
                    <label className="block text-sm font-medium text-gray-700 dark:text-slate-300 mb-1">用户名</label>
                    <input
                      type="text"
                      value={rustfsConfig.username}
                      onChange={(e) => setRustfsConfig({ ...rustfsConfig, username: e.target.value })}
                      className="w-full px-3 py-2 border border-gray-300 dark:border-slate-600 rounded-lg bg-white dark:bg-slate-700 text-slate-800 dark:text-slate-100"
                    />
                  </div>
                  <div>
                    <label className="block text-sm font-medium text-gray-700 dark:text-slate-300 mb-1">密码</label>
                    <input
                      type="password"
                      value={rustfsConfig.password}
                      onChange={(e) => setRustfsConfig({ ...rustfsConfig, password: e.target.value })}
                      className="w-full px-3 py-2 border border-gray-300 dark:border-slate-600 rounded-lg bg-white dark:bg-slate-700 text-slate-800 dark:text-slate-100"
                    />
                  </div>
                </div>
              </div>
            )}

            {/* 本地存储配置表单 */}
            {(sourceStorage === 'local' || targetStorage === 'local') && (
              <div className="space-y-3 p-4 bg-gray-50 dark:bg-slate-800 rounded-lg">
                <div>
                  <label className="block text-sm font-medium text-gray-700 dark:text-slate-300 mb-1">存储目录</label>
                  <input
                    type="text"
                    value={localConfig.basePath}
                    onChange={(e) => setLocalConfig({ ...localConfig, basePath: e.target.value })}
                    placeholder="./uploads"
                    className="w-full px-3 py-2 border border-gray-300 dark:border-slate-600 rounded-lg bg-white dark:bg-slate-700 text-slate-800 dark:text-slate-100"
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700 dark:text-slate-300 mb-1">访问URL前缀</label>
                  <input
                    type="text"
                    value={localConfig.baseURL}
                    onChange={(e) => setLocalConfig({ ...localConfig, baseURL: e.target.value })}
                    placeholder="/uploads"
                    className="w-full px-3 py-2 border border-gray-300 dark:border-slate-600 rounded-lg bg-white dark:bg-slate-700 text-slate-800 dark:text-slate-100"
                  />
                </div>
              </div>
            )}

            {/* 选项 */}
            <div className="flex items-center gap-2">
              <input
                type="checkbox"
                id="updateMarkdown"
                checked={updateMarkdown}
                onChange={(e) => setUpdateMarkdown(e.target.checked)}
                className="rounded border-gray-300"
              />
              <label htmlFor="updateMarkdown" className="text-sm text-gray-600 dark:text-slate-400">
                迁移后更新文章 Markdown 中的图片引用
              </label>
            </div>

            {/* 验证结果 */}
            {validationResult && (
              <div className={`p-3 rounded-lg ${validationResult.valid ? 'bg-green-50 dark:bg-green-900/20 text-green-700 dark:text-green-300' : 'bg-red-50 dark:bg-red-900/20 text-red-700 dark:text-red-300'}`}>
                {validationResult.valid ? '✓ 配置验证通过' : `✗ ${validationResult.error}`}
              </div>
            )}

            {/* 操作按钮 */}
            <div className="flex gap-3 pt-2">
              <button
                onClick={handleValidate}
                disabled={isValidating || isMigrating}
                className="px-4 py-2 bg-gray-100 dark:bg-slate-700 text-slate-700 dark:text-slate-200 rounded-lg hover:bg-gray-200 dark:hover:bg-slate-600 transition-colors disabled:opacity-50"
              >
                {isValidating ? <Loader2 size={16} className="animate-spin" /> : '验证配置'}
              </button>
              <button
                onClick={handleMigrate}
                disabled={!validationResult?.valid || isMigrating}
                className="px-4 py-2 bg-wechat text-white rounded-lg hover:bg-wechat-dark transition-colors disabled:opacity-50"
              >
                {isMigrating ? <Loader2 size={16} className="animate-spin" /> : '开始迁移'}
              </button>
              <button
                onClick={handleUpdateRefs}
                disabled={isUpdatingRefs || isMigrating}
                className="px-4 py-2 bg-blue-500 text-white rounded-lg hover:bg-blue-600 transition-colors disabled:opacity-50"
              >
                {isUpdatingRefs ? <Loader2 size={16} className="animate-spin" /> : '更新引用'}
              </button>
            </div>

            {/* 说明 */}
            <div className="mt-6 p-4 bg-yellow-50 dark:bg-yellow-900/20 rounded-lg border border-yellow-200 dark:border-yellow-800">
              <p className="text-sm text-yellow-700 dark:text-yellow-300 font-medium mb-2">⚠️ 迁移说明</p>
              <ul className="text-xs text-yellow-600 dark:text-yellow-400 space-y-1">
                <li>• 迁移会将所有图片从当前存储复制到目标存储</li>
                <li>• 原存储的图片不会被删除，需要手动清理</li>
                <li>• 迁移完成后需手动更新配置文件并重启服务</li>
                <li>• 建议在迁移前备份数据</li>
              </ul>
            </div>
          </div>
          )}
        </div>
      </div>
    );
  };

  // Dashboard 页面组件
  const DashboardPage = () => (
    <div className="space-y-8 animate-in fade-in duration-500">
      <div>
        <h2 className="text-2xl font-bold text-slate-800 dark:text-slate-100 mb-2">欢迎回来</h2>
        <p className="text-slate-500 dark:text-slate-400 mb-6">共 {totalCount} 篇文章，来自 {authorsList.length} 位作者。</p>
        <Stats authors={authorsList} isDarkMode={isDarkMode} />
      </div>

      <div>
         <div className="flex justify-between items-end mb-4">
          <h3 className="text-lg font-bold text-slate-800 dark:text-slate-100">最近文章</h3>
          <Link
            to="/list"
            onClick={() => { setSelectedAuthor(null); selectedAuthorRef.current = null; }}
            className="text-sm text-wechat font-medium hover:underline flex items-center"
          >
            查看全部 <ChevronRight size={16} />
          </Link>
        </div>

        <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5 gap-4">
          {displayedArticles.slice(0, 10).map(article => (
            <ArticleCard
              key={article.id}
              article={article}
              onDelete={handleDeleteArticle}
              onRead={setReadingArticle}
              onEdit={setEditingArticle}
            />
          ))}
          {displayedArticles.length === 0 && !loading && (
            <div className="col-span-full py-12 text-center text-gray-400 dark:text-slate-500 bg-white dark:bg-slate-800 rounded-xl border border-dashed border-gray-200 dark:border-slate-700">
              <BookOpen size={48} className="mx-auto mb-3 opacity-20" />
              <p>未找到文章,试试添加一篇吧!</p>
            </div>
          )}
        </div>
      </div>
    </div>
  );

  // Authors 页面 - handleAuthorClick
  const handleAuthorClick = useCallback(async (authorName: string) => {
    setSelectedAuthor(authorName);
    selectedAuthorRef.current = authorName;  // 同步设置 ref
    await filterByAuthor(authorName);  // 等待筛选完成
    navigate('/list');
  }, [filterByAuthor, navigate]);

  // Authors 页面 - 响应式分页
  const [authorSearchQuery, setAuthorSearchQuery] = useState('');
  const [authorPage, setAuthorPage] = useState(1);
  const [authorPageSize, setAuthorPageSize] = useState(8); // 默认 2列 × 4行
  const [authorJumpPage, setAuthorJumpPage] = useState('');
  const authorPageSizeRef = useRef(authorPageSize);

  // 根据屏幕宽度计算每行列数（作者卡片网格）
  const getAuthorColumnsPerRow = (width: number): number => {
    if (width >= 1280) return 6;  // xl
    if (width >= 1024) return 5;  // lg
    if (width >= 768) return 4;   // md
    if (width >= 640) return 3;   // sm
    return 2;                     // 移动端
  };

  // 监听屏幕大小变化
  useEffect(() => {
    const calculatePageSize = () => {
      const width = window.innerWidth;
      const columns = getAuthorColumnsPerRow(width);
      const newPageSize = columns * 4; // 每页4行
      setAuthorPageSize(newPageSize);
    };

    let resizeTimeout: NodeJS.Timeout;
    const handleResize = () => {
      clearTimeout(resizeTimeout);
      resizeTimeout = setTimeout(calculatePageSize, 200);
    };

    calculatePageSize();
    window.addEventListener('resize', handleResize);
    return () => {
      window.removeEventListener('resize', handleResize);
      clearTimeout(resizeTimeout);
    };
  }, []);

  // pageSize 变化时重置到第一页
  useEffect(() => {
    if (authorPageSizeRef.current !== authorPageSize) {
      authorPageSizeRef.current = authorPageSize;
      setAuthorPage(1);
    }
  }, [authorPageSize]);

  const filteredAuthors = useMemo(() => {
    if (!authorSearchQuery.trim()) return authorsList;
    const query = authorSearchQuery.toLowerCase();
    return authorsList.filter(author => author.name.toLowerCase().includes(query));
  }, [authorsList, authorSearchQuery]);

  // 搜索时重置页码
  useEffect(() => {
    setAuthorPage(1);
  }, [authorSearchQuery]);

  // 分页后的作者列表
  const paginatedAuthors = useMemo(() => {
    const start = (authorPage - 1) * authorPageSize;
    return filteredAuthors.slice(start, start + authorPageSize);
  }, [filteredAuthors, authorPage, authorPageSize]);

  const authorTotalPages = Math.ceil(filteredAuthors.length / authorPageSize) || 1;

  const handleAuthorPageChange = useCallback((page: number) => {
    if (page < 1 || page > authorTotalPages) return;
    setAuthorPage(page);
  }, [authorTotalPages]);

  // 作者页面跳转
  const handleAuthorJumpPage = useCallback(() => {
    if (!authorJumpPage) return;
    const page = parseInt(authorJumpPage, 10);
    if (isNaN(page) || page < 1 || page > authorTotalPages) return;
    handleAuthorPageChange(page);
    setAuthorJumpPage('');
  }, [authorJumpPage, authorTotalPages, handleAuthorPageChange]);

  // 生成作者页码数组
  const getAuthorPageNumbers = useCallback(() => {
    const pages: (number | string)[] = [];
    const maxVisible = 5;

    if (authorTotalPages <= maxVisible + 2) {
      for (let i = 1; i <= authorTotalPages; i++) {
        pages.push(i);
      }
    } else {
      pages.push(1);

      if (authorPage > 3) {
        pages.push('...');
      }

      const start = Math.max(2, authorPage - 1);
      const end = Math.min(authorTotalPages - 1, authorPage + 1);

      for (let i = start; i <= end; i++) {
        pages.push(i);
      }

      if (authorPage < authorTotalPages - 2) {
        pages.push('...');
      }

      if (authorTotalPages > 1) {
        pages.push(authorTotalPages);
      }
    }

    return pages;
  }, [authorTotalPages, authorPage]);

  const authorsPageElement = useMemo(() => (
    <div className="animate-in slide-in-from-bottom-4 duration-500 pb-32 md:pb-16">
      <div className="mb-6">
        <div className="relative max-w-md">
          <Search size={18} className="absolute left-3 top-1/2 -translate-y-1/2 text-gray-400 dark:text-slate-500" />
          <input
            type="text"
            placeholder="搜索作者..."
            value={authorSearchQuery}
            onChange={(e) => setAuthorSearchQuery(e.target.value)}
            className="w-full bg-gray-100 dark:bg-slate-800 border-transparent focus:bg-white dark:focus:bg-slate-900 focus:ring-2 focus:ring-wechat/50 focus:border-wechat rounded-xl py-2 pl-10 pr-10 text-sm transition-all outline-none dark:text-slate-200 dark:placeholder-slate-500"
          />
          {authorSearchQuery && (
            <button
              type="button"
              onClick={() => setAuthorSearchQuery('')}
              className="absolute right-3 top-1/2 -translate-y-1/2 text-gray-400 hover:text-gray-600 dark:hover:text-slate-300 transition-colors"
            >
              <X size={16} />
            </button>
          )}
        </div>
      </div>

      {paginatedAuthors.length > 0 ? (
        <div className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 xl:grid-cols-6 gap-3">
          {paginatedAuthors.map(author => (
            <div
              key={author.id}
              onClick={() => handleAuthorClick(author.name)}
              onContextMenu={(e) => handleAuthorContextMenu(e, author)}
              className="bg-white dark:bg-slate-800 rounded-xl p-3 border border-gray-100 dark:border-slate-700 hover:border-wechat dark:hover:border-wechat hover:shadow-md hover:shadow-wechat/10 transition-all cursor-pointer group"
            >
              <div className="flex items-center gap-3">
                {/* 头像 */}
                <div className="w-10 h-10 rounded-full bg-gradient-to-br from-wechat to-wechat-dark flex items-center justify-center text-white font-bold text-base shrink-0">
                  {author.avatar ? (
                    <img src={author.avatar} alt={author.name} className="w-full h-full rounded-full object-cover" />
                  ) : (
                    author.name.charAt(0).toUpperCase()
                  )}
                </div>

                {/* 信息 */}
                <div className="flex-1 min-w-0">
                  <h3 className="font-medium text-sm text-slate-800 dark:text-slate-100 truncate group-hover:text-wechat transition-colors">
                    {author.name}
                  </h3>
                  <span className="text-xs text-gray-400 dark:text-slate-500">
                    {author.articleCount} 篇
                  </span>
                </div>
              </div>
            </div>
          ))}
        </div>
      ) : (
        <div className="flex flex-col items-center justify-center py-20 text-gray-400 dark:text-slate-500 bg-gray-50 dark:bg-slate-800 rounded-2xl border-2 border-dashed border-gray-200 dark:border-slate-700">
          <Users size={48} className="mb-4 opacity-20" />
          <p>{authorSearchQuery ? '未找到匹配的作者' : '暂无作者数据'}</p>
          {!authorSearchQuery && <p className="text-sm mt-2">添加文章后，作者信息将自动显示在这里</p>}
        </div>
      )}

      {/* 右键菜单 */}
      {contextMenu.visible && contextMenu.author && (
        <div
          className="fixed z-[100] bg-white dark:bg-slate-800 rounded-xl shadow-xl border border-gray-200 dark:border-slate-700 py-1 min-w-[160px] animate-in fade-in zoom-in-95 duration-150"
          style={{ left: contextMenu.x, top: contextMenu.y }}
          onClick={(e) => e.stopPropagation()}
        >
          <button
            onClick={handleDeleteAuthor}
            disabled={contextMenu.author.articleCount > 0}
            className={`w-full px-4 py-2.5 text-left text-sm flex items-center gap-2 transition-colors ${
              contextMenu.author.articleCount > 0
                ? 'text-gray-300 dark:text-slate-600 cursor-not-allowed'
                : 'text-red-600 dark:text-red-400 hover:bg-red-50 dark:hover:bg-red-900/20'
            }`}
            title={contextMenu.author.articleCount > 0 ? '该作者还有文章，无法删除' : ''}
          >
            <Trash2 size={14} />
            删除作者
          </button>
        </div>
      )}

      {/* 底部固定分页栏 */}
      {filteredAuthors.length > 0 && (
        <div className={`fixed bottom-24 md:bottom-0 left-0 right-0 bg-white/95 dark:bg-slate-900/95 backdrop-blur-sm border-t border-gray-200 dark:border-slate-800 py-2 px-3 md:px-4 z-10 transition-all duration-300 ${
          isSidebarCollapsed ? 'md:left-20' : 'md:left-64'
        }`}>
          {/* 移动端简化布局 */}
          <div className="flex md:hidden items-center justify-between">
            <button
              onClick={() => handleAuthorPageChange(authorPage - 1)}
              disabled={authorPage <= 1}
              className={`p-2 rounded-lg transition-colors ${
                authorPage <= 1 ? 'text-gray-300 dark:text-slate-600' : 'text-gray-500 dark:text-slate-400 active:bg-gray-100 dark:active:bg-slate-800'
              }`}
            >
              <ChevronLeft size={20} />
            </button>
            <div className="flex items-center gap-2">
              <span className="text-sm font-medium text-slate-700 dark:text-slate-300">
                {authorPage} / {authorTotalPages || 1}
              </span>
              <span className="text-xs text-gray-400 dark:text-slate-500">
                共{filteredAuthors.length}位
              </span>
            </div>
            <button
              onClick={() => handleAuthorPageChange(authorPage + 1)}
              disabled={authorPage >= authorTotalPages}
              className={`p-2 rounded-lg transition-colors ${
                authorPage >= authorTotalPages ? 'text-gray-300 dark:text-slate-600' : 'text-gray-500 dark:text-slate-400 active:bg-gray-100 dark:active:bg-slate-800'
              }`}
            >
              <ChevronRight size={20} />
            </button>
          </div>

          {/* 桌面端完整布局 */}
          <div className="hidden md:flex items-center justify-center gap-2">
            <button
              onClick={() => handleAuthorPageChange(authorPage - 1)}
              disabled={authorPage <= 1}
              className={`p-1.5 rounded-lg transition-colors ${
                authorPage <= 1 ? 'text-gray-300 dark:text-slate-600 cursor-not-allowed' : 'text-gray-500 dark:text-slate-400 hover:bg-gray-100 dark:hover:bg-slate-800'
              }`}
            >
              <ChevronLeft size={18} />
            </button>

            <div className="flex items-center gap-1">
              {getAuthorPageNumbers().map((page, index) => (
                typeof page === 'number' ? (
                  <button
                    key={index}
                    onClick={() => handleAuthorPageChange(page)}
                    className={`min-w-[28px] h-7 px-1.5 rounded-lg text-sm font-medium transition-colors ${
                      authorPage === page
                        ? 'bg-wechat text-white'
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

            <button
              onClick={() => handleAuthorPageChange(authorPage + 1)}
              disabled={authorPage >= authorTotalPages}
              className={`p-1.5 rounded-lg transition-colors ${
                authorPage >= authorTotalPages ? 'text-gray-300 dark:text-slate-600 cursor-not-allowed' : 'text-gray-500 dark:text-slate-400 hover:bg-gray-100 dark:hover:bg-slate-800'
              }`}
            >
              <ChevronRight size={18} />
            </button>

            <div className="w-px h-5 bg-gray-200 dark:bg-slate-700 mx-1" />

            {/* 跳转输入框 */}
            <div className="flex items-center gap-1">
              <span className="text-xs text-gray-400 dark:text-slate-500">跳至</span>
              <input
                type="text"
                inputMode="numeric"
                pattern="[0-9]*"
                value={authorJumpPage}
                onChange={(e) => setAuthorJumpPage(e.target.value.replace(/\D/g, ''))}
                onKeyDown={(e) => e.key === 'Enter' && handleAuthorJumpPage()}
                className="w-10 h-7 text-center text-sm bg-gray-50 dark:bg-slate-800 border border-gray-200 dark:border-slate-700 rounded-lg outline-none focus:ring-1 focus:ring-wechat/50 focus:border-wechat dark:text-slate-200"
              />
              <span className="text-xs text-gray-400 dark:text-slate-500">页</span>
              <button
                type="button"
                onClick={handleAuthorJumpPage}
                disabled={!authorJumpPage}
                className="px-2 h-7 text-xs bg-wechat/10 text-wechat rounded-lg hover:bg-wechat/20 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
              >
                Go
              </button>
            </div>

            <span className="text-xs text-gray-400 dark:text-slate-500 whitespace-nowrap ml-2">
              共{filteredAuthors.length}位
            </span>
          </div>
        </div>
      )}
    </div>
  ), [authorSearchQuery, paginatedAuthors, filteredAuthors, authorPage, authorTotalPages, authorJumpPage, handleAuthorClick, handleAuthorContextMenu, contextMenu, handleDeleteAuthor, isSidebarCollapsed, handleAuthorPageChange, handleAuthorJumpPage, getAuthorPageNumbers, totalCount]);

  // List 页面使用独立组件，避免输入框失焦

  // Settings 页面组件
  const SettingsPage = () => (
    <div className="animate-in fade-in duration-500 h-full">
       <div className="flex flex-col md:flex-row gap-8 h-full">
          {/* Settings Sidebar */}
          <div className="w-full md:w-64 flex-shrink-0">
             <h2 className="text-2xl font-bold text-slate-800 dark:text-slate-100 mb-6 px-2">设置</h2>
             <nav className="flex flex-col space-y-1">
               <button
                 onClick={() => { setSettingsTab('general'); navigate('/settings'); }}
                 className={`flex items-center gap-3 px-4 py-3 rounded-xl transition-all font-medium text-left ${
                   settingsTab === 'general'
                     ? 'bg-white dark:bg-slate-800 shadow-sm text-wechat ring-1 ring-gray-100 dark:ring-slate-700'
                     : 'text-slate-600 dark:text-slate-400 hover:bg-white/50 dark:hover:bg-slate-800/50'
                 }`}
               >
                 <Monitor size={18} />
                 通用
               </button>
               <button
                 onClick={() => { setSettingsTab('data'); navigate('/settings#data'); }}
                 className={`flex items-center gap-3 px-4 py-3 rounded-xl transition-all font-medium text-left ${
                   settingsTab === 'data'
                     ? 'bg-white dark:bg-slate-800 shadow-sm text-wechat ring-1 ring-gray-100 dark:ring-slate-700'
                     : 'text-slate-600 dark:text-slate-400 hover:bg-white/50 dark:hover:bg-slate-800/50'
                 }`}
               >
                 <Database size={18} />
                 数据管理
               </button>
               {/* 搜索引擎设置 - 只有启用时才显示 */}
               {searchEnabled === true && (
                 <button
                   onClick={() => { setSettingsTab('search'); navigate('/settings#search'); }}
                   className={`flex items-center gap-3 px-4 py-3 rounded-xl transition-all font-medium text-left ${
                     settingsTab === 'search'
                       ? 'bg-white dark:bg-slate-800 shadow-sm text-wechat ring-1 ring-gray-100 dark:ring-slate-700'
                       : 'text-slate-600 dark:text-slate-400 hover:bg-white/50 dark:hover:bg-slate-800/50'
                   }`}
                 >
                   <Search size={18} />
                   搜索引擎
                 </button>
               )}
               {storageMigrationEnabled && (
              <button
                 onClick={() => { setSettingsTab('storage'); navigate('/settings#storage'); }}
                 className={`flex items-center gap-3 px-4 py-3 rounded-xl transition-all font-medium text-left ${
                   settingsTab === 'storage'
                     ? 'bg-white dark:bg-slate-800 shadow-sm text-wechat ring-1 ring-gray-100 dark:ring-slate-700'
                     : 'text-slate-600 dark:text-slate-400 hover:bg-white/50 dark:hover:bg-slate-800/50'
                 }`}
               >
                 <Database size={18} />
                 存储管理
               </button>
              )}
               <button
                 onClick={() => { setSettingsTab('about'); navigate('/settings#about'); }}
                 className={`flex items-center gap-3 px-4 py-3 rounded-xl transition-all font-medium text-left ${
                   settingsTab === 'about'
                     ? 'bg-white dark:bg-slate-800 shadow-sm text-wechat ring-1 ring-gray-100 dark:ring-slate-700'
                     : 'text-slate-600 dark:text-slate-400 hover:bg-white/50 dark:hover:bg-slate-800/50'
                 }`}
               >
                 <Info size={18} />
                 关于
               </button>
             </nav>
          </div>

          {/* Settings Content Area */}
          <div className="flex-1 bg-white dark:bg-slate-900 rounded-2xl shadow-sm border border-gray-100 dark:border-slate-800 p-8 min-h-[500px]">
             {settingsTab === 'general' && (
               <div className="space-y-8 animate-in slide-in-from-right-2 duration-300">
                  <div>
                     <h3 className="text-xl font-bold text-slate-800 dark:text-slate-100 border-b border-gray-100 dark:border-slate-800 pb-4 mb-6">通用设置</h3>

                     <div className="space-y-6">
                        <div className="flex items-center justify-between">
                           <div>
                             <div className="font-semibold text-slate-800 dark:text-slate-200 mb-1">默认启动视图</div>
                             <div className="text-sm text-gray-500 dark:text-slate-400">选择打开 WeKeep 时显示的页面</div>
                           </div>
                           <div className="flex bg-gray-100 dark:bg-slate-800 p-1 rounded-lg">
                              <button
                                onClick={() => handleSaveDefaultView('dashboard')}
                                className={`px-3 py-1.5 text-sm font-medium rounded-md transition-all ${
                                  defaultView === 'dashboard'
                                    ? 'bg-white dark:bg-slate-700 shadow text-slate-800 dark:text-white'
                                    : 'text-gray-500 dark:text-slate-400 hover:text-slate-700 dark:hover:text-slate-300'
                                }`}
                              >
                                仪表盘
                              </button>
                              <button
                                onClick={() => handleSaveDefaultView('list')}
                                className={`px-3 py-1.5 text-sm font-medium rounded-md transition-all ${
                                  defaultView === 'list'
                                    ? 'bg-white dark:bg-slate-700 shadow text-slate-800 dark:text-white'
                                    : 'text-gray-500 dark:text-slate-400 hover:text-slate-700 dark:hover:text-slate-300'
                                }`}
                              >
                                文章
                              </button>
                           </div>
                        </div>

                        <div className="flex items-center justify-between">
                           <div>
                             <div className="font-semibold text-slate-800 dark:text-slate-200 mb-1">外观</div>
                             <div className="text-sm text-gray-500 dark:text-slate-400">切换浅色/深色主题</div>
                           </div>
                           <button
                             onClick={() => setIsDarkMode(!isDarkMode)}
                             className={`w-14 h-7 rounded-full relative transition-colors duration-300 focus:outline-none ${isDarkMode ? 'bg-slate-700' : 'bg-gray-200'}`}
                           >
                              <div className={`absolute top-1 left-1 w-5 h-5 bg-white rounded-full shadow-sm transition-transform duration-300 flex items-center justify-center ${isDarkMode ? 'translate-x-7' : 'translate-x-0'}`}>
                                {isDarkMode ? <Moon size={12} className="text-slate-700" /> : <Sun size={12} className="text-yellow-500" />}
                              </div>
                           </button>
                        </div>
                     </div>
                  </div>
               </div>
             )}

             {settingsTab === 'data' && (
               <div className="space-y-6 animate-in slide-in-from-right-2 duration-300">
                  <div>
                     <h3 className="text-xl font-bold text-slate-800 dark:text-slate-100 border-b border-gray-100 dark:border-slate-800 pb-4 mb-4">数据管理</h3>

                     <div className="space-y-3">
                        {/* 导出备份 */}
                        <div className="flex items-center justify-between p-4 bg-white dark:bg-slate-800 rounded-lg border border-gray-200 dark:border-slate-700">
                           <div className="flex items-center gap-3">
                              <Download size={18} className="text-blue-500" />
                              <div>
                                 <div className="font-medium text-slate-800 dark:text-slate-100">导出备份</div>
                                 <div className="text-xs text-gray-500">{totalCount} 篇文章</div>
                              </div>
                           </div>
                           <button
                             onClick={handleExportData}
                             className="px-3 py-1.5 text-sm bg-blue-50 dark:bg-blue-900/20 text-blue-600 dark:text-blue-400 font-medium rounded-lg hover:bg-blue-100 dark:hover:bg-blue-900/30 transition-colors"
                           >
                             下载
                           </button>
                        </div>

                        {/* 导入恢复 */}
                        <div className="flex items-center justify-between p-4 bg-white dark:bg-slate-800 rounded-lg border border-gray-200 dark:border-slate-700">
                           <div className="flex items-center gap-3">
                              <Upload size={18} className="text-green-500" />
                              <div>
                                 <div className="font-medium text-slate-800 dark:text-slate-100">导入恢复</div>
                                 <div className="text-xs text-gray-500">从 JSON 备份文件恢复</div>
                              </div>
                           </div>
                           <label className="px-3 py-1.5 text-sm bg-green-50 dark:bg-green-900/20 text-green-600 dark:text-green-400 font-medium rounded-lg hover:bg-green-100 dark:hover:bg-green-900/30 transition-colors cursor-pointer">
                             选择文件
                             <input
                               type="file"
                               accept=".json"
                               onChange={handleImportData}
                               className="hidden"
                             />
                           </label>
                        </div>

                        {/* 数据迁移 */}
                        {needsMigration() && (
                          <div className="flex items-center justify-between p-4 bg-white dark:bg-slate-800 rounded-lg border border-green-200 dark:border-green-900/50">
                             <div className="flex items-center gap-3">
                                <Database size={18} className="text-green-500" />
                                <div>
                                   <div className="font-medium text-slate-800 dark:text-slate-100">数据迁移</div>
                                   <div className="text-xs text-gray-500">
                                     {migrationStatus || '检测到本地数据待迁移'}
                                   </div>
                                </div>
                             </div>
                             <button
                               onClick={async () => {
                                 setIsMigrating(true);
                                 setMigrationStatus('正在迁移...');
                                 try {
                                   const result = await migrateToBackend();
                                   setMigrationStatus(`完成 ${result.success} 篇`);
                                   setTimeout(() => {
                                     loadArticles();
                                     setIsMigrating(false);
                                     setMigrationStatus('');
                                   }, 2000);
                                 } catch (error) {
                                   setMigrationStatus('迁移失败');
                                   setTimeout(() => setIsMigrating(false), 2000);
                                 }
                               }}
                               disabled={isMigrating}
                               className="px-3 py-1.5 text-sm bg-green-500 text-white font-medium rounded-lg hover:bg-green-600 transition-colors disabled:opacity-50"
                             >
                               {isMigrating ? '迁移中' : '迁移'}
                             </button>
                          </div>
                        )}
                     </div>
                  </div>
               </div>
             )}

             {settingsTab === 'search' && searchEnabled && (
               <div className="space-y-6 animate-in slide-in-from-right-2 duration-300">
                  <div>
                     <h3 className="text-xl font-bold text-slate-800 dark:text-slate-100 border-b border-gray-100 dark:border-slate-800 pb-4 mb-4">搜索引擎设置</h3>

                     <div className="space-y-3">
                        {/* 搜索引擎状态 */}
                        <div className="flex items-center justify-between p-4 bg-white dark:bg-slate-800 rounded-lg border border-gray-200 dark:border-slate-700">
                           <div className="flex items-center gap-3">
                              <Search size={18} className="text-wechat" />
                              <div>
                                 <div className="font-medium text-slate-800 dark:text-slate-100">搜索引擎状态</div>
                                 <div className="text-xs text-gray-500">Meilisearch 全文搜索已启用</div>
                              </div>
                           </div>
                           <div className="flex items-center gap-2">
                              <span className="text-sm text-green-600 dark:text-green-400 font-medium">已启用</span>
                           </div>
                        </div>

                        {/* 索引状态 */}
                        <div className="flex items-center justify-between p-4 bg-white dark:bg-slate-800 rounded-lg border border-gray-200 dark:border-slate-700">
                           <div className="flex items-center gap-3">
                              <Database size={18} className="text-blue-500" />
                              <div>
                                 <div className="font-medium text-slate-800 dark:text-slate-100">索引状态</div>
                                 <div className="text-xs text-gray-500">
                                   {searchStatusLoading ? '加载中...' : `已索引 ${indexedCount} / ${totalCount} 篇文章`}
                                 </div>
                              </div>
                           </div>
                           <div className="flex items-center gap-2">
                              {indexedCount === totalCount && totalCount > 0 ? (
                                <span className="text-sm text-green-600 dark:text-green-400 font-medium">已同步</span>
                              ) : indexedCount < totalCount ? (
                                <span className="text-sm text-orange-500 font-medium">待同步</span>
                              ) : null}
                           </div>
                        </div>

                        {/* 同步按钮 */}
                        <div className="flex items-center justify-between p-4 bg-white dark:bg-slate-800 rounded-lg border border-gray-200 dark:border-slate-700">
                           <div className="flex items-center gap-3">
                              <RefreshCw size={18} className={`text-purple-500 ${isIndexing ? 'animate-spin' : ''}`} />
                              <div>
                                 <div className="font-medium text-slate-800 dark:text-slate-100">同步索引</div>
                                 <div className="text-xs text-gray-500">将所有文章同步到搜索引擎</div>
                              </div>
                           </div>
                           <button
                             onClick={handleIndexAllArticles}
                             disabled={isIndexing}
                             className="px-4 py-2 text-sm bg-wechat hover:bg-wechat-dark text-white font-medium rounded-lg transition-colors disabled:opacity-50 flex items-center gap-2"
                           >
                             {isIndexing ? (
                               <>
                                 <Loader2 size={14} className="animate-spin" />
                                 同步中...
                               </>
                             ) : (
                               <>
                                 <RefreshCw size={14} />
                                 同步
                               </>
                             )}
                           </button>
                        </div>

                        {/* 说明 */}
                        <div className="p-4 bg-blue-50 dark:bg-blue-900/20 rounded-lg border border-blue-100 dark:border-blue-900/30">
                           <div className="text-sm text-blue-700 dark:text-blue-300 space-y-2">
                             <p className="font-medium">💡 全文搜索说明</p>
                             <ul className="text-xs space-y-1 text-blue-600 dark:text-blue-400">
                               <li>• 启用后可搜索标题、作者、摘要和内容</li>
                               <li>• 未启用时仅支持标题搜索</li>
                               <li>• 新增/编辑文章会自动更新索引</li>
                               <li>• 如搜索异常请点击同步按钮重建索引</li>
                             </ul>
                           </div>
                        </div>
                     </div>
                  </div>
               </div>
             )}

             {settingsTab === 'storage' && (
               <StorageSettings />
             )}

             {settingsTab === 'about' && (
               <div className="animate-in slide-in-from-right-2 duration-300">
                  <div className="text-center py-10">
                     <div className="w-20 h-20 mx-auto mb-5 bg-gradient-to-br from-wechat to-wechat-dark rounded-2xl flex items-center justify-center text-white shadow-xl">
                        <BookOpen size={36} />
                     </div>
                     <h3 className="text-2xl font-bold text-slate-900 dark:text-white mb-1">WeKeep</h3>
                     <p className="text-sm text-gray-500 dark:text-slate-400 mb-8">微信公众号文章收藏管理工具</p>

                     <div className="flex justify-center gap-8 mb-8">
                        <div className="text-center">
                           <div className="text-2xl font-bold text-wechat">{totalCount}</div>
                           <div className="text-xs text-gray-500 mt-1">文章数</div>
                        </div>
                        <div className="w-px bg-gray-200 dark:bg-slate-700"></div>
                        <div className="text-center">
                           <div className="text-2xl font-bold text-slate-700 dark:text-slate-300">{appVersion || '...'}</div>
                           <div className="text-xs text-gray-500 mt-1">版本</div>
                        </div>
                     </div>

                     <a
                        href="https://github.com/cicbyte/wekeep"
                        target="_blank"
                        rel="noopener noreferrer"
                        className="inline-flex items-center gap-2 px-5 py-2.5 rounded-full bg-gray-100 dark:bg-slate-800 text-slate-600 dark:text-slate-300 hover:bg-gray-200 dark:hover:bg-slate-700 transition-colors text-sm font-medium"
                     >
                        <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="currentColor"><path d="M12 0C5.37 0 0 5.37 0 12c0 5.31 3.435 9.795 8.205 11.385.6.105.825-.255.825-.57 0-.285-.015-1.23-.015-2.235-3.015.555-3.795-.735-4.035-1.41-.135-.345-.72-1.41-1.23-1.695-.42-.225-1.02-.78-.015-.795.945-.015 1.62.87 1.845 1.23 1.08 1.815 2.805 1.305 3.495.99.105-.78.42-1.305.765-1.605-2.67-.3-5.46-1.335-5.46-5.925 0-1.305.465-2.385 1.23-3.225-.12-.3-.54-1.53.12-3.18 0 0 1.005-.315 3.3 1.23.96-.27 1.98-.405 3-.405s2.04.135 3 .405c2.295-1.56 3.3-1.23 3.3-1.23.66 1.65.24 2.88.12 3.18.765.84 1.23 1.905 1.23 3.225 0 4.605-2.805 5.625-5.475 5.925.435.375.81 1.095.81 2.22 0 1.605-.015 2.895-.015 3.3 0 .315.225.69.825.57A12.02 12.02 0 0024 12c0-6.63-5.37-12-12-12z"/></svg>
                        GitHub
                     </a>

                     <p className="text-xs text-gray-400 dark:text-slate-600 mt-8">
                        React + GoFrame + Meilisearch
                     </p>
                  </div>
               </div>
             )}
          </div>
       </div>
    </div>
  );

  return (
    <div className="h-screen flex flex-col md:flex-row bg-[#F7F7F7] dark:bg-slate-950 transition-colors duration-300 overflow-hidden">
      {/* Mobile Header */}
      <div className="md:hidden bg-white dark:bg-slate-900 p-4 flex justify-center items-center border-b border-gray-200 dark:border-slate-800 sticky top-0 z-20">
        <div className="flex items-center gap-2 font-bold text-xl text-slate-800 dark:text-white">
          <div className="w-8 h-8 bg-wechat rounded-lg flex items-center justify-center text-white">
            <BookOpen size={18} />
          </div>
          WeKeep
        </div>
      </div>

      {/* Sidebar Navigation (Desktop) */}
      <aside 
        className={`hidden md:flex flex-col bg-white dark:bg-slate-900 border-r border-gray-200 dark:border-slate-800 h-screen sticky top-0 z-30 transition-all duration-300 ${
          isSidebarCollapsed ? 'w-20' : 'w-64'
        }`}
      >
        <div className={`flex items-center transition-all ${isSidebarCollapsed ? 'justify-center p-4' : 'p-6'}`}>
          <div className="flex items-center gap-3 font-bold text-2xl text-slate-800 dark:text-white overflow-hidden whitespace-nowrap">
            <div className="w-10 h-10 bg-wechat rounded-xl flex items-center justify-center text-white shadow-lg shadow-wechat/30 flex-shrink-0">
              <BookOpen size={24} />
            </div>
            {!isSidebarCollapsed && <span className="transition-opacity duration-300">WeKeep</span>}
          </div>
        </div>

        <nav className="flex-1 px-4 space-y-2 mt-4 overflow-x-hidden">
          <Link
            to="/"
            title={isSidebarCollapsed ? "仪表盘" : ""}
            className={`flex items-center gap-3 w-full px-4 py-3 rounded-xl transition-all font-medium whitespace-nowrap ${
              viewMode === 'dashboard'
                ? 'bg-wechat-light dark:bg-wechat-dark/20 text-wechat-dark dark:text-wechat-light'
                : 'text-gray-500 dark:text-slate-400 hover:bg-gray-50 dark:hover:bg-slate-800'
            } ${isSidebarCollapsed ? 'justify-center px-2' : ''}`}
          >
            <LayoutGrid size={20} className="flex-shrink-0" />
            {!isSidebarCollapsed && <span>仪表盘</span>}
          </Link>
          <Link
            to="/authors"
            onClick={() => { if(authorsList.length > 0 && !selectedAuthor) handleAuthorFilter(authorsList[0].name); }}
            title={isSidebarCollapsed ? "作者" : ""}
            className={`flex items-center gap-3 w-full px-4 py-3 rounded-xl transition-all font-medium whitespace-nowrap ${
              viewMode === 'authors'
                ? 'bg-wechat-light dark:bg-wechat-dark/20 text-wechat-dark dark:text-wechat-light'
                : 'text-gray-500 dark:text-slate-400 hover:bg-gray-50 dark:hover:bg-slate-800'
            } ${isSidebarCollapsed ? 'justify-center px-2' : ''}`}
          >
            <Users size={20} className="flex-shrink-0" />
            {!isSidebarCollapsed && <span>作者</span>}
          </Link>
          <Link
            to="/list"
            onClick={() => { setSelectedAuthor(null); selectedAuthorRef.current = null; }}
            title={isSidebarCollapsed ? "文章" : ""}
            className={`flex items-center gap-3 w-full px-4 py-3 rounded-xl transition-all font-medium whitespace-nowrap ${
              viewMode === 'list'
                ? 'bg-wechat-light dark:bg-wechat-dark/20 text-wechat-dark dark:text-wechat-light'
                : 'text-gray-500 dark:text-slate-400 hover:bg-gray-50 dark:hover:bg-slate-800'
            } ${isSidebarCollapsed ? 'justify-center px-2' : ''}`}
          >
            <Menu size={20} className="flex-shrink-0" />
            {!isSidebarCollapsed && <span>文章</span>}
          </Link>
          <Link
            to="/settings"
            title={isSidebarCollapsed ? "设置" : ""}
            className={`flex items-center gap-3 w-full px-4 py-3 rounded-xl transition-all font-medium whitespace-nowrap ${
              viewMode === 'settings'
                ? 'bg-wechat-light dark:bg-wechat-dark/20 text-wechat-dark dark:text-wechat-light'
                : 'text-gray-500 dark:text-slate-400 hover:bg-gray-50 dark:hover:bg-slate-800'
            } ${isSidebarCollapsed ? 'justify-center px-2' : ''}`}
          >
            <Settings size={20} className="flex-shrink-0" />
            {!isSidebarCollapsed && <span>设置</span>}
          </Link>
        </nav>

        <div className="px-4 py-2 border-t border-gray-100 dark:border-slate-800">
           <button
             onClick={() => setIsAddModalOpen(true)}
             title={isSidebarCollapsed ? "添加文章" : ""}
             className={`w-full bg-slate-900 dark:bg-slate-700 text-white py-3 rounded-xl font-bold shadow-lg hover:shadow-xl transition-all flex items-center justify-center gap-2 hover:bg-black dark:hover:bg-slate-600 whitespace-nowrap overflow-hidden ${isSidebarCollapsed ? 'px-0' : ''}`}
           >
             <Plus size={20} className="flex-shrink-0" />
             {!isSidebarCollapsed && <span>添加文章</span>}
           </button>
        </div>
      </aside>

      {/* Main Content Area */}
      <main className="flex-1 min-w-0 overflow-y-auto">
        {/* Header */}
        <header className="sticky top-0 z-10 bg-white/80 dark:bg-slate-900/80 backdrop-blur-md border-b border-gray-200 dark:border-slate-800 px-6 py-4 flex flex-col md:flex-row md:items-center justify-between gap-4 transition-colors duration-300">

          {/* Breadcrumbs & Toggle */}
          <div className="flex-1 min-w-0 flex items-center gap-3">
             <button
              onClick={() => setIsSidebarCollapsed(!isSidebarCollapsed)}
              className="hidden md:flex p-2 -ml-2 text-gray-500 dark:text-slate-400 hover:bg-gray-100 dark:hover:bg-slate-800 hover:text-gray-700 dark:hover:text-slate-200 rounded-lg transition-colors"
              title={isSidebarCollapsed ? "展开侧边栏" : "收起侧边栏"}
            >
              <PanelLeft size={20} />
            </button>
            {renderBreadcrumbs()}
          </div>
        </header>

        <div className="p-4 md:p-6 pb-24 md:pb-8">
          {/* 全局加载提示 */}
          {loading && articles.length === 0 && (
            <div className="flex items-center justify-center h-64">
              <div className="flex flex-col items-center gap-3">
                <Loader2 size={48} className="text-wechat animate-spin" />
                <p className="text-gray-500 dark:text-slate-400">加载中...</p>
              </div>
            </div>
          )}

          {/* 迁移提示 */}
          {isMigrating && (
            <div className="mb-4 p-4 bg-green-50 dark:bg-green-900/20 border border-green-200 dark:border-green-800 rounded-xl">
              <div className="flex items-center gap-2 text-green-700 dark:text-green-400">
                <Loader2 size={20} className="animate-spin" />
                <p className="font-medium">{migrationStatus || '正在迁移数据...'}</p>
              </div>
            </div>
          )}

          {/* 主内容 */}
          <Routes>
            <Route path="/" element={<DashboardPage />} />
            <Route path="/authors" element={authorsPageElement} />
            <Route path="/list" element={
              <ListPage
                articles={displayedArticles}
                authorsList={authorsList}
                loading={loading}
                searchArticles={searchArticles}
                loadArticles={loadArticles}
                filterByAuthor={filterByAuthor}
                selectedAuthor={selectedAuthor}
                onAuthorFilter={handleAuthorFilter}
                onDeleteArticle={handleDeleteArticle}
                onReparseArticle={handleReparseArticle}
                isSidebarCollapsed={isSidebarCollapsed}
                totalCount={totalCount}
                currentPage={currentPage}
              />
            } />
            <Route path="/settings" element={<SettingsPage />} />
            <Route path="/read/:id" element={<ReaderPage />} />
            <Route path="/edit/:id" element={<EditPage />} />
            <Route path="*" element={<Navigate to="/" replace />} />
          </Routes>
        </div>
      </main>

      {/* Mobile Bottom Nav */}
      <div className="md:hidden fixed bottom-0 left-0 right-0 bg-white dark:bg-slate-900 border-t border-gray-200 dark:border-slate-800 flex justify-around p-3 z-20 pb-safe">
         <Link
            to="/"
            className={`flex flex-col items-center gap-1 p-2 rounded-lg ${viewMode === 'dashboard' ? 'text-wechat' : 'text-gray-400 dark:text-slate-500'}`}
          >
            <LayoutGrid size={20} />
            <span className="text-[10px] font-medium">首页</span>
          </Link>
          <Link
            to="/authors"
            onClick={() => { if(authorsList.length > 0 && !selectedAuthor) handleAuthorFilter(authorsList[0].name); }}
            className={`flex flex-col items-center gap-1 p-2 rounded-lg ${viewMode === 'authors' ? 'text-wechat' : 'text-gray-400 dark:text-slate-500'}`}
          >
            <Users size={20} />
            <span className="text-[10px] font-medium">作者</span>
          </Link>
          <button
             onClick={() => setIsAddModalOpen(true)}
             className="flex flex-col items-center justify-center -mt-8 bg-wechat text-white w-14 h-14 rounded-full shadow-lg shadow-wechat/40 border-4 border-[#F7F7F7] dark:border-slate-950"
          >
             <Plus size={28} />
          </button>
           <Link
            to="/list"
            onClick={() => { setSelectedAuthor(null); selectedAuthorRef.current = null; }}
            className={`flex flex-col items-center gap-1 p-2 rounded-lg ${viewMode === 'list' ? 'text-wechat' : 'text-gray-400 dark:text-slate-500'}`}
          >
            <Menu size={20} />
            <span className="text-[10px] font-medium">列表</span>
          </Link>
          <Link
            to="/settings"
            className={`flex flex-col items-center gap-1 p-2 rounded-lg ${viewMode === 'settings' ? 'text-wechat' : 'text-gray-400 dark:text-slate-500'}`}
          >
            <Settings size={20} />
            <span className="text-[10px] font-medium">设置</span>
          </Link>
      </div>

      <AddArticleModal 
        isOpen={isAddModalOpen} 
        onClose={() => setIsAddModalOpen(false)} 
        onAdd={handleAddArticle}
      />

      {/* New Modals */}
      <ReaderModal 
        article={readingArticle}
        onClose={() => setReadingArticle(null)}
        onEdit={setEditingArticle}
      />

      <EditArticleModal 
        article={editingArticle}
        isOpen={!!editingArticle}
        onClose={() => setEditingArticle(null)}
        onSave={handleUpdateArticle}
      />
      
      {/* Delete Confirmation Dialog */}
      <ConfirmDialog
        isOpen={!!articleToDelete}
        title="删除文章"
        message="确定要删除这篇文章吗?此操作无法撤销。"
        confirmLabel="删除"
        variant="danger"
        onConfirm={confirmDeleteArticle}
        onClose={() => setArticleToDelete(null)}
      />
    </div>
  );
};

// App 组件：提供路由上下文
const App: React.FC = () => {
  return (
    <BrowserRouter>
      <AppContent />
    </BrowserRouter>
  );
};

export default App;
