import { defineComponent, ref, onMounted } from 'vue';
import { getOperas } from '../../api/opera';
import type { OperaListItem } from '../../types';
import { useRouter } from 'vue-router';
import Pagination from '../../components/Pagination.vue';
import { formatDate } from '../../utils/dateUtils';

export default defineComponent({
  name: 'RankingView',
  components: {
    Pagination
  },
  setup() {
    const router = useRouter();
    const rankings = ref<OperaListItem[]>([]);
    const loading = ref(true);
    const currentPage = ref(1);
    const pageSize = ref(20);
    const totalItems = ref(0);

    const fetchRankings = async (page = 1) => {
      try {
        loading.value = true;
        const data = await getOperas({
          page: page,
          page_size: pageSize.value,
        });
        if (data && data.data && data.data.list) {
          rankings.value = data.data.list;
          if (data.data.pagination) {
            totalItems.value = data.data.pagination.total;
          }
          currentPage.value = page;
        }
      } catch (error) {
        console.error('Failed to fetch rankings', error);
      } finally {
        loading.value = false;
      }
    };

    const navigateToVideo = (id: number) => {
      router.push(`/video/${id}`);
    };

    const handlePageChange = (newPage: number) => {
      fetchRankings(newPage);
      window.scrollTo({ top: 0, behavior: 'smooth' });
    };

    onMounted(() => {
      fetchRankings();
    });

    return {
      rankings,
      loading,
      navigateToVideo,
      currentPage,
      pageSize,
      totalItems,
      handlePageChange,
      formatDate
    };
  }
});
