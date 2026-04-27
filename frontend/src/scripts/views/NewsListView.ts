import { defineComponent, onMounted, ref, type Ref } from 'vue';
import { useRouter } from 'vue-router';
import Pagination from '../../components/Pagination.vue';
import { getNews } from '../../api/news';
import type { NewsListItem } from '../../types';
import { formatDate } from '../../utils/dateUtils';
import { renderMarkdown } from '../../utils/markdown';

export default defineComponent({
  name: 'NewsListView',
  components: {
    Pagination
  },
  setup() {
    const router = useRouter();
    const newsList: Ref<NewsListItem[]> = ref([]);
    const loading = ref(false);
    const currentPage = ref(1);
    const pageSize = ref(9);
    const totalItems = ref(0);

    // 前台只消费已发布新闻，后台草稿不会从公开接口返回。
    const fetchNews = async (page = 1) => {
      loading.value = true;
      try {
        const res = await getNews({ page, page_size: pageSize.value });
        newsList.value = res.data.list || [];
        totalItems.value = res.data.pagination.total || 0;
        currentPage.value = res.data.pagination.page || page;
      } finally {
        loading.value = false;
      }
    };

    const handlePageChange = (page: number) => {
      fetchNews(page);
      window.scrollTo({ top: 0, behavior: 'smooth' });
    };

    const goDetail = (id: number) => {
      router.push(`/news/${id}`);
    };

    const formatNewsDate = (value: string) => formatDate(value);

    const renderSummary = (summary: string) => {
      return renderMarkdown(summary || '点击查看完整资讯内容');
    };

    onMounted(() => {
      fetchNews();
    });

    return {
      newsList,
      loading,
      currentPage,
      pageSize,
      totalItems,
      handlePageChange,
      goDetail,
      formatNewsDate,
      renderSummary
    };
  }
});
