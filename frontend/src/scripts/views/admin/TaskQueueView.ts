import { defineComponent, ref, computed, onMounted, onUnmounted } from 'vue';
import axios from 'axios';
import { View, Refresh, DataAnalysis, Document, Loading, VideoPlay, ChatLineSquare, Picture } from '@element-plus/icons-vue';
import { formatDateTime } from '../../../utils/dateUtils';

export default defineComponent({
  name: 'TaskQueueView',
  components: {
    View,
    Refresh,
    DataAnalysis,
    Document,
    Loading,
    VideoPlay,
    ChatLineSquare,
    Picture
  },
  setup() {
    const queueStatus = ref<any>({
      embedding_queue_length: 0,
      summary_queue_length: 0,
      transcode_queue_length: 0,
      subtitle_queue_length: 0,
      cover_queue_length: 0,
      pending_embedding_tasks: [],
      pending_summary_tasks: [],
      pending_transcode_tasks: [],
      pending_subtitle_tasks: [],
      pending_cover_tasks: [],
      processing_embedding_tasks: [],
      processing_summary_tasks: [],
      processing_transcode_tasks: [],
      processing_subtitle_tasks: [],
      processing_cover_tasks: [],
      completed_embedding_tasks: [],
      completed_summary_tasks: [],
      completed_transcode_tasks: [],
      completed_subtitle_tasks: [],
      completed_cover_tasks: []
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

    const transcodeStats = computed(() => {
      const processing = queueStatus.value.processing_transcode_tasks?.length || 0;
      const pending = queueStatus.value.pending_transcode_tasks?.length || 0;
      const completed = queueStatus.value.completed_transcode_tasks?.length || 0;
      return {
        active: processing + pending,
        total: processing + pending + completed
      };
    });

    const subtitleStats = computed(() => {
      const processing = queueStatus.value.processing_subtitle_tasks?.length || 0;
      const pending = queueStatus.value.pending_subtitle_tasks?.length || 0;
      const completed = queueStatus.value.completed_subtitle_tasks?.length || 0;
      return {
        active: processing + pending,
        total: processing + pending + completed
      };
    });

    const coverStats = computed(() => {
      const processing = queueStatus.value.processing_cover_tasks?.length || 0;
      const pending = queueStatus.value.pending_cover_tasks?.length || 0;
      const completed = queueStatus.value.completed_cover_tasks?.length || 0;
      return {
        active: processing + pending,
        total: processing + pending + completed
      };
    });

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
            transcode_queue_length: data.transcode_queue_length || 0,
            subtitle_queue_length: data.subtitle_queue_length || 0,
            cover_queue_length: data.cover_queue_length || 0,
            pending_embedding_tasks: data.pending_embedding_tasks || [],
            pending_summary_tasks: data.pending_summary_tasks || [],
            pending_transcode_tasks: data.pending_transcode_tasks || [],
            pending_subtitle_tasks: data.pending_subtitle_tasks || [],
            pending_cover_tasks: data.pending_cover_tasks || [],
            processing_embedding_tasks: data.processing_embedding_tasks || [],
            processing_summary_tasks: data.processing_summary_tasks || [],
            processing_transcode_tasks: data.processing_transcode_tasks || [],
            processing_subtitle_tasks: data.processing_subtitle_tasks || [],
            processing_cover_tasks: data.processing_cover_tasks || [],
            completed_embedding_tasks: data.completed_embedding_tasks || [],
            completed_summary_tasks: data.completed_summary_tasks || [],
            completed_transcode_tasks: data.completed_transcode_tasks || [],
            completed_subtitle_tasks: data.completed_subtitle_tasks || [],
            completed_cover_tasks: data.completed_cover_tasks || []
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
      transcodeStats,
      subtitleStats,
      coverStats,
      fetchQueueStatus,
      formatTime: formatDateTime
    };
  }
});
