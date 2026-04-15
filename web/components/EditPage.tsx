import React, { useState, useEffect } from 'react';
import { useParams, useNavigate, Link } from 'react-router-dom';
import {
  ArrowLeft,
  Save,
  Settings,
  AlignLeft,
  Moon,
  Sun,
  Home,
  Loader2,
  X,
  Tag
} from 'lucide-react';
import { Article } from '../types';
import { MdEditor } from 'md-editor-rt';
import 'md-editor-rt/lib/style.css';
import { articlesApi, authorsApi } from '../services/apiService';
import { useToast } from './Toast';

const EditPage: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const { showToast } = useToast();

  const [article, setArticle] = useState<Article | null>(null);
  const [formData, setFormData] = useState<Partial<Article>>({});
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const [showMetadata, setShowMetadata] = useState(false);
  const [isDark, setIsDark] = useState(() => {
    const saved = localStorage.getItem('wekeep_dark_mode');
    if (saved !== null) return saved === 'true';
    return window.matchMedia('(prefers-color-scheme: dark)').matches;
  });
  const [previewTheme, setPreviewTheme] = useState<'default' | 'github' | 'vuepress' | 'mk-cute' | 'smart-blue' | 'cyanosis'>('github');

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
          const articleData: Article = {
            id: String(found.id),
            title: found.title,
            author: found.authorName || `作者${found.authorId || 0}`,
            authorId: found.authorId,
            url: found.url || '',
            summary: found.summary || '',
            content: found.content || '',
            tags: found.tags || [],
            dateAdded: found.dateAdded || Date.now(),
          };
          setArticle(articleData);
          setFormData(articleData);
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

  const handleSave = async () => {
    if (!formData.title?.trim()) {
      showToast('请输入文章标题', 'error');
      return;
    }

    if (!formData.author?.trim()) {
      showToast('请输入作者名称', 'error');
      return;
    }

    setSaving(true);

    try {
      // 获取或创建作者
      let authorId = article?.authorId;
      if (formData.author !== article?.author) {
        // 作者变更，查找或创建新作者
        const authorsResponse = await authorsApi.list({ name: formData.author, pageSize: 1 });
        if (authorsResponse.authorsList && authorsResponse.authorsList.length > 0) {
          const existingAuthor = authorsResponse.authorsList.find(a => a.name === formData.author);
          if (existingAuthor) {
            authorId = existingAuthor.id;
          } else {
            const newAuthor = await authorsApi.add({ name: formData.author });
            authorId = newAuthor.id;
          }
        } else {
          const newAuthor = await authorsApi.add({ name: formData.author });
          authorId = newAuthor.id;
        }
      }

      // 更新文章
      await articlesApi.edit(parseInt(id!), {
        id: parseInt(id!),
        title: formData.title,
        authorId: authorId!,
        url: formData.url || '',
        summary: formData.summary || '',
        content: formData.content || '',
        tags: formData.tags || [],
        dateAdded: formData.dateAdded || Date.now(),
      });

      showToast('保存成功', 'success');
      // 保存后留在编辑页面，不跳转
    } catch (err: any) {
      console.error('保存文章失败:', err);
      showToast(err?.message || '保存失败', 'error');
    } finally {
      setSaving(false);
    }
  };

  const handleBack = () => {
    if (window.history.length > 1) {
      navigate(-1);
    } else {
      navigate(`/read/${id}`);
    }
  };

  const handleContentChange = (content: string) => {
    setFormData({ ...formData, content });
  };

  // 编辑器工具栏
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

  // 加载状态
  if (loading) {
    return (
      <div className={`fixed inset-0 z-50 flex items-center justify-center ${isDark ? 'bg-slate-900' : 'bg-white'}`}>
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
      <div className={`fixed inset-0 z-50 flex flex-col items-center justify-center ${isDark ? 'bg-slate-900' : 'bg-white'}`}>
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
    <div className={`fixed inset-0 z-50 flex flex-col ${isDark ? 'bg-slate-900' : 'bg-white'}`}>
      {/* Top Navigation Bar */}
      <div className={`h-14 border-b ${isDark ? 'border-slate-700' : 'border-gray-200'} flex items-center justify-between px-4 shrink-0 z-20 shadow-sm`}>
        <div className="flex items-center gap-3 flex-1 overflow-hidden">
          <button
            onClick={handleBack}
            className={`p-2 -ml-2 rounded-full ${isDark ? 'hover:bg-slate-800' : 'hover:bg-gray-100'} ${isDark ? 'text-slate-400' : 'text-gray-500'} transition-colors`}
            title="返回"
          >
            <ArrowLeft size={20} />
          </button>

          <Link
            to="/"
            className={`p-2 rounded-full ${isDark ? 'hover:bg-slate-800' : 'hover:bg-gray-100'} ${isDark ? 'text-slate-400' : 'text-gray-500'} transition-colors`}
            title="首页"
          >
            <Home size={18} />
          </Link>

          <div className="flex flex-col flex-1 min-w-0">
            <div className={`font-bold text-base truncate ${isDark ? 'text-white' : 'text-slate-800'}`}>
              {formData.title || '未命名文章'}
            </div>
            <div className={`text-xs truncate ${isDark ? 'text-slate-400' : 'text-gray-500'}`}>
              {formData.author || '未知作者'}
            </div>
          </div>
        </div>

        <div className="flex items-center gap-2">
          {/* 预览主题选择 */}
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
            className={`p-2 rounded-lg ${isDark ? 'hover:bg-slate-800' : 'hover:bg-gray-100'} ${isDark ? 'text-slate-400' : 'text-gray-500'} transition-colors`}
            title={isDark ? '切换到浅色模式' : '切换到深色模式'}
          >
            {isDark ? <Sun size={18} /> : <Moon size={18} />}
          </button>

          {/* 元数据弹框按钮 */}
          <button
            onClick={() => setShowMetadata(!showMetadata)}
            className={`p-2 rounded-lg transition-colors ${showMetadata ? 'bg-wechat text-white' : `${isDark ? 'text-slate-400 hover:bg-slate-800' : 'text-gray-500 hover:bg-gray-100'}`}`}
            title="编辑文章信息"
          >
            <Settings size={18} />
          </button>

          {/* 保存按钮 */}
          <button
            onClick={handleSave}
            disabled={saving}
            className="px-4 py-2 bg-wechat hover:bg-wechat-dark disabled:bg-wechat/50 text-white font-bold rounded-lg shadow-lg shadow-wechat/20 transition-all flex items-center gap-2"
          >
            {saving ? <Loader2 size={16} className="animate-spin" /> : <Save size={16} />}
            <span className="hidden sm:inline">{saving ? '保存中...' : '保存'}</span>
          </button>
        </div>
      </div>

      {/* Metadata Modal */}
      {showMetadata && (
        <div className="fixed inset-0 z-[60] flex items-center justify-center p-4">
          {/* Backdrop */}
          <div
            className={`absolute inset-0 ${isDark ? 'bg-black/60' : 'bg-black/30'}`}
            onClick={() => setShowMetadata(false)}
          />

          {/* Modal */}
          <div className={`relative w-full max-w-2xl max-h-[80vh] overflow-y-auto rounded-2xl shadow-2xl ${isDark ? 'bg-slate-800' : 'bg-white'} animate-in fade-in zoom-in-95 duration-200`}>
            {/* Modal Header */}
            <div className={`sticky top-0 flex items-center justify-between p-6 border-b ${isDark ? 'border-slate-700' : 'border-gray-200'} ${isDark ? 'bg-slate-800' : 'bg-white'}`}>
              <h2 className={`text-xl font-bold ${isDark ? 'text-white' : 'text-slate-800'}`}>
                文章信息
              </h2>
              <button
                onClick={() => setShowMetadata(false)}
                className={`p-2 rounded-full ${isDark ? 'hover:bg-slate-700' : 'hover:bg-gray-100'} transition-colors`}
              >
                <X size={20} className={isDark ? 'text-slate-400' : 'text-gray-500'} />
              </button>
            </div>

            {/* Modal Body */}
            <div className="p-6 space-y-6">
              {/* 摘要 */}
              <div className="space-y-2">
                <label className={`text-sm font-semibold ${isDark ? 'text-slate-300' : 'text-gray-700'} flex items-center gap-2`}>
                  <AlignLeft size={14} /> 摘要
                </label>
                <textarea
                  value={formData.summary || ''}
                  onChange={(e) => setFormData({ ...formData, summary: e.target.value })}
                  rows={6}
                  className={`w-full p-3 ${isDark ? 'bg-slate-900 border-slate-700 text-white placeholder-slate-500' : 'bg-gray-50 border-gray-200 text-slate-800 placeholder-gray-400'} border rounded-xl text-base outline-none focus:border-wechat focus:ring-2 focus:ring-wechat/20 resize-none`}
                  placeholder="输入文章摘要..."
                />
              </div>

              {/* 标签 */}
              <div className="space-y-2">
                <label className={`text-sm font-semibold ${isDark ? 'text-slate-300' : 'text-gray-700'} flex items-center gap-2`}>
                  <Tag size={14} /> 标签
                </label>
                <input
                  value={(formData.tags || []).join(', ')}
                  onChange={(e) => setFormData({ ...formData, tags: e.target.value.split(',').map(t => t.trim()).filter(Boolean) })}
                  className={`w-full p-3 ${isDark ? 'bg-slate-900 border-slate-700 text-white placeholder-slate-500' : 'bg-gray-50 border-gray-200 text-slate-800 placeholder-gray-400'} border rounded-xl text-base outline-none focus:border-wechat focus:ring-2 focus:ring-wechat/20`}
                  placeholder="标签1, 标签2, 标签3..."
                />
                <p className={`text-xs ${isDark ? 'text-slate-500' : 'text-gray-400'}`}>多个标签用逗号分隔</p>
              </div>
            </div>

            {/* Modal Footer */}
            <div className={`sticky bottom-0 flex justify-end gap-3 p-6 border-t ${isDark ? 'border-slate-700' : 'border-gray-200'} ${isDark ? 'bg-slate-800' : 'bg-white'}`}>
              <button
                onClick={() => setShowMetadata(false)}
                className={`px-5 py-2.5 rounded-xl font-medium ${isDark ? 'bg-slate-700 text-slate-300 hover:bg-slate-600' : 'bg-gray-100 text-gray-700 hover:bg-gray-200'} transition-colors`}
              >
                取消
              </button>
              <button
                onClick={() => setShowMetadata(false)}
                className="px-5 py-2.5 bg-wechat hover:bg-wechat-dark text-white font-bold rounded-xl shadow-lg shadow-wechat/20 transition-all"
              >
                确定
              </button>
            </div>
          </div>
        </div>
      )}

      {/* Markdown Editor */}
      <div className="flex-1 min-h-0">
        <div className={`h-full ${isDark ? 'dark' : ''}`}>
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
              handleSave();
            }}
          />
        </div>
      </div>
    </div>
  );
};

export default EditPage;
