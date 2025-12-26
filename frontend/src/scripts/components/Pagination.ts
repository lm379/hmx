import { defineComponent, computed, toRefs } from 'vue';

export default defineComponent({
  name: 'Pagination',
  props: {
    currentPage: {
      type: Number,
      required: true
    },
    totalItems: {
      type: Number,
      required: true
    },
    pageSize: {
      type: Number,
      default: 10
    }
  },
  emits: ['page-change'],
  setup(props, { emit }) {
    const { currentPage, totalItems, pageSize } = toRefs(props);

    const totalPages = computed(() => {
      return Math.ceil(totalItems.value / pageSize.value);
    });

    const handlePageChange = (newPage: number) => {
      if (newPage >= 1 && newPage <= totalPages.value && newPage !== currentPage.value) {
        emit('page-change', newPage);
      }
    };

    return {
      totalPages,
      handlePageChange
    };
  }
});
