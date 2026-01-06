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
    const uploadProgress = ref(0);

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

      if (avatarUploadRef.value) avatarUploadRef.value.clearFiles();
      avatarFile.value = null;

      dialogVisible.value = true;
    };

    const handleEdit = (row: any) => {
      isEdit.value = true;
      form.id = row.artist_id;
      form.name = row.name;
      form.bio = row.bio;
      form.avatar = '';

      if (avatarUploadRef.value) avatarUploadRef.value.clearFiles();
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
        // If name is empty, auto-fill with filename (without extension)
        if (!form.name) {
          const name = file.name.split('.').slice(0, -1).join('.') || file.name;
          form.name = name;
        }
      }
    };

    const handleElFileRemove = () => {
      avatarFile.value = null;
    };

    const handleExceed: UploadProps['onExceed'] = () => {
      ElMessage.warning(`限制选择 1 个文件，请先移除旧文件`);
    };

    const handlePaste = (event: ClipboardEvent) => {
      const items = event.clipboardData?.items;
      if (!items) return;

      for (let i = 0; i < items.length; i++) {
        const item = items[i];
        if (item && item.type.indexOf('image') !== -1) {
          const file = item.getAsFile();
          if (file) {
            avatarFile.value = file;
            if (avatarUploadRef.value) {
              avatarUploadRef.value.clearFiles();
              // @ts-ignore
              avatarUploadRef.value.handleStart(file);
            }
            ElMessage.success('已从剪贴板读取图片');
            break;
          }
        }
      }
    };

    const submitForm = async () => {
      submitting.value = true;
      uploadProgress.value = 0;
      try {
        let aPath = form.avatar;

        // 如果是新建艺术家且有头像文件，先创建艺术家再上传头像
        if (!isEdit.value && avatarFile.value) {
          // 先创建艺术家（不带头像）
          const createRes = await axios.post('/api/v1/admin/artists', {
            name: form.name,
            bio: form.bio,
            avatar: ''
          });
          const artistId = createRes.data.data.artist_id;

          // 上传头像到艺术家专属路径
          const res = await axios.post('/api/v1/uploads/presign', {
            upload_type: 'artist_avatar',
            target_id: artistId,
            filename: avatarFile.value.name,
            content_type: avatarFile.value.type
          });
          const { upload_url, object_key } = res.data.data;
          await axios.put(upload_url, avatarFile.value, {
            headers: { 'Content-Type': avatarFile.value.type },
            onUploadProgress: (progressEvent) => {
              if (progressEvent.total) {
                uploadProgress.value = Math.round((progressEvent.loaded * 100) / progressEvent.total);
              }
            }
          });

          // 更新艺术家头像
          await axios.put(`/api/v1/admin/artists/${artistId}`, {
            avatar: object_key
          });

          ElMessage.success('创建成功');
          dialogVisible.value = false;
          fetchData();
        } else {
          // 编辑模式或没有新头像
          if (avatarFile.value) {
            const res = await axios.post('/api/v1/uploads/presign', {
              upload_type: 'artist_avatar',
              target_id: form.id,
              filename: avatarFile.value.name,
              content_type: avatarFile.value.type
            });
            const { upload_url, object_key } = res.data.data;
            await axios.put(upload_url, avatarFile.value, {
              headers: { 'Content-Type': avatarFile.value.type },
              onUploadProgress: (progressEvent) => {
                if (progressEvent.total) {
                  uploadProgress.value = Math.round((progressEvent.loaded * 100) / progressEvent.total);
                }
              }
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
        }
      } catch (e) {
        ElMessage.error('操作失败');
      } finally {
        submitting.value = false;
        avatarFile.value = null;
        uploadProgress.value = 0;
      }
    };

    return {
      tableData,
      currentPage,
      pageSize,
      total,
      loading,
      avatarUploadRef,
      uploadProgress,
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
      handlePaste,
      submitForm,
      avatarFile
    };
  }
});
