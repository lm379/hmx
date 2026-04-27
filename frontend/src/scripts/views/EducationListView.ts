import { defineComponent, onMounted, ref, type Ref } from 'vue';
import { useRouter } from 'vue-router';
import Pagination from '../../components/Pagination.vue';
import { getEducationBooks } from '../../api/education';
import type { EducationBookListItem } from '../../types';

export default defineComponent({
  name: 'EducationListView',
  components: {
    Pagination
  },
  setup() {
    const router = useRouter();
    const books: Ref<EducationBookListItem[]> = ref([]);
    const loading = ref(false);
    const currentPage = ref(1);
    const pageSize = ref(12);
    const totalItems = ref(0);

    // 前台只加载已发布 PDF，草稿资源只在管理后台可见。
    const fetchBooks = async (page = 1) => {
      loading.value = true;
      try {
        const res = await getEducationBooks({ page, page_size: pageSize.value });
        books.value = res.data.list || [];
        totalItems.value = res.data.pagination.total || 0;
        currentPage.value = res.data.pagination.page || page;
      } finally {
        loading.value = false;
      }
    };

    const handlePageChange = (page: number) => {
      fetchBooks(page);
      window.scrollTo({ top: 0, behavior: 'smooth' });
    };

    const goReader = (id: number) => {
      router.push(`/education/${id}`);
    };

    onMounted(() => {
      fetchBooks();
    });

    return {
      books,
      loading,
      currentPage,
      pageSize,
      totalItems,
      handlePageChange,
      goReader
    };
  }
});
