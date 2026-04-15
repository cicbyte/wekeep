import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { Article } from '../types';
import { ExternalLink, Tag, Trash2, BookOpen, Edit3, RefreshCw } from 'lucide-react';

interface ArticleCardProps {
  article: Article;
  onDelete: (id: string) => void;
  onRead?: (article: Article) => void; // 保留兼容性，但不再使用
  onEdit?: (article: Article) => void;
  onReparse?: (article: Article) => void;
}

const ArticleCard: React.FC<ArticleCardProps> = ({ article, onDelete, onReparse }) => {
  const navigate = useNavigate();
  const [contextMenu, setContextMenu] = useState<{ visible: boolean; x: number; y: number }>({ visible: false, x: 0, y: 0 });

  const handleRead = () => {
    navigate(`/read/${article.id}`);
  };

  const handleEdit = () => {
    navigate(`/edit/${article.id}`);
  };

  // 右键菜单处理
  const handleContextMenu = (e: React.MouseEvent) => {
    e.preventDefault();
    e.stopPropagation();
    setContextMenu({
      visible: true,
      x: e.clientX,
      y: e.clientY
    });
  };

  const handleCloseContextMenu = () => {
    setContextMenu({ visible: false, x: 0, y: 0 });
  };

  const handleDelete = () => {
    handleCloseContextMenu();
    onDelete(article.id);
  };

  const handleReparse = () => {
    handleCloseContextMenu();
    if (onReparse) {
      onReparse(article);
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

  // 判断是否有搜索高亮
  const hasHighlight = article.formattedTitle || article.formattedSummary || article.contextSnippet;
  // 检查是否有正文匹配：如果有 contextSnippet 且包含 <mark> 标签，说明是正文匹配
  const hasContentMatch = article.contextSnippet && article.contextSnippet.includes('<mark>');

  // 显示标题（优先使用高亮版本）
  const displayTitle = article.formattedTitle || article.title;
  // 显示摘要（优先使用高亮版本）
  const displaySummary = article.formattedSummary || article.summary;

  return (
    <>
      <div
        onClick={handleRead}
        onContextMenu={handleContextMenu}
        className="bg-white dark:bg-slate-800 rounded-lg p-3 shadow-sm border border-gray-100 dark:border-slate-700 hover:shadow-md transition-all duration-200 group relative flex flex-col h-full cursor-pointer"
      >
        <h3
          className="font-semibold text-sm text-slate-900 dark:text-slate-100 leading-tight line-clamp-2 mb-1.5"
          dangerouslySetInnerHTML={{ __html: displayTitle }}
        />

        <div className="text-xs text-gray-400 dark:text-slate-500 mb-2 flex items-center gap-2">
          <span className="text-wechat dark:text-wechat-light font-medium truncate max-w-[60%]">
            {article.author}
          </span>
          <span className="shrink-0">
            {new Date(article.dateAdded).toLocaleDateString()}
          </span>
        </div>

        {/* 摘要或正文匹配上下文 */}
        <div className="flex-grow mb-2">
          {hasContentMatch ? (
            // 正文匹配时显示上下文
            <p
              className="text-xs text-slate-600 dark:text-slate-300 line-clamp-4 bg-yellow-50 dark:bg-yellow-900/20 px-2 py-1 rounded"
              dangerouslySetInnerHTML={{ __html: article.contextSnippet }}
            />
          ) : (
            // 显示摘要（可能带高亮）
            <p
              className="text-gray-500 dark:text-slate-400 text-xs line-clamp-4"
              dangerouslySetInnerHTML={{ __html: displaySummary || "暂无摘要。" }}
            />
          )}
        </div>

        <div className="mt-auto space-y-2">
          {article.tags.length > 0 && (
            <div className="flex items-center gap-1 overflow-hidden flex-wrap">
              {article.tags.slice(0, 2).map(tag => (
                <span key={tag} className="inline-flex items-center text-[10px] text-gray-500 dark:text-slate-400 bg-gray-100 dark:bg-slate-700 px-1.5 py-0.5 rounded">
                  <Tag size={8} className="mr-0.5" />
                  {tag}
                </span>
              ))}
              {article.tags.length > 2 && (
                <span className="text-[10px] text-gray-400">+{article.tags.length - 2}</span>
              )}
            </div>
          )}

          <div className="flex items-center justify-between pt-2 border-t border-gray-50 dark:border-slate-700/50" onClick={(e) => e.stopPropagation()}>
            <div className="flex gap-1.5">
              <button
                onClick={handleEdit}
                className="flex items-center gap-1 text-xs font-medium text-gray-500 dark:text-slate-400 hover:text-slate-700 dark:hover:text-slate-200 transition-colors px-2 py-1 rounded hover:bg-gray-100 dark:hover:bg-slate-700"
              >
                <Edit3 size={12} />
                编辑
              </button>
            </div>

            {article.url && (
              <a
                href={article.url}
                target="_blank"
                rel="noopener noreferrer"
                onClick={(e) => e.stopPropagation()}
                className="text-gray-300 dark:text-slate-600 hover:text-wechat dark:hover:text-wechat transition-colors"
                title="打开原文"
              >
                <ExternalLink size={14} />
              </a>
            )}
          </div>
        </div>
      </div>

      {/* 右键菜单 */}
      {contextMenu.visible && (
        <div
          className="fixed z-[100] bg-white dark:bg-slate-800 rounded-xl shadow-xl border border-gray-200 dark:border-slate-700 py-1 min-w-[140px] animate-in fade-in zoom-in-95 duration-150"
          style={{ left: contextMenu.x, top: contextMenu.y }}
          onClick={(e) => e.stopPropagation()}
        >
          <button
            onClick={handleRead}
            className="w-full px-4 py-2.5 text-left text-sm flex items-center gap-2 text-slate-700 dark:text-slate-300 hover:bg-gray-50 dark:hover:bg-slate-700 transition-colors"
          >
            <BookOpen size={14} />
            阅读
          </button>
          <button
            onClick={handleEdit}
            className="w-full px-4 py-2.5 text-left text-sm flex items-center gap-2 text-slate-700 dark:text-slate-300 hover:bg-gray-50 dark:hover:bg-slate-700 transition-colors"
          >
            <Edit3 size={14} />
            编辑
          </button>
          {article.url && onReparse && (
            <button
              onClick={handleReparse}
              className="w-full px-4 py-2.5 text-left text-sm flex items-center gap-2 text-slate-700 dark:text-slate-300 hover:bg-gray-50 dark:hover:bg-slate-700 transition-colors"
            >
              <RefreshCw size={14} />
              重新解析
            </button>
          )}
          <div className="border-t border-gray-100 dark:border-slate-700 my-1"></div>
          <button
            onClick={handleDelete}
            className="w-full px-4 py-2.5 text-left text-sm flex items-center gap-2 text-red-600 dark:text-red-400 hover:bg-red-50 dark:hover:bg-red-900/20 transition-colors"
          >
            <Trash2 size={14} />
            删除
          </button>
        </div>
      )}
    </>
  );
};

// 添加高亮样式
const style = document.createElement('style');
style.textContent = `
  mark {
    background-color: #fef08a;
    color: #854d0e;
    padding: 0 2px;
    border-radius: 2px;
    font-weight: 500;
  }
  .dark mark {
    background-color: rgba(250, 204, 21, 0.3);
    color: #fde047;
  }
`;
document.head.appendChild(style);

export default ArticleCard;
