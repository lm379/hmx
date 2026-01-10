import { defineComponent, ref, computed, onMounted, onUnmounted } from 'vue';
import axios from 'axios';
import { View, Refresh, DataAnalysis, Document, Loading } from '@element-plus/icons-vue';

export default defineComponent({
  name: 'TaskQueueView',
  components: {
    View,
    Refresh,
    DataAnalysis,
    Document,
    Loading
  },
  setup() {
    const queueStatus = ref<any>({
      embedding_queue_length: 0,
      summary_queue_length: 0,
      pending_embedding_tasks: [],
      pending_summary_tasks: [],
      processing_embedding_tasks: [],
      processing_summary_tasks: [],
      completed_embedding_tasks: [],
      completed_summary_tasks: []
    });
    const statusPollingTimer = ref<any>(null);

    // 计算统计数据
    const embeddingStats = computed(() => {
      const processing = queueStatus.value.processing_embedding_tasks?.length || 0;
      const pending = queueStatus.value.pending_embedding_tasks?.length || 0;
      const completed = queueStatus.value.completed_embedding_tasks?.length || 0;
      return {
        active: processing + pending,
        total: processing + pending + completed
      };
    });

    const summaryStats = computed(() => {
      const processing = queueStatus.value.processing_summary_tasks?.length || 0;
      const pending = queueStatus.value.pending_summary_tasks?.length || 0;
      const completed = queueStatus.value.completed_summary_tasks?.length || 0;
      return {
        active: processing + pending,
        total: processing + pending + completed
      };
    });

    // 格式化时间
    const formatTime = (time: string) => {
      if (!time) return '';
      const date = new Date(time);
      const year = date.getFullYear();
      const month = String(date.getMonth() + 1).padStart(2, '0');
      const day = String(date.getDate()).padStart(2, '0');
      const hours = String(date.getHours()).padStart(2, '0');
      const minutes = String(date.getMinutes()).padStart(2, '0');
      const seconds = String(date.getSeconds()).padStart(2, '0');
      return `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`;
    };

    // 获取队列状态
    const fetchQueueStatus = async () => {
      try {
        const res = await axios.get('/api/v1/admin/tasks/queue/status');
        if (res.data && res.data.data) {
          const data = res.data.data;
          queueStatus.value = {
            ...data,
            embedding_queue_length: data.embedding_queue_length || 0,
            summary_queue_length: data.summary_queue_length || 0,
            pending_embedding_tasks: data.pending_embedding_tasks || [],
            pending_summary_tasks: data.pending_summary_tasks || [],
            processing_embedding_tasks: data.processing_embedding_tasks || [],
            processing_summary_tasks: data.processing_summary_tasks || [],
            completed_embedding_tasks: data.completed_embedding_tasks || [],
            completed_summary_tasks: data.completed_summary_tasks || []
          };
        }
      } catch (e) {
        console.error('获取队列状态失败', e);
      }
    };

    // 轮询逻辑
    const pollingLoop = async () => {
      await fetchQueueStatus();
      if (statusPollingTimer.value !== null) {
        statusPollingTimer.value = setTimeout(pollingLoop, 5000);
      }
    };

    // 开始轮询队列状态
    const startQueueStatusPolling = () => {
      if (statusPollingTimer.value) return;
      statusPollingTimer.value = setTimeout(pollingLoop, 0);
    };

    // 停止轮询
    const stopQueueStatusPolling = () => {
      if (statusPollingTimer.value) {
        clearTimeout(statusPollingTimer.value);
        statusPollingTimer.value = null;
      }
    };

    // 组件挂载时开始轮询
    onMounted(() => {
      startQueueStatusPolling();
    });

    // 组件卸载时清理定时器
    onUnmounted(() => {
      stopQueueStatusPolling();
    });

    return {
      queueStatus,
      embeddingStats,
      summaryStats,
      fetchQueueStatus,
      formatTime
    };
  }
});
