import { defineComponent, ref, onMounted, type Ref } from 'vue';
import { useRouter } from 'vue-router';
import VideoCard from '../../components/VideoCard.vue';
import Pagination from '../../components/Pagination.vue';
import axios from 'axios';
import type { Opera } from '../../types';

export default defineComponent({
  name: 'HomeView',
  components: {
    VideoCard,
    Pagination
  },
  setup() {
    const router = useRouter();
    const operas: Ref<Opera[]> = ref([]);
    const loading = ref(true);
    const currentPage = ref(1);
    const pageSize = ref(15); // Grid usually fits 12 better (4x3)
    const totalItems = ref(0);

    const fetchOperas = async (page = 1) => {
      try {
        loading.value = true;
        const response = await axios.get('/api/v1/operas/', {
          params: {
            page: page,
            page_size: pageSize.value
          }
        });

        if (response.data && response.data.data) {
          operas.value = response.data.data.list || [];
          if (response.data.data.pagination) {
            totalItems.value = response.data.data.pagination.total;
            currentPage.value = response.data.data.pagination.page;
          }
        } else {
          console.error("Unexpected API response structure:", response.data);
        }
      } catch (error) {
        console.error("Failed to fetch operas:", error);
      } finally {
        loading.value = false;
      }
    };

    const navigateToVideo = (id: number) => {
      router.push(`/video/${id}`);
    };

    const handlePageChange = (newPage: number) => {
      fetchOperas(newPage);
      window.scrollTo({ top: 0, behavior: 'smooth' });
    };

    onMounted(() => {
      fetchOperas();
    });

    return {
      operas,
      loading,
      currentPage,
      pageSize,
      totalItems,
      fetchOperas,
      navigateToVideo,
      handlePageChange
    };
  }
});
