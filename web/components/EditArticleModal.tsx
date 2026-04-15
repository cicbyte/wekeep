
import React, { useState, useEffect } from 'react';
import {
  X,
  Save,
  Settings,
  Eye,
  PenLine,
  Link as LinkIcon,
  AlignLeft,
  Moon,
  Sun
} from 'lucide-react';
import { Article } from '../types';
import { MdEditor } from 'md-editor-rt';
import 'md-editor-rt/lib/style.css';

interface EditArticleModalProps {
  article: Article | null;
  isOpen: boolean;
  onClose: () => void;
  onSave: (updatedArticle: Article) => void;
}

const EditArticleModal: React.FC<EditArticleModalProps> = ({ article, isOpen, onClose, onSave }) => {
  const [formData, setFormData] = useState<Partial<Article>>({});
  const [showMetadata, setShowMetadata] = useState(true);
  const [isDark, setIsDark] = useState(false);
  const [previewTheme, setPreviewTheme] = useState<'default' | 'github' | 'vuepress' | 'mk-cute' | 'smart-blue' | 'cyanosis'>('github');

  // 检测系统主题
  useEffect(() => {
    const darkMode = window.matchMedia('(prefers-color-scheme: dark)').matches;
    setIsDark(darkMode);
  }, []);

  useEffect(() => {
    if (article) {
      setFormData({ ...article });
      setShowMetadata(false);
    }
  }, [article]);

  // 锁定 body 滚动
  useEffect(() => {
    if (isOpen) {
      const originalOverflow = document.body.style.overflow;
      document.body.style.overflow = 'hidden';
      return () => {
        document.body.style.overflow = originalOverflow;
      };
    }
  }, [isOpen]);

  // 监听 Escape 键关闭
  useEffect(() => {
    const handleEscape = (e: KeyboardEvent) => {
      if (e.key === 'Escape' && isOpen) {
        onClose();
      }
    };
    window.addEventListener('keydown', handleEscape);
    return () => window.removeEventListener('keydown', handleEscape);
  }, [isOpen, onClose]);

  if (!isOpen || !article) return null;

  const handleSubmit = (e?: React.FormEvent) => {
    if (e) e.preventDefault();
    if (formData) {
      onSave({ ...article, ...formData } as Article);
      onClose();
    }
  };

  const handleContentChange = (content: string) => {
    setFormData({ ...formData, content });
  };

  // 自定义工具栏
  const editorToolbars = [
    'bold',
    'underline',
    'italic',
    '-',
    'strikeThrough',
    'title',
    'sub',
    'sup',
    'quote',
    'unorderedList',
    'orderedList',
    'task',
    '-',
    'codeRow',
    'code',
    'link',
    'image',
    'table',
    'mermaid',
    '-',
    'revoke',
    'next',
    '=',
    'pageFullscreen',
    'fullscreen',
    'preview',
    'catalog',
  ];

  return (
    <div className="fixed inset-0 bg-white dark:bg-slate-900 z-[80] flex flex-col animate-in fade-in duration-200">
      {/* Top Navigation Bar */}
      <div className="h-14 border-b border-gray-200 dark:border-slate-700 flex items-center justify-between px-4 bg-white dark:bg-slate-900 shrink-0 z-20 shadow-sm">
        <div className="flex items-center gap-3 flex-1 overflow-hidden">
          <button
            onClick={onClose}
            className="p-2 -ml-2 rounded-full hover:bg-gray-100 dark:hover:bg-slate-800 text-gray-500 dark:text-slate-400 transition-colors"
            title="关闭编辑器 (Esc)"
          >
            <X size={20} />
          </button>

          <div className="flex flex-col flex-1 min-w-0">
             <input
               value={formData.title || ''}
               onChange={(e) => setFormData({...formData, title: e.target.value})}
               className="font-bold text-base text-slate-800 dark:text-white border-none outline-none bg-transparent placeholder-gray-400 dark:placeholder-slate-600 w-full truncate focus:ring-0 p-0"
               placeholder="文章标题"
             />
             <input
               value={formData.author || ''}
               onChange={(e) => setFormData({...formData, author: e.target.value})}
               className="text-xs text-gray-500 dark:text-slate-400 border-none outline-none bg-transparent placeholder-gray-300 dark:placeholder-slate-600 w-full truncate focus:ring-0 p-0"
               placeholder="作者名称"
             />
          </div>
        </div>

        <div className="flex items-center gap-2">
          {/* 预览主题选择 */}
          <select
            value={previewTheme}
            onChange={(e) => setPreviewTheme(e.target.value as typeof previewTheme)}
            className="text-xs bg-gray-100 dark:bg-slate-800 border-none rounded-lg px-2 py-1.5 text-gray-600 dark:text-slate-300 outline-none cursor-pointer"
          >
            <option value="github">GitHub</option>
            <option value="vuepress">VuePress</option>
            <option value="mk-cute"> Cute</option>
            <option value="smart-blue">Smart Blue</option>
            <option value="cyanosis">Cyanosis</option>
            <option value="default">Default</option>
          </select>

          {/* 主题切换 */}
          <button
            onClick={() => setIsDark(!isDark)}
            className="p-2 rounded-lg hover:bg-gray-100 dark:hover:bg-slate-800 text-gray-500 dark:text-slate-400 transition-colors"
            title={isDark ? '切换到亮色模式' : '切换到暗色模式'}
          >
            {isDark ? <Sun size={18} /> : <Moon size={18} />}
          </button>

          {/* 元数据 */}
          <button
            onClick={() => setShowMetadata(!showMetadata)}
            className={`p-2 rounded-lg transition-colors ${showMetadata ? 'bg-gray-100 dark:bg-slate-800 text-slate-800 dark:text-white' : 'text-gray-500 dark:text-slate-400 hover:bg-gray-100 dark:hover:bg-slate-800'}`}
            title="切换元数据"
          >
            <Settings size={18} />
          </button>

          {/* 保存按钮 */}
          <button
            onClick={() => handleSubmit()}
            className="px-4 py-2 bg-wechat hover:bg-wechat-dark text-white font-bold rounded-lg shadow-lg shadow-wechat/20 transition-all flex items-center gap-2"
          >
            <Save size={16} />
            <span className="hidden sm:inline">保存</span>
          </button>
        </div>
      </div>

      {/* Metadata Panel */}
      {showMetadata && (
        <div className="bg-gray-50 dark:bg-slate-800 border-b border-gray-200 dark:border-slate-700 p-4 grid grid-cols-1 md:grid-cols-2 gap-4 animate-in slide-in-from-top-2 duration-200 shrink-0">
          <div className="space-y-1">
             <label className="text-xs font-bold text-gray-400 dark:text-slate-500 uppercase flex items-center gap-1.5">
                <LinkIcon size={12} /> 来源链接
             </label>
             <input
                value={formData.url || ''}
                onChange={(e) => setFormData({...formData, url: e.target.value})}
                className="w-full p-2 bg-white dark:bg-slate-900 border border-gray-200 dark:border-slate-700 rounded-lg text-sm text-slate-600 dark:text-slate-300 outline-none focus:border-wechat focus:ring-1 focus:ring-wechat"
                placeholder="https://..."
             />
          </div>
          <div className="space-y-1">
             <label className="text-xs font-bold text-gray-400 dark:text-slate-500 uppercase flex items-center gap-1.5">
                <AlignLeft size={12} /> 摘要
             </label>
             <textarea
                value={formData.summary || ''}
                onChange={(e) => setFormData({...formData, summary: e.target.value})}
                rows={1}
                className="w-full p-2 bg-white dark:bg-slate-900 border border-gray-200 dark:border-slate-700 rounded-lg text-sm text-slate-600 dark:text-slate-300 outline-none focus:border-wechat focus:ring-1 focus:ring-wechat resize-none"
                placeholder="简短摘要..."
             />
          </div>
        </div>
      )}

      {/* Markdown Editor */}
      <div className="flex-1 overflow-hidden">
        <MdEditor
          modelValue={formData.content || ''}
          onChange={handleContentChange}
          language="zh-CN"
          theme={isDark ? 'dark' : 'light'}
          previewTheme={previewTheme}
          style={{ height: '100%' }}
          className="!border-none"
          toolbars={editorToolbars}
          placeholder="开始编写你的文章..."
          showCodeRowNumber
          previewOnly={false}
          onSave={(v) => {
            setFormData({ ...formData, content: v });
            handleSubmit();
          }}
        />
      </div>
    </div>
  );
};

export default EditArticleModal;
