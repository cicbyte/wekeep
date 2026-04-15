
export interface Article {
  id: string;
  title: string;
  author: string;        // 作者名称（用于显示）
  authorId?: number;     // 作者ID（用于 API）
  url?: string;
  summary?: string;
  content?: string; // Markdown content
  tags: string[];
  dateAdded: number; // timestamp
  // 搜索高亮相关字段
  formattedTitle?: string;    // 高亮后的标题
  formattedSummary?: string;  // 高亮后的摘要
  contextSnippet?: string;    // 正文匹配上下文
  matchFields?: string[];     // 匹配的字段列表
}

export type ViewMode = 'dashboard' | 'authors' | 'list' | 'settings';

export interface ParseResult {
  title: string;
  author: string;
  summary: string;
  content: string;
  tags: string[];
  url?: string;
}

export const SAMPLE_ARTICLES: Article[] = [
  {
    id: '1',
    title: 'Deep Dive into React 19',
    author: 'Frontend Weekly',
    url: 'https://mp.weixin.qq.com/s/example1',
    summary: 'An extensive look at the new features in React 19, including the compiler and server components.',
    content: '## React 19 Features\n\nReact 19 brings the compiler...',
    tags: ['React', 'Tech', 'Frontend'],
    dateAdded: Date.now() - 100000000
  },
  {
    id: '2',
    title: 'The Future of AI in 2025',
    author: 'TechDaily',
    url: 'https://mp.weixin.qq.com/s/example2',
    summary: 'Predictions about Gemini and other LLMs shaping the future software landscape.',
    content: '## AI Predictions\n\nBy 2025, we expect...',
    tags: ['AI', 'Future', 'Gemini'],
    dateAdded: Date.now() - 50000000
  },
  {
    id: '3',
    title: 'Tailwind CSS Best Practices',
    author: 'Frontend Weekly',
    summary: 'How to organize your utility classes and avoid clutter.',
    content: '## Organizing Utilities\n\nAvoid long strings...',
    tags: ['CSS', 'Design'],
    dateAdded: Date.now() - 20000000
  }
];
