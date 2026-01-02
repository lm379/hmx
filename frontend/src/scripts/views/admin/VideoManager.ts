import { defineComponent, ref, onMounted, reactive } from 'vue';
import axios from 'axios';
import { ElMessage, ElMessageBox, type UploadFile, type UploadInstance, type UploadProps } from 'element-plus';

export default defineComponent({
  name: 'VideoManager',
  setup() {
    const tableData = ref([]);
    const currentPage = ref(1);
    const pageSize = ref(10);
    const total = ref(0);
    const loading = ref(false);
    const videoUploadRef = ref<UploadInstance>();
    const avatarUploadRef = ref<UploadInstance>();

    const dialogVisible = ref(false);
    const isEdit = ref(false);
    const submitting = ref(false);

    const form = reactive({
      id: 0,
      title: '',
      description: '',
      video_path: '',
      avatar_path: '',
      artist_ids: [] as number[],
      is_hidden: false,
    });

    const artistOptions = ref<any[]>([]);

    const fetchData = async () => {
      loading.value = true;
      try {
        const res = await axios.get('/api/v1/admin/operas', {
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

    const fetchArtists = async () => {
      try {
        const res = await axios.get('/api/v1/artists');
        artistOptions.value = res.data.data.list;
      } catch (e) {
        console.error(e);
      }
    };

    onMounted(() => {
      fetchData();
      fetchArtists();
    });

    const handleCreate = () => {
      isEdit.value = false;
      form.id = 0;
      form.title = '';
      form.description = '';
      form.video_path = '';
      form.avatar_path = '';
      form.artist_ids = [];
      form.is_hidden = false;
      
      if(videoUploadRef.value) videoUploadRef.value.clearFiles();
      if(avatarUploadRef.value) avatarUploadRef.value.clearFiles();
      videoFile.value = null;
      avatarFile.value = null;

      dialogVisible.value = true;
    };

    const handleEdit = (row: any) => {
      isEdit.value = true;
      form.id = row.opera_id;
      form.title = row.opera_title;
      form.description = row.description;
      form.video_path = '';
      form.avatar_path = '';
      form.artist_ids = row.artists.map((a: any) => a.artist_id);
      form.is_hidden = row.is_hidden;
      
      if(videoUploadRef.value) videoUploadRef.value.clearFiles();
      if(avatarUploadRef.value) avatarUploadRef.value.clearFiles();
      videoFile.value = null;
      avatarFile.value = null;

      dialogVisible.value = true;
    };

    const handleToggleHidden = async (row: any) => {
      try {
        await axios.put(`/api/v1/admin/operas/${row.opera_id}`, {
          is_hidden: !row.is_hidden
        });
        ElMessage.success('操作成功');
        fetchData();
      } catch (e) {
        ElMessage.error('操作失败');
      }
    };

    const handleDelete = (row: any) => {
      ElMessageBox.confirm('确定删除该视频吗？', '提示', {
        type: 'warning'
      }).then(async () => {
        try {
          await axios.delete(`/api/v1/admin/operas/${row.opera_id}`);
          ElMessage.success('删除成功');
          fetchData();
        } catch (e) {
          ElMessage.error('删除失败');
        }
      });
    };

    const videoFile = ref<File | null>(null);
    const avatarFile = ref<File | null>(null);

    const handleElFileChange = (file: UploadFile, type: string) => {
      if (file.raw) {
        if (type === 'video') videoFile.value = file.raw;
        if (type === 'avatar') avatarFile.value = file.raw;
      }
    };

    const handleElFileRemove = (type: string) => {
      if (type === 'video') videoFile.value = null;
      if (type === 'avatar') avatarFile.value = null;
    };

    const handleExceed: UploadProps['onExceed'] = () => {
      ElMessage.warning(`限制选择 1 个文件，请先移除旧文件`);
    };

    const uploadFile = async (file: File, type: string) => {
      const res = await axios.post('/api/v1/uploads/presign', {
        filename: file.name,
        content_type: file.type,
        upload_type: type === 'video' ? 'videos' : 'avatars'
      });
      const { upload_url, object_key } = res.data.data;
      
      await axios.put(upload_url, file, {
        headers: { 'Content-Type': file.type }
      });
      
      return object_key;
    };

    const submitForm = async () => {
      submitting.value = true;
      try {
        let vPath = form.video_path;
        let aPath = form.avatar_path;
        
        if (!isEdit.value && !videoFile.value) {
          ElMessage.error("请选择视频文件");
          submitting.value = false;
          return;
        }

        if (videoFile.value) {
          vPath = await uploadFile(videoFile.value, 'video');
        }
        if (avatarFile.value) {
          aPath = await uploadFile(avatarFile.value, 'avatar');
        }

        const data = {
          title: form.title,
          description: form.description,
          video_path: vPath,
          avatar_path: aPath,
          artist_ids: form.artist_ids,
          is_hidden: form.is_hidden
        };

        if (isEdit.value) {
          await axios.put(`/api/v1/admin/operas/${form.id}`, data);
        } else {
          await axios.post('/api/v1/operas/', data);
        }

        ElMessage.success(isEdit.value ? '更新成功' : '创建成功');
        dialogVisible.value = false;
        fetchData();
      } catch (e) {
        console.error(e);
        ElMessage.error('操作失败');
      } finally {
        submitting.value = false;
        videoFile.value = null;
        avatarFile.value = null;
      }
    };

    return {
      tableData,
      currentPage,
      pageSize,
      total,
      loading,
      videoUploadRef,
      avatarUploadRef,
      dialogVisible,
      isEdit,
      submitting,
      form,
      artistOptions,
      fetchData,
      handleCreate,
      handleEdit,
      handleToggleHidden,
      handleDelete,
      handleElFileChange,
      handleElFileRemove,
      handleExceed,
      submitForm,
      videoFile,
      avatarFile
    };
  }
});
