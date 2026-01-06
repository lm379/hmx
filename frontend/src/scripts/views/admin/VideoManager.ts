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
    const videoProgress = ref(0);
    const avatarProgress = ref(0);

    const dialogVisible = ref(false);
    const isEdit = ref(false);
    const submitting = ref(false);

    const form = reactive({
      id: 0,
      title: '',
      description: '',
      video_path: '',
      avatar_path: '',
      artist_ids: [] as (number | string)[],
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
        const res = await axios.get('/api/v1/artists', {
          params: {
            page: 1,
            page_size: 100,
            // If backend supports filtering by name, add query here. 
            // Currently backend returns list, frontend filtering might be needed if list is small,
            // or rely on el-select filtering.
          }
        });
        artistOptions.value = res.data.data.list;
      } catch (e) {
        console.error(e);
        ElMessage.error('获取艺术家列表失败');
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

      if (videoUploadRef.value) videoUploadRef.value.clearFiles();
      if (avatarUploadRef.value) avatarUploadRef.value.clearFiles();
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

      if (videoUploadRef.value) videoUploadRef.value.clearFiles();
      if (avatarUploadRef.value) avatarUploadRef.value.clearFiles();
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
        if (type === 'video') {
          // 验证视频文件类型
          if (!file.raw.type.startsWith('video/')) {
            ElMessage.error('只能上传视频文件！');
            if (videoUploadRef.value) videoUploadRef.value.clearFiles();
            return;
          }
          const maxSize = 4 * 1024 * 1024 * 1024; // 4GB
          if (file.raw.size > maxSize) {
            ElMessage.error('视频文件大小不能超过 4GB！');
            if (videoUploadRef.value) videoUploadRef.value.clearFiles();
            return;
          }
          videoFile.value = file.raw;
          // If title is empty, auto-fill with filename (without extension)
          if (!form.title) {
            const name = file.name.split('.').slice(0, -1).join('.') || file.name;
            form.title = name;
          }
        }
        if (type === 'avatar') {
          // 验证图片文件类型
          if (!file.raw.type.startsWith('image/')) {
            ElMessage.error('只能上传图片文件！');
            if (avatarUploadRef.value) avatarUploadRef.value.clearFiles();
            return;
          }
          const maxSize = 5 * 1024 * 1024; // 5MB
          if (file.raw.size > maxSize) {
            ElMessage.error('图片文件大小不能超过 5MB！');
            if (avatarUploadRef.value) avatarUploadRef.value.clearFiles();
            return;
          }
          avatarFile.value = file.raw;
        }
      }
    };

    const handleElFileRemove = (type: string) => {
      if (type === 'video') videoFile.value = null;
      if (type === 'avatar') avatarFile.value = null;
    };

    const handleExceed: UploadProps['onExceed'] = () => {
      ElMessage.warning(`限制选择 1 个文件，请先移除旧文件`);
    };

    const beforeVideoUpload: UploadProps['beforeUpload'] = (rawFile) => {
      if (!rawFile.type.startsWith('video/')) {
        ElMessage.error('只能上传视频文件！');
        return false;
      }
      // Optional: Add file size limit (e.g., 500MB)
      const maxSize = 4 * 1024 * 1024 * 1024; // 4GB
      if (rawFile.size > maxSize) {
        ElMessage.error('视频文件大小不能超过 4GB！');
        return false;
      }
      return true;
    };

    const beforeAvatarUpload: UploadProps['beforeUpload'] = (rawFile) => {
      if (!rawFile.type.startsWith('image/')) {
        ElMessage.error('只能上传图片文件！');
        return false;
      }
      const maxSize = 5 * 1024 * 1024; // 5MB
      if (rawFile.size > maxSize) {
        ElMessage.error('图片文件大小不能超过 5MB！');
        return false;
      }
      return true;
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

    const uploadFile = async (file: File, type: string, customFilename?: string, operaId?: number) => {
      const filename = customFilename || file.name;

      let uploadType: string;
      let targetId: number;

      if (type === 'video') {
        // 视频上传到 /tmp 不需要 opera_id
        uploadType = 'video_upload';
        targetId = 0; // video_upload 不需要 targetId，传0
      } else {
        // 封面上传需要 opera_id
        if (!operaId) {
          throw new Error('Opera ID is required for cover upload');
        }
        uploadType = 'opera_cover';
        targetId = operaId;
      }

      const res = await axios.post('/api/v1/uploads/presign', {
        upload_type: uploadType,
        target_id: targetId,
        filename: filename,
        content_type: file.type
      });
      const { upload_url, object_key } = res.data.data;

      await axios.put(upload_url, file, {
        headers: { 'Content-Type': file.type },
        onUploadProgress: (progressEvent) => {
          const total = progressEvent.total || file.size;
          if (total) {
            const percent = Math.round((progressEvent.loaded * 100) / total);
            if (type === 'video') {
              videoProgress.value = percent;
            } else {
              avatarProgress.value = percent;
            }
          }
        }
      });

      return object_key;
    };

    const submitForm = async () => {
      submitting.value = true;
      videoProgress.value = 0;
      avatarProgress.value = 0;
      try {
        let vPath = form.video_path;
        let aPath = form.avatar_path;

        if (!isEdit.value && !videoFile.value) {
          ElMessage.error("请选择视频文件");
          submitting.value = false;
          return;
        }

        // Separate IDs and new names
        const artistIDs: number[] = [];
        const newArtistNames: string[] = [];
        form.artist_ids.forEach(val => {
          if (typeof val === 'number') {
            artistIDs.push(val);
          } else {
            newArtistNames.push(val);
          }
        });

        if (!isEdit.value) {
          // 新建模式：先上传视频到tmp，然后创建Opera并上传封面

          // 上传视频到 tmp（不需要opera_id）
          if (videoFile.value) {
            const ext = videoFile.value.name.split('.').pop();
            const videoFilename = ext ? `${form.title}.${ext}` : form.title;
            vPath = await uploadFile(videoFile.value, 'video', videoFilename);
          }

          // 创建Opera（使用上传后的视频路径）
          const createRes = await axios.post('/api/v1/operas/', {
            title: form.title,
            description: form.description,
            video_path: vPath,
            avatar_path: '',
            artist_ids: artistIDs,
            new_artist_names: newArtistNames,
            is_hidden: form.is_hidden
          });

          const operaId = createRes.data.data.opera_id;

          // 上传封面（如果有）
          if (avatarFile.value) {
            aPath = await uploadFile(avatarFile.value, 'cover', undefined, operaId);
            // 更新Opera的封面路径
            await axios.put(`/api/v1/admin/operas/${operaId}`, {
              avatar_path: aPath
            });
          }

        } else {
          // 编辑模式：视频上传到tmp，封面需要opera_id
          if (videoFile.value) {
            const ext = videoFile.value.name.split('.').pop();
            const videoFilename = ext ? `${form.title}.${ext}` : form.title;
            vPath = await uploadFile(videoFile.value, 'video', videoFilename);
          }
          if (avatarFile.value) {
            aPath = await uploadFile(avatarFile.value, 'cover', undefined, form.id);
          }

          const data = {
            title: form.title,
            description: form.description,
            video_path: vPath,
            avatar_path: aPath,
            artist_ids: artistIDs,
            new_artist_names: newArtistNames,
            is_hidden: form.is_hidden
          };

          await axios.put(`/api/v1/admin/operas/${form.id}`, data);
        }

        ElMessage.success(isEdit.value ? '更新成功' : '创建成功');
        dialogVisible.value = false;
        fetchData();
        // Refresh artist list to include newly created ones
        fetchArtists();
      } catch (e) {
        console.error(e);
        ElMessage.error('操作失败');
      } finally {
        submitting.value = false;
        videoFile.value = null;
        avatarFile.value = null;
        videoProgress.value = 0;
        avatarProgress.value = 0;
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
      videoProgress,
      avatarProgress,
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
      beforeVideoUpload,
      beforeAvatarUpload,
      handlePaste,
      submitForm,
      videoFile,
      avatarFile
    };
  }
});
