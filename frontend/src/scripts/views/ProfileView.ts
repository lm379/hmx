import { defineComponent, ref, onMounted, watch } from 'vue';
import axios from 'axios';
import VideoCard from '../../components/VideoCard.vue';
import Pagination from '../../components/Pagination.vue';
import type { User, Opera } from '../../types';

export default defineComponent({
  name: 'ProfileView',
  components: {
    VideoCard,
    Pagination
  },
  setup() {
    const user = ref<User | null>(null);
    const loadingUser = ref(true);
    const activeTab = ref('history');
    const list = ref<Opera[]>([]);
    const loading = ref(false);
    
    const currentPage = ref(1);
    const pageSize = ref(10);
    const totalItems = ref(0);

    const fetchUser = async () => {
      try {
        loadingUser.value = true;
        const response = await axios.get('/api/v1/users/me');
        if (response.data && response.data.data) {
          user.value = response.data.data;
          fetchTabData();
        }
      } catch (error) {
        console.error("Failed to fetch user:", error);
      } finally {
        loadingUser.value = false;
      }
    };

    const fetchTabData = async (page = 1) => {
      try {
        loading.value = true;
        let endpoint = '';
        switch (activeTab.value) {
          case 'history': endpoint = '/api/v1/users/me/history'; break;
          case 'likes': endpoint = '/api/v1/users/me/likes'; break;
          case 'favorites': endpoint = '/api/v1/users/me/favorites'; break;
        }

        const response = await axios.get(endpoint, {
          params: {
            page,
            page_size: pageSize.value
          }
        });

        if (response.data && response.data.data) {
          list.value = response.data.data.list || [];
          if (response.data.data.pagination) {
            totalItems.value = response.data.data.pagination.total;
            currentPage.value = response.data.data.pagination.page;
          }
        }
      } catch (error) {
        console.error("Failed to fetch tab data:", error);
        list.value = [];
        totalItems.value = 0;
      } finally {
        loading.value = false;
      }
    };

    const handlePageChange = (newPage: number) => {
      fetchTabData(newPage);
      window.scrollTo({ top: 0, behavior: 'smooth' });
    };

    watch(activeTab, () => {
      currentPage.value = 1;
      fetchTabData();
    });

    onMounted(() => {
      fetchUser();
    });

    const formatTime = (time: string) => {
      if (!time) return '';
      const date = new Date(time);
      return date.toLocaleString('zh-CN');
    };

    return {
      user,
      loadingUser,
      activeTab,
      list,
      loading,
      currentPage,
      pageSize,
      totalItems,
      handlePageChange,
      formatTime
    };
  }
});
