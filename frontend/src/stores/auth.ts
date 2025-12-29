import { defineStore } from 'pinia';
import { ref } from 'vue';
import axios from 'axios';
import type { User } from '../types';

export const useAuthStore = defineStore('auth', () => {
  const isLoggedIn = ref(false);
  const user = ref<User | null>(null);
  const loading = ref(false);

  const checkLoginStatus = async () => {
    const token = localStorage.getItem('token');
    if (token) {
      axios.defaults.headers.common['Authorization'] = `Bearer ${token}`;
      
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

  const login = async (token: string) => {
    localStorage.setItem('token', token);
    axios.defaults.headers.common['Authorization'] = `Bearer ${token}`;
    await checkLoginStatus();
  };

  const logout = () => {
    localStorage.removeItem('token');
    delete axios.defaults.headers.common['Authorization'];
    isLoggedIn.value = false;
    user.value = null;
  };

  return {
    isLoggedIn,
    user,
    loading,
    checkLoginStatus,
    login,
    logout
  };
});
