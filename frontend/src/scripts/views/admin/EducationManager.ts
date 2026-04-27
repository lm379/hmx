import { defineComponent, onMounted, reactive, ref } from 'vue';
import axios from 'axios';
import * as pdfjsLib from 'pdfjs-dist/build/pdf.mjs';
import pdfWorkerUrl from 'pdfjs-dist/build/pdf.worker.mjs?url';
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
  createEducationBook,
  deleteEducationBook,
  getAdminEducationBookById,
  getAdminEducationBooks,
  updateEducationBook,
  type SaveEducationBookPayload
} from '../../../api/education';
import type { EducationBookListItem } from '../../../types';
import { formatDateTime } from '../../../utils/dateUtils';

pdfjsLib.GlobalWorkerOptions.workerSrc = pdfWorkerUrl;

export default defineComponent({
  name: 'EducationManager',
  setup() {
    const tableData = ref<EducationBookListItem[]>([]);
    const loading = ref(false);
    const currentPage = ref(1);
    const pageSize = ref(10);
    const total = ref(0);

    const dialogVisible = ref(false);
    const isEdit = ref(false);
    const submitting = ref(false);
    const formRef = ref<FormInstance>();
    const pdfUploadRef = ref<UploadInstance>();
    const coverUploadRef = ref<UploadInstance>();
    const pdfFile = ref<File | null>(null);
    const coverFile = ref<File | null>(null);
    const coverPreviewUrl = ref('');
    const coverGenerating = ref(false);
    const coverMode = ref<'none' | 'auto' | 'custom'>('none');
    let coverGenerationTask: Promise<void> | null = null;

    const form = reactive<SaveEducationBookPayload & { book_id: number }>({
      book_id: 0,
      title: '',
      description: '',
      pdf_path: '',
      cover_path: '',
      is_published: false
    });

    const rules: FormRules = {
      title: [
        { required: true, message: '请输入书籍名称', trigger: 'blur' },
        { min: 2, max: 200, message: '书籍名称长度在 2 到 200 个字符', trigger: 'blur' }
      ],
      pdf_path: [
        { required: true, message: '请上传 PDF 文件', trigger: 'change' }
      ]
    };

    const resetForm = () => {
      form.book_id = 0;
      form.title = '';
      form.description = '';
      form.pdf_path = '';
      form.cover_path = '';
      form.is_published = false;
      pdfFile.value = null;
      coverFile.value = null;
      coverPreviewUrl.value = '';
      coverGenerating.value = false;
      coverMode.value = 'none';
      coverGenerationTask = null;
      pdfUploadRef.value?.clearFiles();
      coverUploadRef.value?.clearFiles();
      formRef.value?.clearValidate();
    };

    const fetchData = async () => {
      loading.value = true;
      try {
        const res = await getAdminEducationBooks({
          page: currentPage.value,
          page_size: pageSize.value
        });
        tableData.value = res.data.list || [];
        total.value = res.data.pagination.total || 0;
      } catch {
        ElMessage.error('获取教育资源失败');
      } finally {
        loading.value = false;
      }
    };

    const handleCreate = () => {
      isEdit.value = false;
      resetForm();
      dialogVisible.value = true;
    };

    const handleEdit = async (row: EducationBookListItem) => {
      isEdit.value = true;
      resetForm();
      try {
        const res = await getAdminEducationBookById(row.book_id);
        const detail = res.data;
        form.book_id = detail.book_id;
        form.title = detail.title;
        form.description = detail.description || '';
        // 后台保存对象存储 key，前台展示才使用完整 URL。
        form.pdf_path = detail.pdf_path;
        form.cover_path = detail.cover_path || '';
        coverPreviewUrl.value = detail.cover_url || '';
        coverMode.value = detail.cover_path ? 'custom' : 'none';
        form.is_published = detail.is_published;
        dialogVisible.value = true;
      } catch {
        ElMessage.error('获取教育资源详情失败');
      }
    };

    const handlePdfChange = (file: UploadFile) => {
      if (!file.raw) return;
      if (file.raw.type !== 'application/pdf') {
        ElMessage.error('只能上传 PDF 文件');
        pdfUploadRef.value?.clearFiles();
        pdfFile.value = null;
        form.pdf_path = '';
        return;
      }
      pdfFile.value = file.raw;
      if (!form.title.trim()) {
        form.title = file.name.replace(/\.pdf$/i, '');
      }
      // 先填一个占位值触发表单校验，真正 object key 在保存时上传后写入。
      form.pdf_path = file.name;
      // 未手动设置封面时，自动截取 PDF 第一页作为默认封面。
      if (coverMode.value !== 'custom') {
        coverGenerationTask = generateCoverFromPdf(file.raw);
      }
    };

    const handlePdfRemove = () => {
      pdfFile.value = null;
      form.pdf_path = '';
      coverGenerating.value = false;
      coverGenerationTask = null;
      if (coverMode.value === 'auto') {
        coverFile.value = null;
        coverPreviewUrl.value = '';
        form.cover_path = '';
        coverMode.value = 'none';
      }
    };

    const handlePdfExceed: UploadProps['onExceed'] = () => {
      ElMessage.warning('只能选择一个 PDF 文件，请先移除当前文件');
    };

    const uploadPdfIfNeeded = async () => {
      if (!pdfFile.value) return form.pdf_path;

      const presignRes = await axios.post('/api/v1/uploads/presign', {
        upload_type: 'education_pdf',
        target_id: 0,
        filename: pdfFile.value.name,
        content_type: pdfFile.value.type
      });
      const { upload_url, object_key } = presignRes.data.data;
      await axios.put(upload_url, pdfFile.value, {
        headers: { 'Content-Type': pdfFile.value.type }
      });
      return object_key;
    };

    const generateCoverFromPdf = async (file: File) => {
      coverGenerating.value = true;
      try {
        const data = new Uint8Array(await file.arrayBuffer());
        const pdf = await pdfjsLib.getDocument({ data }).promise;
        const page = await pdf.getPage(1);
        const viewport = page.getViewport({ scale: 1.2 });
        const canvas = document.createElement('canvas');
        const context = canvas.getContext('2d');
        if (!context) throw new Error('Canvas context unavailable');

        canvas.width = Math.floor(viewport.width);
        canvas.height = Math.floor(viewport.height);
        await page.render({ canvasContext: context, viewport }).promise;

        const blob = await new Promise<Blob | null>((resolve) => {
          canvas.toBlob(resolve, 'image/jpeg', 0.86);
        });
        if (!blob) throw new Error('Failed to create cover image');
        // 自动封面生成是兜底能力，不能覆盖管理员随后手动上传的封面。
        if (coverMode.value === 'custom') return;

        coverFile.value = new File([blob], `${file.name.replace(/\.pdf$/i, '')}-cover.jpg`, {
          type: 'image/jpeg'
        });
        if (coverPreviewUrl.value.startsWith('blob:')) {
          URL.revokeObjectURL(coverPreviewUrl.value);
        }
        coverPreviewUrl.value = URL.createObjectURL(blob);
        form.cover_path = coverFile.value.name;
        coverMode.value = 'auto';
      } catch (error) {
        console.error('Failed to generate PDF cover', error);
        if (coverMode.value === 'custom') return;
        coverFile.value = null;
        coverPreviewUrl.value = '';
        form.cover_path = '';
        coverMode.value = 'none';
        ElMessage.warning('PDF 封面生成失败，可保存 PDF 后稍后重新上传');
      } finally {
        coverGenerating.value = false;
      }
    };

    const setCustomCoverFile = (file: File) => {
      if (!file.type.startsWith('image/')) {
        ElMessage.error('封面只能上传图片文件');
        coverUploadRef.value?.clearFiles();
        return;
      }
      if (coverPreviewUrl.value.startsWith('blob:')) {
        URL.revokeObjectURL(coverPreviewUrl.value);
      }
      coverFile.value = file;
      coverPreviewUrl.value = URL.createObjectURL(file);
      form.cover_path = file.name;
      coverMode.value = 'custom';
    };

    const handleCoverChange = (file: UploadFile) => {
      if (!file.raw) return;
      setCustomCoverFile(file.raw);
    };

    const handleCoverRemove = () => {
      coverFile.value = null;
      form.cover_path = '';
      coverPreviewUrl.value = '';
      coverMode.value = 'none';
      if (pdfFile.value) {
        coverGenerationTask = generateCoverFromPdf(pdfFile.value);
      }
    };

    const handleCoverExceed: UploadProps['onExceed'] = () => {
      ElMessage.warning('只能选择一张封面图，请先移除当前图片');
    };

    const handlePaste = (event: ClipboardEvent) => {
      const items = event.clipboardData?.items;
      if (!items) return;

      for (let i = 0; i < items.length; i++) {
        const item = items[i];
        if (!item || !item.type.startsWith('image/')) continue;

        const file = item.getAsFile();
        if (!file) return;

        // 粘贴图片也复用 el-upload 文件列表，让界面、预览和保存流程保持一致。
        const uploadFile = file as UploadRawFile;
        uploadFile.uid = Date.now();
        coverUploadRef.value?.clearFiles();
        coverUploadRef.value?.handleStart(uploadFile);
        setCustomCoverFile(uploadFile);
        ElMessage.success('已从剪贴板读取 PDF 封面');
        break;
      }
    };

    const uploadCoverIfNeeded = async () => {
      if (!coverFile.value) return form.cover_path;

      const presignRes = await axios.post('/api/v1/uploads/presign', {
        upload_type: 'education_cover',
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

    const buildPayload = (): SaveEducationBookPayload => ({
      title: form.title.trim(),
      description: form.description.trim(),
      pdf_path: form.pdf_path,
      cover_path: form.cover_path,
      is_published: form.is_published
    });

    const handleSubmit = async () => {
      if (!formRef.value) return;
      try {
        await formRef.value.validate();
      } catch {
        return;
      }

      submitting.value = true;
      try {
        if (coverGenerationTask) {
          await coverGenerationTask;
        }
        form.pdf_path = await uploadPdfIfNeeded();
        form.cover_path = await uploadCoverIfNeeded();
        const payload = buildPayload();
        if (isEdit.value) {
          await updateEducationBook(form.book_id, payload);
          ElMessage.success('教育资源已更新');
        } else {
          await createEducationBook(payload);
          ElMessage.success('教育资源已创建');
          currentPage.value = 1;
        }
        dialogVisible.value = false;
        pdfFile.value = null;
        coverFile.value = null;
        pdfUploadRef.value?.clearFiles();
        coverUploadRef.value?.clearFiles();
        await fetchData();
      } catch (err: any) {
        ElMessage.error(err.response?.data?.msg || '保存失败');
      } finally {
        submitting.value = false;
      }
    };

    const handleDelete = (row: EducationBookListItem) => {
      ElMessageBox.confirm(`确定删除「${row.title}」吗？`, '提示', {
        type: 'warning'
      }).then(async () => {
        try {
          await deleteEducationBook(row.book_id);
          ElMessage.success('删除成功');
          await fetchData();
        } catch {
          ElMessage.error('删除失败');
        }
      });
    };

    const openPdf = (row: EducationBookListItem) => {
      window.open(row.pdf_url, '_blank', 'noopener');
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
      pdfUploadRef,
      coverUploadRef,
      pdfFile,
      coverPreviewUrl,
      coverGenerating,
      coverMode,
      form,
      rules,
      fetchData,
      handleCreate,
      handleEdit,
      handlePdfChange,
      handlePdfRemove,
      handlePdfExceed,
      handleCoverChange,
      handleCoverRemove,
      handleCoverExceed,
      handlePaste,
      handleSubmit,
      handleDelete,
      openPdf,
      formatTime
    };
  }
});
