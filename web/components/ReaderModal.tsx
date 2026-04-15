
import React, { useState, useEffect } from 'react';
import { X, ExternalLink, Calendar, User, Edit3, Moon, Sun } from 'lucide-react';
import { Article } from '../types';
import { MdPreview, MdCatalog } from 'md-editor-rt';
import 'md-editor-rt/lib/preview.css';

const IMAGE_BASE_URL = (import.meta.env.VITE_API_BASE_URL || '/api/v1') + '/images/file/';

function processContent(content: string): string {
  return content.replace(
    /!\[([^\]]*)\]\(([^)]+)\)/g,
    (_, alt, src) => {
      if (src.startsWith('http://') || src.startsWith('https://') || src.startsWith('/')) return `![${alt}](${src})`;
      return `![${alt}](${IMAGE_BASE_URL}${src})`;
    }
  );
}

interface ReaderModalProps {
  article: Article | null;
  onClose: () => void;
  onEdit: (article: Article) => void;
}

const ReaderModal: React.FC<ReaderModalProps> = ({ article, onClose, onEdit }) => {
  const [isDark, setIsDark] = useState(false);
  const [previewTheme, setPreviewTheme] = useState<'default' | 'github' | 'vuepress' | 'mk-cute' | 'smart-blue' | 'cyanosis'>('github');
  const [showCatalog, setShowCatalog] = useState(false);
  const editorId = 'reader-preview';

  // 检测系统主题
  useEffect(() => {
    const darkMode = window.matchMedia('(prefers-color-scheme: dark)').matches;
    setIsDark(darkMode);
  }, []);

  // 锁定 body 滚动
  useEffect(() => {
    const originalOverflow = document.body.style.overflow;
    document.body.style.overflow = 'hidden';
    return () => {
      document.body.style.overflow = originalOverflow;
    };
  }, []);

  // 监听 Escape 键关闭
  useEffect(() => {
    const handleEscape = (e: KeyboardEvent) => {
      if (e.key === 'Escape') {
        onClose();
      }
    };
    window.addEventListener('keydown', handleEscape);
    return () => window.removeEventListener('keydown', handleEscape);
  }, [onClose]);

  if (!article) return null;

  return (
    <div className={`fixed inset-0 z-[70] overflow-hidden flex ${isDark ? 'dark' : ''}`}>
      {/* 目录侧边栏 */}
      {showCatalog && article.content && (
        <div className="w-64 bg-gray-50 dark:bg-slate-800 border-r border-gray-200 dark:border-slate-700 overflow-y-auto shrink-0 hidden lg:block">
          <div className="p-4">
            <h3 className="text-sm font-bold text-gray-500 dark:text-slate-400 uppercase mb-4">目录</h3>
            <MdCatalog
              editorId={editorId}
              theme={isDark ? 'dark' : 'light'}
              className="!bg-transparent"
            />
          </div>
        </div>
      )}

      {/* 主内容区 */}
      <div className={`flex-1 overflow-y-auto ${isDark ? 'bg-slate-900' : 'bg-white'}`}>
        {/* Header */}
        <div className={`sticky top-0 ${isDark ? 'bg-slate-900/80' : 'bg-white/80'} backdrop-blur-md border-b ${isDark ? 'border-slate-800' : 'border-gray-100'} px-4 py-3 flex justify-between items-center max-w-4xl mx-auto w-full z-10`}>
          <div className="flex items-center gap-2">
            <button
              onClick={onClose}
              className={`p-2 hover:${isDark ? 'bg-slate-800' : 'bg-gray-100'} rounded-full transition-colors ${isDark ? 'text-slate-400' : 'text-gray-500'}`}
            >
              <X size={24} />
            </button>

            {/* 主题选择 */}
            <select
              value={previewTheme}
              onChange={(e) => setPreviewTheme(e.target.value as typeof previewTheme)}
              className={`text-xs ${isDark ? 'bg-slate-800 text-slate-300' : 'bg-gray-100 text-gray-600'} border-none rounded-lg px-2 py-1.5 outline-none cursor-pointer`}
            >
              <option value="github">GitHub</option>
              <option value="vuepress">VuePress</option>
              <option value="mk-cute">Cute</option>
              <option value="smart-blue">Smart Blue</option>
              <option value="cyanosis">Cyanosis</option>
            </select>

            {/* 主题切换 */}
            <button
              onClick={() => setIsDark(!isDark)}
              className={`p-2 rounded-lg hover:${isDark ? 'bg-slate-800' : 'bg-gray-100'} transition-colors ${isDark ? 'text-slate-400' : 'text-gray-500'}`}
            >
              {isDark ? <Sun size={18} /> : <Moon size={18} />}
            </button>

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
          </div>

          <div className="flex gap-2">
            <button
               onClick={() => {
                 onClose();
                 onEdit(article);
               }}
               className={`flex items-center gap-2 px-4 py-2 text-sm font-medium ${isDark ? 'text-slate-300 hover:bg-slate-800' : 'text-gray-600 hover:bg-gray-100'} rounded-full transition-colors`}
            >
              <Edit3 size={18} />
              编辑
            </button>
            {article.url && (
              <a
                href={article.url}
                target="_blank"
                rel="noopener noreferrer"
                className="flex items-center gap-2 px-4 py-2 text-sm font-medium text-wechat hover:bg-wechat-light dark:hover:bg-wechat-dark/20 rounded-full transition-colors"
              >
                <ExternalLink size={18} />
                原文
              </a>
            )}
          </div>
        </div>

        {/* Content */}
        <div className="max-w-3xl mx-auto px-6 py-12">
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
            <MdPreview
              id={editorId}
              modelValue={processContent(article.content)}
              theme={isDark ? 'dark' : 'light'}
              previewTheme={previewTheme}
              className="!bg-transparent"
            />
          ) : (
            <div className={`text-center ${isDark ? 'text-slate-500' : 'text-gray-400'} py-20 flex flex-col items-center`}>
              <p className="mb-4">本文暂无内容。</p>
              <button
                onClick={() => {
                  onClose();
                  onEdit(article);
                }}
                className="text-wechat hover:underline"
              >
                手动添加内容
              </button>
            </div>
          )}
        </div>
      </div>
    </div>
  );
};

export default ReaderModal;
