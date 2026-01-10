import axios from 'axios';

const API_BASE_URL = '/api/v1';

// 推荐参数接口
export interface RecommendationParams {
  limit?: number;
  similarity_weight?: number;
  interest_weight?: number;
}

// 相似视频参数接口
export interface SimilarOperasParams {
  limit?: number;
}

// 获取个性化推荐
export const getRecommendations = async (params: RecommendationParams = {}) => {
  const queryParams = new URLSearchParams();
  
  if (params.limit) queryParams.append('limit', params.limit.toString());
  if (params.similarity_weight !== undefined) {
    queryParams.append('similarity_weight', params.similarity_weight.toString());
  }
  if (params.interest_weight !== undefined) {
    queryParams.append('interest_weight', params.interest_weight.toString());
  }

  const response = await axios.get(
    `${API_BASE_URL}/recommendations?${queryParams.toString()}`
  );
  return response.data;
};

// 获取相似作品推荐
export const getSimilarOperas = async (operaId: number, params: SimilarOperasParams = {}) => {
  const queryParams = new URLSearchParams();
  
  if (params.limit) queryParams.append('limit', params.limit.toString());

  const response = await axios.get(
    `${API_BASE_URL}/operas/${operaId}/similar?${queryParams.toString()}`
  );
  return response.data;
};

// 管理员：生成单个向量
export const generateEmbedding = async (operaId: number) => {
  const response = await axios.post(
    `${API_BASE_URL}/admin/operas/${operaId}/embedding`
  );
  return response.data;
};

// 管理员：批量生成向量
export const batchGenerateEmbeddings = async () => {
  const response = await axios.post(
    `${API_BASE_URL}/admin/operas/batch-embedding`
  );
  return response.data;
};
