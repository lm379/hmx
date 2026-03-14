import { defineComponent, ref, computed, onMounted, onUnmounted } from 'vue';
import axios from 'axios';
import { ElMessage } from 'element-plus';
import type { FormInstance } from 'element-plus';
import { Refresh } from '@element-plus/icons-vue';
import { marked } from 'marked';
import DOMPurify from 'dompurify';

interface KnowledgeDocument {
  doc_id: string;
  title: string;
  source_type: 'opera' | 'professional';
  chunks_count: number;
  embedding_status: 'pending' | 'processing' | 'completed' | 'failed';
  is_active: boolean;
  created_at: string;
  content?: string;
}

interface ChunkItem {
  chunk_id: string;
  chunk_index: number;
  chunk_text: string;
}

interface DocumentDetail extends KnowledgeDocument {
  updated_at: string;
  chunks: ChunkItem[];
}

interface KnowledgeStats {
  total_documents: number;
  professional_docs: number;
  opera_docs: number;
  completed_docs: number;
  pending_docs: number;
  failed_docs: number;
  total_chunks: number;
  avg_chunks_per_doc: number;
}

// 配置 marked：安全、换行友好
marked.setOptions({ breaks: true });

function renderMarkdown(text: string): string {
  const raw = marked.parse(text) as string;
  return DOMPurify.sanitize(raw);
}

const API = '/api/v1/admin/knowledge';

export default defineComponent({
  name: 'KnowledgeManager',

  setup() {
    const tableData = ref<KnowledgeDocument[]>([]);
    const loading = ref(false);
    const uploading = ref(false);
    const importing = ref(false);
    const currentPage = ref(1);
    const pageSize = ref(20);
    const total = ref(0);
    const filterSource = ref('');
    const filterStatus = ref('');

    const uploadDialogVisible = ref(false);
    const detailDialogVisible = ref(false);

    // 轮询相关
    const pollingTimer = ref<ReturnType<typeof setTimeout> | null>(null);
    const POLL_INTERVAL = 5000;

    // 详情弹窗 tab
    const detailActiveTab = ref<'full' | 'chunks'>('full');

    const uploadFormRef = ref<FormInstance>();
    const uploadForm = ref({
      title: '',
      source_type: 'professional' as 'professional' | 'opera',
      content: '',
    });

    const uploadFormRules = {
      title: [
        { required: true, message: '请输入文档标题', trigger: 'blur' },
        { min: 2, max: 255, message: '标题长度在 2 到 255 个字符', trigger: 'blur' },
      ],
      source_type: [{ required: true, message: '请选择文档来源', trigger: 'change' }],
      content: [
        { required: true, message: '请输入文档内容', trigger: 'blur' },
        { min: 10, message: '内容至少需要 10 个字符', trigger: 'blur' },
      ],
    };

    const stats = ref<Partial<KnowledgeStats>>({});
    const currentDocument = ref<DocumentDetail | null>(null);

    // 前端过滤
    const filteredData = computed(() => {
      return tableData.value.filter((doc) => {
        const matchSource = !filterSource.value || doc.source_type === filterSource.value;
        const matchStatus = !filterStatus.value || doc.embedding_status === filterStatus.value;
        return matchSource && matchStatus;
      });
    });

    // 完整内容的 markdown HTML
    const fullContentHtml = computed(() => {
      const content = currentDocument.value?.content ?? '';
      return renderMarkdown(content);
    });

    // 各分块的 markdown HTML（缓存避免重复计算）
    const chunksHtml = computed(() => {
      return (currentDocument.value?.chunks ?? []).map((c) => renderMarkdown(c.chunk_text));
    });

    const fetchStats = async () => {
      try {
        const res = await axios.get(`${API}/stats`);
        stats.value = res.data.data ?? {};
      } catch {
        // 静默失败
      }
    };

    const fetchDocuments = async () => {
      loading.value = true;
      try {
        const res = await axios.get(`${API}/documents`, {
          params: { page: currentPage.value, page_size: pageSize.value },
        });
        const body = res.data;
        tableData.value = body.data?.list ?? body.data ?? [];
        total.value = body.data?.pagination?.total ?? body.data?.length ?? 0;
      } catch {
        ElMessage.error('获取文档列表失败');
      } finally {
        loading.value = false;
      }
    };

    const hasActiveJobs = computed(() =>
      tableData.value.some(
        (d) => d.embedding_status === 'pending' || d.embedding_status === 'processing',
      ),
    );

    const stopPolling = () => {
      if (pollingTimer.value !== null) {
        clearTimeout(pollingTimer.value);
        pollingTimer.value = null;
      }
    };

    const startPolling = () => {
      stopPolling();
      const loop = async () => {
        await Promise.all([fetchDocuments(), fetchStats()]);
        if (hasActiveJobs.value) {
          pollingTimer.value = setTimeout(loop, POLL_INTERVAL);
        } else {
          pollingTimer.value = null;
        }
      };
      pollingTimer.value = setTimeout(loop, POLL_INTERVAL);
    };

    const handleRefresh = async () => {
      await Promise.all([fetchDocuments(), fetchStats()]);
      if (hasActiveJobs.value) {
        startPolling();
      }
    };

    const handleFilter = () => {};

    const handleShowUploadDialog = () => {
      uploadForm.value = { title: '', source_type: 'professional', content: '' };
      uploadDialogVisible.value = true;
    };

    const handleSubmitUpload = async () => {
      if (!uploadFormRef.value) return;
      try {
        await uploadFormRef.value.validate();
      } catch {
        return;
      }
      uploading.value = true;
      try {
        await axios.post(`${API}/documents`, uploadForm.value);
        ElMessage.success('文档上传成功，正在后台向量化处理...');
        uploadDialogVisible.value = false;
        currentPage.value = 1;
        await Promise.all([fetchDocuments(), fetchStats()]);
        startPolling();
      } catch (err: any) {
        ElMessage.error(err.response?.data?.msg || '上传失败');
      } finally {
        uploading.value = false;
      }
    };

    const handleViewDetail = async (doc: KnowledgeDocument) => {
      detailActiveTab.value = 'full';
      try {
        const res = await axios.get(`${API}/documents/${doc.doc_id}`);
        currentDocument.value = res.data.data ?? res.data;
        detailDialogVisible.value = true;
      } catch {
        ElMessage.error('获取文档详情失败');
      }
    };

    const handleDeleteDocument = async (docId: string) => {
      try {
        await axios.delete(`${API}/documents/${docId}`);
        ElMessage.success('文档已删除');
        await Promise.all([fetchDocuments(), fetchStats()]);
      } catch {
        ElMessage.error('删除失败');
      }
    };

    const handleImportOperas = async () => {
      importing.value = true;
      try {
        await axios.post(`${API}/import-operas`);
        ElMessage.success('作品字幕导入任务已启动，请稍后刷新查看进度');
        setTimeout(() => {
          fetchDocuments();
          fetchStats();
          startPolling();
        }, 2000);
      } catch (err: any) {
        ElMessage.error(err.response?.data?.msg || '导入失败');
      } finally {
        importing.value = false;
      }
    };

    const handlePageChange = (page: number) => {
      currentPage.value = page;
      fetchDocuments();
    };

    const handlePageSizeChange = (size: number) => {
      pageSize.value = size;
      currentPage.value = 1;
      fetchDocuments();
    };

    const getStatusType = (status: string): string => {
      const map: Record<string, string> = {
        pending: 'info',
        processing: 'warning',
        completed: 'success',
        failed: 'danger',
      };
      return map[status] ?? 'info';
    };

    const formatStatus = (status: string): string => {
      const map: Record<string, string> = {
        pending: '待处理',
        processing: '处理中',
        completed: '已完成',
        failed: '失败',
      };
      return map[status] ?? status;
    };

    const formatTime = (time: string): string => {
      if (!time) return '-';
      return new Date(time).toLocaleString('zh-CN', {
        year: 'numeric',
        month: '2-digit',
        day: '2-digit',
        hour: '2-digit',
        minute: '2-digit',
      });
    };

    onMounted(async () => {
      await Promise.all([fetchStats(), fetchDocuments()]);
      if (hasActiveJobs.value) {
        startPolling();
      }
    });

    onUnmounted(() => {
      stopPolling();
    });

    return {
      tableData,
      filteredData,
      loading,
      uploading,
      importing,
      currentPage,
      pageSize,
      total,
      filterSource,
      filterStatus,
      uploadDialogVisible,
      detailDialogVisible,
      detailActiveTab,
      uploadFormRef,
      uploadForm,
      uploadFormRules,
      stats,
      currentDocument,
      fullContentHtml,
      chunksHtml,
      handleFilter,
      handleRefresh,
      handleShowUploadDialog,
      Refresh,
      handleSubmitUpload,
      handleViewDetail,
      handleDeleteDocument,
      handleImportOperas,
      handlePageChange,
      handlePageSizeChange,
      getStatusType,
      formatStatus,
      formatTime,
    };
  },
});

