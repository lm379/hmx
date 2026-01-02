import { defineComponent, ref, onMounted } from 'vue';
import axios from 'axios';
import { ElMessage } from 'element-plus';

export default defineComponent({
  name: 'UserManager',
  setup() {
    const tableData = ref([]);
    const currentPage = ref(1);
    const pageSize = ref(10);
    const total = ref(0);
    const loading = ref(false);

    const fetchData = async () => {
      loading.value = true;
      try {
        const res = await axios.get('/api/v1/admin/users', {
          params: { page: currentPage.value, page_size: pageSize.value }
        });
        tableData.value = res.data.data.list;
        total.value = res.data.data.pagination.total;
      } catch (err) {
        ElMessage.error('获取列表失败');
      } finally {
        loading.value = false;
      }
    };

    const toggleRole = async (row: any) => {
      const newRole = row.role === 'User' ? 'Administrator' : 'User';
      try {
        await axios.put(`/api/v1/admin/users/${row.user_id}/role`, {
          role: newRole
        });
        ElMessage.success('角色更新成功');
        fetchData();
      } catch (e) {
        ElMessage.error('更新失败');
      }
    };

    onMounted(() => {
      fetchData();
    });

    return {
      tableData,
      currentPage,
      pageSize,
      total,
      loading,
      fetchData,
      toggleRole
    };
  }
});
