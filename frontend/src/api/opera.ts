import axios from 'axios';
import type { OperaListResponse, OperaListItem } from '../types';

const API_BASE_URL = '/api/v1';

export interface GetOperasParams {
  page?: number;
  page_size?: number;
  sort?: string;
}

export const getOperas = async (params: GetOperasParams = {}) => {
  const queryParams = new URLSearchParams();
  if (params.page) queryParams.append('page', params.page.toString());
  if (params.page_size) queryParams.append('page_size', params.page_size.toString());
  if (params.sort) queryParams.append('sort', params.sort);

  const response = await axios.get<OperaListResponse>(
    `${API_BASE_URL}/operas?${queryParams.toString()}`
  );
  return response.data;
};

/**
 * 根据ID获取单个作品详情
 */
export const getOperaById = async (id: number): Promise<{ code: number; data: OperaListItem }> => {
  const response = await axios.get(`${API_BASE_URL}/operas/${id}`);
  return response.data;
};

/**
 * 批量获取作品详情
 */
export const getOperasByIds = async (ids: number[]): Promise<OperaListItem[]> => {
  const requests = ids.map(id => getOperaById(id));
  const responses = await Promise.allSettled(requests);
  
  return responses
    .filter((result): result is PromiseFulfilledResult<{ code: number; data: OperaListItem }> => 
      result.status === 'fulfilled' && result.value.code === 200
    )
    .map(result => result.value.data);
};

