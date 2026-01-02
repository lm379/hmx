import { defineStore } from 'pinia';
import { ref } from 'vue';
import axios from 'axios';
import type { User } from '../types';

export const useAuthStore = defineStore('auth', () => {
  const isLoggedIn = ref(false);
  const user = ref<User | null>(null);
  const loading = ref(false);

  // 设置 axios 拦截器处理 token 刷新
  const setupAxiosInterceptors = () => {
    // 请求拦截器：自动添加 access token
    axios.interceptors.request.use(
      (config) => {
        const accessToken = localStorage.getItem('access_token');
        if (accessToken) {
          config.headers.Authorization = `Bearer ${accessToken}`;
        }
        return config;
      },
      (error) => Promise.reject(error)
    );

    // 响应拦截器：处理 401 错误，自动刷新 token
    axios.interceptors.response.use(
      (response) => response,
      async (error) => {
        const originalRequest = error.config;

        // 跳过登录、登出、注册等不需要处理 401 的请求
        const skipUrls = ['/api/v1/auth/login', '/api/v1/auth/logout', '/api/v1/auth/register', '/api/v1/auth/refresh'];
        if (skipUrls.some(url => originalRequest.url?.includes(url))) {
          return Promise.reject(error);
        }

        // 如果是 401 错误且不是刷新 token 请求，尝试刷新
        if (error.response?.status === 401 && !originalRequest._retry) {
          originalRequest._retry = true;

          const refreshToken = localStorage.getItem('refresh_token');
          if (refreshToken) {
            try {
              const response = await axios.post('/api/v1/auth/refresh', {
                refresh_token: refreshToken
              });

              if (response.data?.data) {
                const { access_token, refresh_token } = response.data.data;
                localStorage.setItem('access_token', access_token);
                localStorage.setItem('refresh_token', refresh_token);

                // 重试原始请求
                originalRequest.headers.Authorization = `Bearer ${access_token}`;
                return axios(originalRequest);
              }
            } catch (refreshError) {
              // 刷新失败，清除所有 token
              logout();
              return Promise.reject(refreshError);
            }
          } else {
            logout();
          }
        }

        return Promise.reject(error);
      }
    );
  };

  const checkLoginStatus = async () => {
    const accessToken = localStorage.getItem('access_token');
    if (accessToken) {
      try {
        loading.value = true;
        const response = await axios.get('/api/v1/users/me');
        if (response.data && response.data.data) {
          user.value = response.data.data;
          isLoggedIn.value = true;
        }
      } catch (error) {
        console.error("Failed to fetch user info:", error);
        logout();
      } finally {
        loading.value = false;
      }
    } else {
      isLoggedIn.value = false;
      user.value = null;
    }
  };

  const login = async (accessToken: string, refreshToken: string) => {
    localStorage.setItem('access_token', accessToken);
    localStorage.setItem('refresh_token', refreshToken);
    await checkLoginStatus();
  };

  const logout = async () => {
    try {
      // 调用后端登出接口
      await axios.post('/api/v1/auth/logout');
    } catch (error) {
      console.error("Logout error:", error);
    } finally {
      localStorage.removeItem('access_token');
      localStorage.removeItem('refresh_token');
      delete axios.defaults.headers.common['Authorization'];
      isLoggedIn.value = false;
      user.value = null;
    }
  };

  // 初始化拦截器
  setupAxiosInterceptors();

  return {
    isLoggedIn,
    user,
    loading,
    checkLoginStatus,
    login,
    logout
  };
});
