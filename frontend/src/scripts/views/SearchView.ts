import { defineComponent, ref, onMounted, watch } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { ElMessage } from 'element-plus';
import { Loading } from '@element-plus/icons-vue';
import { searchOperas } from '../../api/search';
import { getOperasByIds } from '../../api/opera';
import type { OperaListItem } from '../../types';
import VideoCard from '../../components/VideoCard.vue';
import Pagination from '../../components/Pagination.vue';

export default defineComponent({
  name: 'SearchView',
  components: {
    VideoCard,
    Pagination,
    Loading
  },
  setup() {
    const route = useRoute();
    const router = useRouter();

    const searchQuery = ref('');
    const loading = ref(false);
    const hasSearched = ref(false);
    const currentPage = ref(1);
    const pageSize = ref(24);
    const totalResults = ref(0);
    const processingTime = ref(0);

    const operaResults = ref<OperaListItem[]>([]);

    // 执行搜索
    const handleSearch = async () => {
      if (!searchQuery.value.trim()) {
        return;
      }

      loading.value = true;
      hasSearched.value = true;

      try {
        // 第一步：搜索获取ID列表
        const searchResponse = await searchOperas({
          query: searchQuery.value,
          page: currentPage.value,
          page_size: pageSize.value,
        });

        // SearchOperasResponse: { code: number, data: { results..., total... } }
        totalResults.value = searchResponse.data.total;
        processingTime.value = searchResponse.data.processing_time_ms;

        // 第二步：根据ID批量获取完整信息
        if (searchResponse.data.results.length > 0) {
          const operaIds = searchResponse.data.results.map(item => item.opera_id);
          operaResults.value = await getOperasByIds(operaIds);
        } else {
          operaResults.value = [];
        }

        // 更新 URL
        router.push({
          query: {
            q: searchQuery.value,
            page: currentPage.value,
          },
        });
      } catch (error: any) {
        console.error('搜索失败:', error);
        ElMessage.error(error.response?.data?.msg || '搜索失败，请稍后重试');
      } finally {
        loading.value = false;
      }
    };

    // 从 URL 参数初始化
    onMounted(() => {
      const query = route.query.q as string;
      
      if (query) {
        searchQuery.value = query;
        const page = parseInt(route.query.page as string);
        if (!isNaN(page)) {
          currentPage.value = page;
        }
        handleSearch();
      }
    });

    // 监听路由变化，当搜索关键词改变时重新搜索
    watch(
      () => route.query.q,
      (newQuery) => {
        if (newQuery && newQuery !== searchQuery.value) {
          searchQuery.value = newQuery as string;
          currentPage.value = 1; // 重置页码
          handleSearch();
        }
      }
    );

    // 页码变化
    const handlePageChange = (page: number) => {
      currentPage.value = page;
      handleSearch();
      window.scrollTo({ top: 0, behavior: 'smooth' });
    };

    return {
      searchQuery,
      loading,
      hasSearched,
      currentPage,
      pageSize,
      totalResults,
      processingTime,
      operaResults,
      handleSearch,
      handlePageChange
    };
  }
});
