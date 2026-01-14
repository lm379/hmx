import axios from 'axios';
import type { OperaSearchResult, ArtistSearchResult } from '../types';

const API_BASE_URL = '/api/v1';

// 搜索作品请求参数
export interface SearchOperasParams {
  query?: string;
  page?: number;
  page_size?: number;
  include_hidden?: boolean;
  sort_by?: string;
}

// 搜索作品响应
export interface SearchOperasResponse {
  code: number;
  data: {
    results: OperaSearchResult[];
    total: number;
    page: number;
    page_size: number;
    total_pages: number;
    processing_time_ms: number;
  };
}

// 搜索艺术家请求参数
export interface SearchArtistsParams {
  query?: string;
  page?: number;
  page_size?: number;
  sort_by?: string;
}

// 搜索艺术家响应
export interface SearchArtistsResponse {
  code: number;
  data: {
    results: ArtistSearchResult[];
    total: number;
    page: number;
    page_size: number;
    total_pages: number;
    processing_time_ms: number;
  };
}

// 搜索统计响应
export interface SearchStatsResponse {
  code: number;
  data: {
    operas?: {
      numberOfDocuments: number;
      isIndexing: boolean;
      fieldDistribution?: Record<string, number>;
    };
    artists?: {
      numberOfDocuments: number;
      isIndexing: boolean;
      fieldDistribution?: Record<string, number>;
    };
  };
}

// 重新索引响应
export interface ReindexResponse {
  code: number;
  data: {
    message: string;
  };
}

/**
 * 搜索作品
 */
export const searchOperas = async (params: SearchOperasParams = {}): Promise<SearchOperasResponse> => {
  const queryParams = new URLSearchParams();
  if (params.query) queryParams.append('query', params.query);
  if (params.page) queryParams.append('page', params.page.toString());
  if (params.page_size) queryParams.append('page_size', params.page_size.toString());
  if (params.include_hidden !== undefined) queryParams.append('include_hidden', params.include_hidden.toString());
  if (params.sort_by) queryParams.append('sort_by', params.sort_by);

  const response = await axios.get<SearchOperasResponse>(
    `${API_BASE_URL}/search/operas?${queryParams.toString()}`
  );
  return response.data;
};

/**
 * 搜索艺术家
 */
export const searchArtists = async (params: SearchArtistsParams = {}): Promise<SearchArtistsResponse> => {
  const queryParams = new URLSearchParams();
  if (params.query) queryParams.append('query', params.query);
  if (params.page) queryParams.append('page', params.page.toString());
  if (params.page_size) queryParams.append('page_size', params.page_size.toString());
  if (params.sort_by) queryParams.append('sort_by', params.sort_by);

  const response = await axios.get<SearchArtistsResponse>(
    `${API_BASE_URL}/search/artists?${queryParams.toString()}`
  );
  return response.data;
};

/**
 * 重新索引所有作品（管理员）
 */
export const reindexOperas = async (): Promise<ReindexResponse> => {
  const response = await axios.post<ReindexResponse>(
    `${API_BASE_URL}/admin/search/reindex/operas`
  );
  return response.data;
};

/**
 * 重新索引所有艺术家（管理员）
 */
export const reindexArtists = async (): Promise<ReindexResponse> => {
  const response = await axios.post<ReindexResponse>(
    `${API_BASE_URL}/admin/search/reindex/artists`
  );
  return response.data;
};

/**
 * 获取搜索统计（管理员）
 */
export const getSearchStats = async (): Promise<SearchStatsResponse> => {
  const response = await axios.get<SearchStatsResponse>(
    `${API_BASE_URL}/admin/search/stats`
  );
  return response.data;
};
