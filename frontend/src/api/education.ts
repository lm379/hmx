import axios from 'axios';
import type { EducationBookDetail, EducationBookListItem } from '../types';

const API_BASE_URL = '/api/v1';

export interface EducationListResponse {
  code: number;
  msg: string;
  data: {
    list: EducationBookListItem[];
    pagination: {
      total: number;
      page: number;
      page_size: number;
    };
  };
}

export interface EducationDetailResponse {
  code: number;
  msg: string;
  data: EducationBookDetail;
}

export interface EducationListParams {
  page?: number;
  page_size?: number;
}

export interface SaveEducationBookPayload {
  title: string;
  description: string;
  pdf_path: string;
  cover_path: string;
  is_published: boolean;
}

export const getEducationBooks = async (params: EducationListParams = {}) => {
  const response = await axios.get<EducationListResponse>(`${API_BASE_URL}/education/books`, { params });
  return response.data;
};

export const getEducationBookById = async (id: number) => {
  const response = await axios.get<EducationDetailResponse>(`${API_BASE_URL}/education/books/${id}`);
  return response.data;
};

export const getAdminEducationBooks = async (params: EducationListParams = {}) => {
  const response = await axios.get<EducationListResponse>(`${API_BASE_URL}/admin/education/books`, { params });
  return response.data;
};

export const getAdminEducationBookById = async (id: number) => {
  const response = await axios.get<EducationDetailResponse>(`${API_BASE_URL}/admin/education/books/${id}`);
  return response.data;
};

export const createEducationBook = async (payload: SaveEducationBookPayload) => {
  const response = await axios.post<EducationDetailResponse>(`${API_BASE_URL}/admin/education/books`, payload);
  return response.data;
};

export const updateEducationBook = async (id: number, payload: SaveEducationBookPayload) => {
  const response = await axios.put<EducationDetailResponse>(`${API_BASE_URL}/admin/education/books/${id}`, payload);
  return response.data;
};

export const deleteEducationBook = async (id: number) => {
  const response = await axios.delete(`${API_BASE_URL}/admin/education/books/${id}`);
  return response.data;
};
