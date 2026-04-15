
import React from 'react';
import { AlertTriangle } from 'lucide-react';

interface ConfirmDialogProps {
  isOpen: boolean;
  onClose: () => void;
  onConfirm: () => void;
  title: string;
  message: string;
  confirmLabel?: string;
  cancelLabel?: string;
  variant?: 'danger' | 'default';
}

const ConfirmDialog: React.FC<ConfirmDialogProps> = ({
  isOpen,
  onClose,
  onConfirm,
  title,
  message,
  confirmLabel = '确认',
  cancelLabel = '取消',
  variant = 'default'
}) => {
  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 bg-black/50 backdrop-blur-sm flex items-center justify-center p-4 z-[60] animate-in fade-in duration-200">
      <div className="bg-white dark:bg-slate-800 rounded-2xl w-full max-w-sm shadow-2xl overflow-hidden animate-in zoom-in-95 duration-200 transition-colors">
        <div className="p-6 text-center">
          <div className={`w-12 h-12 rounded-full mx-auto flex items-center justify-center mb-4 ${
            variant === 'danger' ? 'bg-red-50 dark:bg-red-900/20 text-red-500 dark:text-red-400' : 'bg-blue-50 dark:bg-blue-900/20 text-blue-500 dark:text-blue-400'
          }`}>
            <AlertTriangle size={24} />
          </div>
          
          <h3 className="text-lg font-bold text-gray-900 dark:text-white mb-2">{title}</h3>
          <p className="text-gray-500 dark:text-slate-400 text-sm mb-6 leading-relaxed">
            {message}
          </p>

          <div className="flex gap-3">
            <button
              onClick={onClose}
              className="flex-1 py-2.5 px-4 bg-gray-100 dark:bg-slate-700 hover:bg-gray-200 dark:hover:bg-slate-600 text-gray-700 dark:text-slate-200 font-medium rounded-xl transition-colors outline-none focus:ring-2 focus:ring-gray-300 dark:focus:ring-slate-500"
            >
              {cancelLabel}
            </button>
            <button
              onClick={() => {
                onConfirm();
                onClose();
              }}
              className={`flex-1 py-2.5 px-4 text-white font-medium rounded-xl shadow-lg transition-all transform active:scale-95 outline-none focus:ring-2 focus:ring-offset-2 ${
                variant === 'danger' 
                  ? 'bg-red-500 hover:bg-red-600 shadow-red-500/30 focus:ring-red-500' 
                  : 'bg-wechat hover:bg-wechat-dark shadow-wechat/30 focus:ring-wechat'
              }`}
            >
              {confirmLabel}
            </button>
          </div>
        </div>
      </div>
    </div>
  );
};

export default ConfirmDialog;
