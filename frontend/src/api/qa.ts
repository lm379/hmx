import axios from 'axios';

const API_BASE_URL = '/api/v1';

export interface Source {
  chunk_id: string;
  doc_id: string;
  doc_title: string;
  chunk_text: string;
  score: number;
}

export interface AskResponse {
  qa_id: string;
  answer: string;
  sources: Source[];
  session_id: string;
}

export interface QAHistoryItem {
  qa_id: string;
  query: string;
  answer: string;
  sources: Source[];
  created_at: string;
}

export interface QAHistoryResponse {
  code: number;
  data: {
    list: QAHistoryItem[];
    pagination: {
      total: number;
      page: number;
      page_size: number;
    };
  };
}

export const askQuestion = async (
  query: string,
  sessionId?: string,
  operaId?: number
): Promise<{ code: number; data: AskResponse }> => {
  const response = await axios.post(`${API_BASE_URL}/qa/ask`, {
    query,
    session_id: sessionId,
    opera_id: operaId
  });
  return response.data;
};

export const getQAHistory = async (
  page = 1,
  pageSize = 20,
  sessionId?: string
): Promise<QAHistoryResponse> => {
  const params: Record<string, string | number> = { page, page_size: pageSize };
  if (sessionId) params.session_id = sessionId;
  const response = await axios.get(`${API_BASE_URL}/qa/history`, { params });
  return response.data;
};

export const submitFeedback = async (
  qaId: string,
  rating: number,
  comments?: string
): Promise<void> => {
  await axios.post(`${API_BASE_URL}/qa/${qaId}/feedback`, { rating, comments });
};
