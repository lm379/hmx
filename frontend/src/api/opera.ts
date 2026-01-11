import axios from 'axios';
import type { OperaListResponse } from '../types';

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
