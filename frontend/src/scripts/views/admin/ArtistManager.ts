import { defineComponent, ref, onMounted, reactive } from 'vue';
import axios from 'axios';
import { ElMessage, ElMessageBox, type UploadFile, type UploadInstance, type UploadProps } from 'element-plus';

export default defineComponent({
  name: 'ArtistManager',
  setup() {
    const tableData = ref([]);
    const currentPage = ref(1);
    const pageSize = ref(10);
    const total = ref(0);
    const loading = ref(false);
    const avatarUploadRef = ref<UploadInstance>();

    const dialogVisible = ref(false);
    const isEdit = ref(false);
    const submitting = ref(false);

    const form = reactive({
      id: 0,
      name: '',
      bio: '',
      avatar: '',
    });

    const fetchData = async () => {
      loading.value = true;
      try {
        const res = await axios.get('/api/v1/artists', {
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

    onMounted(() => {
      fetchData();
    });

    const handleCreate = () => {
      isEdit.value = false;
      form.id = 0;
      form.name = '';
      form.bio = '';
      form.avatar = '';
      
      if(avatarUploadRef.value) avatarUploadRef.value.clearFiles();
      avatarFile.value = null;

      dialogVisible.value = true;
    };

    const handleEdit = (row: any) => {
      isEdit.value = true;
      form.id = row.artist_id;
      form.name = row.name;
      form.bio = row.bio;
      form.avatar = '';
      
      if(avatarUploadRef.value) avatarUploadRef.value.clearFiles();
      avatarFile.value = null;
      
      dialogVisible.value = true;
    };

    const handleDelete = (row: any) => {
      ElMessageBox.confirm('确定删除该艺术家吗？', '提示', {
        type: 'warning'
      }).then(async () => {
        try {
          await axios.delete(`/api/v1/admin/artists/${row.artist_id}`);
          ElMessage.success('删除成功');
          fetchData();
        } catch (e) {
          ElMessage.error('删除失败');
        }
      });
    };

    const avatarFile = ref<File | null>(null);

    const handleElFileChange = (file: UploadFile) => {
      if (file.raw) {
        avatarFile.value = file.raw;
      }
    };

    const handleElFileRemove = () => {
      avatarFile.value = null;
    };

    const handleExceed: UploadProps['onExceed'] = () => {
      ElMessage.warning(`限制选择 1 个文件，请先移除旧文件`);
    };

    const submitForm = async () => {
      submitting.value = true;
      try {
        let aPath = form.avatar;
        if (avatarFile.value) {
          const res = await axios.post('/api/v1/uploads/presign', {
            filename: avatarFile.value.name,
            content_type: avatarFile.value.type,
            upload_type: 'avatars'
          });
          const { upload_url, object_key } = res.data.data;
          await axios.put(upload_url, avatarFile.value, {
            headers: { 'Content-Type': avatarFile.value.type }
          });
          aPath = object_key;
        }

        const data = {
          name: form.name,
          bio: form.bio,
          avatar: aPath
        };

        if (isEdit.value) {
          await axios.put(`/api/v1/admin/artists/${form.id}`, data);
        } else {
          await axios.post('/api/v1/admin/artists', data);
        }

        ElMessage.success(isEdit.value ? '更新成功' : '创建成功');
        dialogVisible.value = false;
        fetchData();
      } catch (e) {
        ElMessage.error('操作失败');
      } finally {
        submitting.value = false;
        avatarFile.value = null;
      }
    };

    return {
      tableData,
      currentPage,
      pageSize,
      total,
      loading,
      avatarUploadRef,
      dialogVisible,
      isEdit,
      submitting,
      form,
      fetchData,
      handleCreate,
      handleEdit,
      handleDelete,
      handleElFileChange,
      handleElFileRemove,
      handleExceed,
      submitForm,
      avatarFile
    };
  }
});
