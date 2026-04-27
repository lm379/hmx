import axios from 'axios';
import type { NewsDetail, NewsListItem } from '../types';

const API_BASE_URL = '/api/v1';

export interface NewsListResponse {
  code: number;
  msg: string;
  data: {
    list: NewsListItem[];
    pagination: {
      total: number;
      page: number;
      page_size: number;
    };
  };
}

export interface NewsDetailResponse {
  code: number;
  msg: string;
  data: NewsDetail;
}

export interface NewsListParams {
  page?: number;
  page_size?: number;
}

export interface SaveNewsPayload {
  title: string;
  summary: string;
  content: string;
  cover: string;
  source: string;
  author: string;
  is_published: boolean;
  published_at?: string | null;
}

export const getNews = async (params: NewsListParams = {}) => {
  const response = await axios.get<NewsListResponse>(`${API_BASE_URL}/news/`, { params });
  return response.data;
};

export const getNewsById = async (id: number) => {
  const response = await axios.get<NewsDetailResponse>(`${API_BASE_URL}/news/${id}`);
  return response.data;
};

export const getAdminNews = async (params: NewsListParams = {}) => {
  const response = await axios.get<NewsListResponse>(`${API_BASE_URL}/admin/news`, { params });
  return response.data;
};

export const getAdminNewsById = async (id: number) => {
  const response = await axios.get<NewsDetailResponse>(`${API_BASE_URL}/admin/news/${id}`);
  return response.data;
};

export const createNews = async (payload: SaveNewsPayload) => {
  const response = await axios.post<NewsDetailResponse>(`${API_BASE_URL}/admin/news`, payload);
  return response.data;
};

export const updateNews = async (id: number, payload: SaveNewsPayload) => {
  const response = await axios.put<NewsDetailResponse>(`${API_BASE_URL}/admin/news/${id}`, payload);
  return response.data;
};

export const deleteNews = async (id: number) => {
  const response = await axios.delete(`${API_BASE_URL}/admin/news/${id}`);
  return response.data;
};
