import axios from 'axios';
import type { ArtistListItem } from '../types';

const API_BASE_URL = '/api/v1';

/**
 * 根据ID获取单个艺术家详情
 */
export const getArtistById = async (id: number): Promise<{ code: number; data: ArtistListItem }> => {
  const response = await axios.get(`${API_BASE_URL}/artists/${id}`);
  return response.data;
};

/**
 * 批量获取艺术家详情
 */
export const getArtistsByIds = async (ids: number[]): Promise<ArtistListItem[]> => {
  const requests = ids.map(id => getArtistById(id));
  const responses = await Promise.allSettled(requests);
  
  return responses
    .filter((result): result is PromiseFulfilledResult<{ code: number; data: ArtistListItem }> => 
      result.status === 'fulfilled' && result.value.code === 200
    )
    .map(result => result.value.data);
};
