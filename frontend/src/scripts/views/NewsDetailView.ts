import { computed, defineComponent, onMounted, ref } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { getNewsById } from '../../api/news';
import type { NewsDetail } from '../../types';
import { formatDateTime } from '../../utils/dateUtils';
import { renderMarkdown } from '../../utils/markdown';

export default defineComponent({
  name: 'NewsDetailView',
  setup() {
    const route = useRoute();
    const router = useRouter();
    const news = ref<NewsDetail | null>(null);
    const loading = ref(false);

    // 新闻正文由管理员维护，渲染前统一清洗，避免富文本内容带来脚本风险。
    const contentHtml = computed(() => {
      return renderMarkdown(news.value?.content || '');
    });

    const summaryHtml = computed(() => {
      return renderMarkdown(news.value?.summary || '');
    });

    const fetchDetail = async () => {
      const id = Number(route.params.id);
      if (!id) {
        router.replace('/404');
        return;
      }

      loading.value = true;
      try {
        const res = await getNewsById(id);
        news.value = res.data;
      } catch {
        news.value = null;
      } finally {
        loading.value = false;
      }
    };

    const formatNewsDate = (value: string) => formatDateTime(value);

    onMounted(() => {
      fetchDetail();
    });

    return {
      news,
      loading,
      contentHtml,
      summaryHtml,
      formatNewsDate
    };
  }
});
