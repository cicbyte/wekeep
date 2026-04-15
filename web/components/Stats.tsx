
import React from 'react';
import { BarChart, Bar, XAxis, YAxis, Tooltip, ResponsiveContainer, Cell } from 'recharts';
import { Author } from '../services/apiService';

interface StatsProps {
  authors: Author[];
  isDarkMode?: boolean;
}

const Stats: React.FC<StatsProps> = ({ authors, isDarkMode }) => {
  // 基于作者列表的 articleCount 计算
  const data = authors
    .map(author => ({ name: author.name, count: author.articleCount }))
    .sort((a, b) => b.count - a.count)
    .slice(0, 5); // Top 5 authors

  if (data.length === 0) return null;

  return (
    <div className="bg-white dark:bg-slate-800 p-6 rounded-2xl shadow-sm border border-gray-100 dark:border-slate-700 mb-8 transition-colors">
      <h3 className="text-lg font-bold text-gray-800 dark:text-slate-100 mb-4">热门作者</h3>
      <div className="h-[200px] w-full">
        <ResponsiveContainer width="100%" height="100%">
          <BarChart data={data} layout="vertical" margin={{ left: 20, right: 20 }}>
            <XAxis type="number" hide />
            <YAxis
              dataKey="name"
              type="category"
              width={100}
              tick={{ fontSize: 12, fill: '#94a3b8' }}
              axisLine={false}
              tickLine={false}
            />
            <Tooltip
              cursor={{ fill: 'transparent' }}
              contentStyle={{
                borderRadius: '8px',
                border: 'none',
                boxShadow: '0 4px 6px -1px rgb(0 0 0 / 0.1)',
                backgroundColor: isDarkMode ? '#1e293b' : '#fff',
                color: isDarkMode ? '#fff' : '#334155'
              }}
            />
            <Bar dataKey="count" radius={[0, 4, 4, 0]} barSize={20}>
               {data.map((entry, index) => (
                  <Cell key={`cell-${index}`} fill="#07C160" />
                ))}
            </Bar>
          </BarChart>
        </ResponsiveContainer>
      </div>
    </div>
  );
};

export default Stats;
