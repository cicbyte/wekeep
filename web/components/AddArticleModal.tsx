
import React, { useState } from 'react';
import { X, Globe, Loader2 } from 'lucide-react';
import { articlesApi } from '../services/apiService';
import { Article } from '../types';
import { useToast } from './Toast';

interface AddArticleModalProps {
  isOpen: boolean;
  onClose: () => void;
  onAdd: (article: Article) => void;
}

const AddArticleModal: React.FC<AddArticleModalProps> = ({ isOpen, onClose, onAdd }) => {
  const [input, setInput] = useState('');
  const [loading, setLoading] = useState(false);
  const { showToast } = useToast();

  if (!isOpen) return null;

  const handleUrlParse = async () => {
    if (!input.trim()) {
      showToast('请输入至少一个 URL 链接', 'error');
      return;
    }

    setLoading(true);

    const urls = input.split('\n')
      .map(url => url.trim())
      .filter(url => url.length > 0 && url.startsWith('http'));

    if (urls.length === 0) {
      showToast('请输入有效的 URL 链接', 'error');
      setLoading(false);
      return;
    }

    try {
      let successCount = 0;
      let failCount = 0;
      let failMessages: string[] = [];

      for (const url of urls) {
        try {
          const result = await articlesApi.parseByUrl(url);
          const newArticle: Article = {
            id: `article-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`,
            title: result.title || '未知标题',
            author: result.author || '微信公众号',
            url: url,
            summary: '',
            content: result.content || '',
            tags: [],
            dateAdded: Date.now()
          };
          await onAdd(newArticle);
          successCount++;
        } catch (err: any) {
          console.error(`处理 URL ${url} 失败:`, err);
          failCount++;
          // 解析错误原因
          const errMsg = err?.message || err?.toString() || '';
          if (errMsg.includes('已经收藏') || errMsg.includes('重复')) {
            failMessages.push('部分文章已收藏');
          } else if (errMsg.includes('解析') || errMsg.includes('parse')) {
            failMessages.push('部分链接解析失败');
          } else {
            failMessages.push(errMsg.slice(0, 30) || '添加失败');
          }
        }
      }

      if (successCount > 0 && failCount === 0) {
        showToast(`成功添加 ${successCount} 篇文章`, 'success');
        onClose();
        setInput('');
      } else if (successCount > 0 && failCount > 0) {
        showToast(`成功 ${successCount} 篇，${[...new Set(failMessages)].join('、')}`, 'error');
        onClose();
        setInput('');
      } else {
        // 全部失败，显示具体原因
        const uniqueMessages = [...new Set(failMessages)];
        showToast(uniqueMessages[0] || '添加失败', 'error');
      }
    } catch (e: any) {
      console.error('解析过程中出错:', e);
      showToast('解析失败，请检查 URL 是否有效', 'error');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="fixed inset-0 bg-black/50 backdrop-blur-sm flex items-center justify-center p-4 z-50">
      <div className="bg-white dark:bg-slate-900 rounded-2xl w-full max-w-lg shadow-2xl overflow-hidden flex flex-col max-h-[90vh] transition-colors">
        <div className="p-4 border-b border-gray-100 dark:border-slate-800 flex justify-between items-center bg-gray-50/50 dark:bg-slate-800/50 flex-shrink-0">
          <h2 className="text-lg font-bold text-gray-800 dark:text-slate-100">添加文章</h2>
          <button onClick={onClose} className="p-1 rounded-full hover:bg-gray-200 dark:hover:bg-slate-700 transition-colors">
            <X size={20} className="text-gray-500 dark:text-slate-400" />
          </button>
        </div>

        <div className="p-6 overflow-y-auto">
          <div className="space-y-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-slate-300 mb-1">
                微信文章链接
              </label>
              <textarea
                value={input}
                onChange={(e) => setInput(e.target.value)}
                placeholder="请输入微信公众号文章的 URL 链接，支持多行输入，一行一个 URL。"
                className="w-full h-24 p-3 border border-gray-200 dark:border-slate-700 bg-white dark:bg-slate-800 rounded-xl focus:ring-2 focus:ring-wechat focus:border-transparent outline-none resize-none text-sm leading-relaxed text-slate-700 dark:text-slate-300 placeholder-gray-400 dark:placeholder-slate-500"
              />
            </div>

            <button
              onClick={handleUrlParse}
              disabled={loading || !input.trim()}
              className="w-full py-3 bg-gradient-to-r from-wechat to-wechat-dark text-white rounded-xl font-bold hover:shadow-lg hover:shadow-wechat/30 transition-all disabled:opacity-50 disabled:cursor-not-allowed flex justify-center items-center gap-2"
            >
              {loading ? <Loader2 className="animate-spin" /> : <Globe size={18} />}
              {loading ? '解析中...' : '解析并保存'}
            </button>

            <p className="text-xs text-gray-400 dark:text-slate-500 text-center">
              智能解析微信公众号文章内容。支持多行输入，一行一个 URL。
            </p>
          </div>
        </div>
      </div>
    </div>
  );
};

export default AddArticleModal;
