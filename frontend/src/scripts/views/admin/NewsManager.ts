import { defineComponent, onMounted, reactive, ref } from 'vue';
import axios from 'axios';
import {
  ElMessage,
  ElMessageBox,
  type FormInstance,
  type FormRules,
  type UploadFile,
  type UploadInstance,
  type UploadProps,
  type UploadRawFile
} from 'element-plus';
import {
  createNews,
  deleteNews,
  getAdminNews,
  getAdminNewsById,
  updateNews,
  type SaveNewsPayload
} from '../../../api/news';
import type { NewsListItem } from '../../../types';
import { formatDateTime } from '../../../utils/dateUtils';

export default defineComponent({
  name: 'NewsManager',
  setup() {
    const tableData = ref<NewsListItem[]>([]);
    const loading = ref(false);
    const currentPage = ref(1);
    const pageSize = ref(10);
    const total = ref(0);

    const dialogVisible = ref(false);
    const isEdit = ref(false);
    const submitting = ref(false);
    const formRef = ref<FormInstance>();
    const coverUploadRef = ref<UploadInstance>();
    const coverFile = ref<File | null>(null);

    const form = reactive<SaveNewsPayload & { news_id: number }>({
      news_id: 0,
      title: '',
      summary: '',
      content: '',
      cover: '',
      source: '',
      author: '',
      is_published: false,
      published_at: null
    });

    const rules: FormRules = {
      title: [
        { required: true, message: '请输入新闻标题', trigger: 'blur' },
        { min: 2, max: 200, message: '标题长度在 2 到 200 个字符', trigger: 'blur' }
      ],
      summary: [
        { max: 500, message: '摘要不能超过 500 个字符', trigger: 'blur' }
      ],
      content: [
        { required: true, message: '请输入新闻正文', trigger: 'blur' },
        { min: 10, message: '正文至少需要 10 个字符', trigger: 'blur' }
      ]
    };

    const resetForm = () => {
      form.news_id = 0;
      form.title = '';
      form.summary = '';
      form.content = '';
      form.cover = '';
      form.source = '';
      form.author = '';
      form.is_published = false;
      form.published_at = null;
      coverFile.value = null;
      coverUploadRef.value?.clearFiles();
      formRef.value?.clearValidate();
    };

    const fetchData = async () => {
      loading.value = true;
      try {
        const res = await getAdminNews({
          page: currentPage.value,
          page_size: pageSize.value
        });
        tableData.value = res.data.list || [];
        total.value = res.data.pagination.total || 0;
      } catch {
        ElMessage.error('获取新闻列表失败');
      } finally {
        loading.value = false;
      }
    };

    const handleCreate = () => {
      isEdit.value = false;
      resetForm();
      dialogVisible.value = true;
    };

    const handleEdit = async (row: NewsListItem) => {
      isEdit.value = true;
      resetForm();
      try {
        // 编辑时拉详情，避免列表摘要字段不足导致正文丢失。
        const res = await getAdminNewsById(row.news_id);
        const detail = res.data;
        form.news_id = detail.news_id;
        form.title = detail.title;
        form.summary = detail.summary || '';
        form.content = detail.content || '';
        form.cover = detail.cover || '';
        form.source = detail.source || '';
        form.author = detail.author || '';
        form.is_published = detail.is_published;
        form.published_at = detail.published_at || null;
        coverFile.value = null;
        coverUploadRef.value?.clearFiles();
        dialogVisible.value = true;
      } catch {
        ElMessage.error('获取新闻详情失败');
      }
    };

    const buildPayload = (): SaveNewsPayload => ({
      title: form.title.trim(),
      summary: form.summary.trim(),
      content: form.content.trim(),
      cover: form.cover.trim(),
      source: form.source.trim(),
      author: form.author.trim(),
      is_published: form.is_published,
      published_at: form.published_at || null
    });

    const handleCoverChange = (file: UploadFile) => {
      if (!file.raw) return;
      if (!file.raw.type.startsWith('image/')) {
        ElMessage.error('只能上传图片文件');
        coverUploadRef.value?.clearFiles();
        coverFile.value = null;
        return;
      }
      coverFile.value = file.raw;
    };

    const handleCoverRemove = () => {
      coverFile.value = null;
    };

    const handleCoverExceed: UploadProps['onExceed'] = () => {
      ElMessage.warning('只能选择一张封面，请先移除当前图片');
    };

    const handlePaste = (event: ClipboardEvent) => {
      const items = event.clipboardData?.items;
      if (!items) return;

      for (let i = 0; i < items.length; i++) {
        const item = items[i];
        if (!item || !item.type.startsWith('image/')) continue;

        const file = item.getAsFile();
        if (!file) return;

        // 统一塞进 el-upload 的文件列表，后续保存仍走 uploadCoverIfNeeded。
        const uploadFile = file as UploadRawFile;
        uploadFile.uid = Date.now();
        coverFile.value = uploadFile;
        coverUploadRef.value?.clearFiles();
        coverUploadRef.value?.handleStart(uploadFile);
        ElMessage.success('已从剪贴板读取封面图片');
        break;
      }
    };

    const uploadCoverIfNeeded = async () => {
      if (!coverFile.value) return form.cover;

      // 新闻封面先走统一预签名接口，保存新闻时只提交对象存储 key。
      const presignRes = await axios.post('/api/v1/uploads/presign', {
        upload_type: 'news_cover',
        target_id: 0,
        filename: coverFile.value.name,
        content_type: coverFile.value.type
      });
      const { upload_url, object_key } = presignRes.data.data;
      await axios.put(upload_url, coverFile.value, {
        headers: { 'Content-Type': coverFile.value.type }
      });
      return object_key;
    };

    const handleSubmit = async () => {
      if (!formRef.value) return;
      try {
        await formRef.value.validate();
      } catch {
        return;
      }

      submitting.value = true;
      try {
        const uploadedCover = await uploadCoverIfNeeded();
        form.cover = uploadedCover;
        const payload = buildPayload();
        if (isEdit.value) {
          await updateNews(form.news_id, payload);
          ElMessage.success('新闻已更新');
        } else {
          await createNews(payload);
          ElMessage.success('新闻已创建');
          currentPage.value = 1;
        }
        dialogVisible.value = false;
        coverFile.value = null;
        coverUploadRef.value?.clearFiles();
        await fetchData();
      } catch (err: any) {
        ElMessage.error(err.response?.data?.msg || '保存失败');
      } finally {
        submitting.value = false;
      }
    };

    const handleDelete = (row: NewsListItem) => {
      ElMessageBox.confirm(`确定删除新闻「${row.title}」吗？`, '提示', {
        type: 'warning'
      }).then(async () => {
        try {
          await deleteNews(row.news_id);
          ElMessage.success('删除成功');
          await fetchData();
        } catch {
          ElMessage.error('删除失败');
        }
      });
    };

    const formatTime = (value: string) => formatDateTime(value);

    onMounted(() => {
      fetchData();
    });

    return {
      tableData,
      loading,
      currentPage,
      pageSize,
      total,
      dialogVisible,
      isEdit,
      submitting,
      formRef,
      coverUploadRef,
      coverFile,
      form,
      rules,
      fetchData,
      handleCreate,
      handleEdit,
      handleSubmit,
      handleDelete,
      handleCoverChange,
      handleCoverRemove,
      handleCoverExceed,
      handlePaste,
      formatTime
    };
  }
});
