import { marked } from 'marked';
import DOMPurify from 'dompurify';

marked.setOptions({ breaks: true });

// 管理端内容会以 Markdown 形式录入；所有前台 HTML 输出都必须先清洗。
export function renderMarkdown(text: string): string {
  const raw = marked.parse(text || '') as string;
  return DOMPurify.sanitize(raw);
}
