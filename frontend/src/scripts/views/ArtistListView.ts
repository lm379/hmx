import { defineComponent, ref, onMounted, type Ref } from 'vue';
import axios from 'axios';
import Pagination from '../../components/Pagination.vue';
import type { Artist } from '../../types';

export default defineComponent({
  name: 'ArtistListView',
  components: {
    Pagination
  },
  setup() {
    const artists: Ref<Artist[]> = ref([]);
    const loading = ref(true);
    const currentPage = ref(1);
    const pageSize = ref(10);
    const totalItems = ref(0);

    const fetchArtists = async (page = 1) => {
      try {
        loading.value = true;
        const response = await axios.get('/api/v1/artists/', {
            params: {
                page: page,
                page_size: pageSize.value
            }
        });
        if (response.data && response.data.data) {
          artists.value = response.data.data.list || [];
          if (response.data.data.pagination) {
            totalItems.value = response.data.data.pagination.total;
            currentPage.value = response.data.data.pagination.page;
          }
        }
      } catch (error) {
        console.error("Failed to fetch artists:", error);
      } finally {
        loading.value = false;
      }
    };

    const handlePageChange = (newPage: number) => {
      fetchArtists(newPage);
      window.scrollTo({ top: 0, behavior: 'smooth' });
    };

    onMounted(() => {
      fetchArtists();
    });

    return {
      artists,
      loading,
      currentPage,
      pageSize,
      totalItems,
      handlePageChange
    };
  }
});
