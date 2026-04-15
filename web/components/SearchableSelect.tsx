import React, { useState, useRef, useEffect } from 'react';
import { Search, ChevronDown, X, User } from 'lucide-react';

interface Option {
  value: string;
  label: string;
  count?: number;
  avatar?: string;
}

interface SearchableSelectProps {
  options: Option[];
  value: string;
  onChange: (value: string) => void;
  placeholder?: string;
  emptyText?: string;
}

const SearchableSelect: React.FC<SearchableSelectProps> = ({
  options,
  value,
  onChange,
  placeholder = '请选择',
  emptyText = '暂无选项'
}) => {
  const [isOpen, setIsOpen] = useState(false);
  const [searchQuery, setSearchQuery] = useState('');
  const [highlightedIndex, setHighlightedIndex] = useState(-1);
  const containerRef = useRef<HTMLDivElement>(null);
  const listRef = useRef<HTMLDivElement>(null);

  const selectedOption = options.find(opt => opt.value === value);

  const filteredOptions = options.filter(opt =>
    opt.label.toLowerCase().includes(searchQuery.toLowerCase())
  );

  // 所有可选项（包含"全部作者"）
  const allOptions = [
    { value: '', label: '全部作者', isAll: true },
    ...filteredOptions
  ];

  // 点击外部关闭
  useEffect(() => {
    const handleClickOutside = (e: MouseEvent) => {
      if (containerRef.current && !containerRef.current.contains(e.target as Node)) {
        setIsOpen(false);
        setSearchQuery('');
        setHighlightedIndex(-1);
      }
    };

    if (isOpen) {
      document.addEventListener('mousedown', handleClickOutside);
      return () => document.removeEventListener('mousedown', handleClickOutside);
    }
  }, [isOpen]);

  // 键盘导航
  useEffect(() => {
    if (!isOpen) return;

    const handleKeyDown = (e: KeyboardEvent) => {
      switch (e.key) {
        case 'Escape':
          setIsOpen(false);
          setSearchQuery('');
          setHighlightedIndex(-1);
          break;
        case 'ArrowDown':
          e.preventDefault();
          setHighlightedIndex(prev =>
            prev < allOptions.length - 1 ? prev + 1 : 0
          );
          break;
        case 'ArrowUp':
          e.preventDefault();
          setHighlightedIndex(prev =>
            prev > 0 ? prev - 1 : allOptions.length - 1
          );
          break;
        case 'Enter':
          e.preventDefault();
          if (highlightedIndex >= 0 && highlightedIndex < allOptions.length) {
            handleSelect(allOptions[highlightedIndex].value);
          }
          break;
      }
    };

    document.addEventListener('keydown', handleKeyDown);
    return () => document.removeEventListener('keydown', handleKeyDown);
  }, [isOpen, highlightedIndex, allOptions]);

  // 滚动到高亮项
  useEffect(() => {
    if (highlightedIndex >= 0 && listRef.current) {
      const items = listRef.current.querySelectorAll('[data-option-index]');
      const targetItem = items[highlightedIndex];
      if (targetItem) {
        targetItem.scrollIntoView({ block: 'nearest' });
      }
    }
  }, [highlightedIndex]);

  // 重置高亮
  useEffect(() => {
    setHighlightedIndex(-1);
  }, [searchQuery]);

  const handleSelect = (optionValue: string) => {
    onChange(optionValue);
    setIsOpen(false);
    setSearchQuery('');
    setHighlightedIndex(-1);
  };

  const handleClear = (e: React.MouseEvent) => {
    e.stopPropagation();
    onChange('');
  };

  const getTotalCount = () => options.reduce((sum, opt) => sum + (opt.count || 0), 0);

  return (
    <div ref={containerRef} className="relative min-w-[160px]">
      {/* 触发按钮 */}
      <button
        type="button"
        onClick={() => setIsOpen(!isOpen)}
        className={`w-full flex items-center gap-2 px-3 py-2 rounded-xl text-sm transition-all outline-none ${
          isOpen
            ? 'bg-white dark:bg-slate-900 ring-2 ring-wechat/50'
            : 'bg-gray-100 dark:bg-slate-800 hover:bg-gray-200 dark:hover:bg-slate-700'
        } ${selectedOption ? 'text-slate-700 dark:text-slate-200' : 'text-gray-500 dark:text-slate-400'}`}
      >
        <User size={16} className="text-gray-400 dark:text-slate-500 shrink-0" />
        <span className="flex-1 text-left truncate">
          {selectedOption ? selectedOption.label : placeholder}
        </span>
        {selectedOption ? (
          <button
            type="button"
            onClick={handleClear}
            className="text-gray-400 hover:text-gray-600 dark:hover:text-slate-300 transition-colors"
          >
            <X size={14} />
          </button>
        ) : (
          <ChevronDown size={16} className={`text-gray-400 transition-transform ${isOpen ? 'rotate-180' : ''}`} />
        )}
      </button>

      {/* 下拉面板 */}
      {isOpen && (
        <div className="absolute top-full left-0 right-0 mt-1 bg-white dark:bg-slate-800 rounded-xl shadow-lg border border-gray-200 dark:border-slate-700 z-50 animate-in fade-in slide-in-from-top-2 duration-200">
          {/* 搜索框 */}
          <div className="p-2 border-b border-gray-100 dark:border-slate-700">
            <div className="relative">
              <Search size={14} className="absolute left-2.5 top-1/2 -translate-y-1/2 text-gray-400" />
              <input
                type="text"
                placeholder="搜索作者..."
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                className="w-full bg-gray-50 dark:bg-slate-900 border border-gray-200 dark:border-slate-700 rounded-lg py-1.5 pl-8 pr-3 text-sm outline-none focus:ring-2 focus:ring-wechat/30 focus:border-wechat dark:text-slate-200 dark:placeholder-slate-500"
                autoFocus
              />
            </div>
          </div>

          {/* 选项列表 */}
          <div ref={listRef} className="max-h-60 overflow-y-auto py-1 scrollbar-hide [-ms-overflow-style:none] [scrollbar-width:none] [&::-webkit-scrollbar]:hidden">
            {/* 全部选项 */}
            <button
              type="button"
              data-option-index="0"
              onClick={() => handleSelect('')}
              onMouseEnter={() => setHighlightedIndex(0)}
              className={`w-full flex items-center gap-2 px-3 py-2 text-sm transition-colors ${
                highlightedIndex === 0
                  ? 'bg-gray-100 dark:bg-slate-700'
                  : ''
              } ${
                !value
                  ? 'text-wechat-dark dark:text-wechat-light'
                  : 'text-slate-700 dark:text-slate-300'
              }`}
            >
              <div className="w-6 h-6 rounded-full bg-gray-200 dark:bg-slate-600 flex items-center justify-center">
                <User size={12} className="text-gray-500 dark:text-slate-400" />
              </div>
              <span className="flex-1 text-left">全部作者</span>
              <span className="text-xs text-gray-400 dark:text-slate-500">
                {getTotalCount()}
              </span>
            </button>

            {/* 作者列表 */}
            {filteredOptions.length > 0 ? (
              filteredOptions.map((option, index) => (
                <button
                  key={option.value}
                  type="button"
                  data-option-index={index + 1}
                  onClick={() => handleSelect(option.value)}
                  onMouseEnter={() => setHighlightedIndex(index + 1)}
                  className={`w-full flex items-center gap-2 px-3 py-2 text-sm transition-colors ${
                    highlightedIndex === index + 1
                      ? 'bg-gray-100 dark:bg-slate-700'
                      : ''
                  } ${
                    value === option.value
                      ? 'text-wechat-dark dark:text-wechat-light'
                      : 'text-slate-700 dark:text-slate-300'
                  }`}
                >
                  <div className="w-6 h-6 rounded-full bg-gradient-to-br from-wechat to-wechat-dark flex items-center justify-center text-white text-xs font-bold">
                    {option.avatar ? (
                      <img src={option.avatar} alt={option.label} className="w-full h-full rounded-full object-cover" />
                    ) : (
                      option.label.charAt(0).toUpperCase()
                    )}
                  </div>
                  <span className="flex-1 text-left truncate">{option.label}</span>
                  {option.count !== undefined && (
                    <span className="text-xs text-gray-400 dark:text-slate-500">{option.count}</span>
                  )}
                </button>
              ))
            ) : (
              <div className="px-3 py-4 text-center text-sm text-gray-400 dark:text-slate-500">
                {emptyText}
              </div>
            )}
          </div>
        </div>
      )}
    </div>
  );
};

export default SearchableSelect;
