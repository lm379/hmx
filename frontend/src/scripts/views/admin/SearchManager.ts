import { defineComponent, ref, onMounted } from 'vue';
import { ElMessage, ElMessageBox } from 'element-plus';
import { Refresh, VideoPlay, User } from '@element-plus/icons-vue';
import { getSearchStats, reindexOperas, reindexArtists } from '../../../api/search';

export default defineComponent({
  name: 'SearchManager',
  components: {
    Refresh,
    VideoPlay,
    User
  },
  setup() {
    const statsLoading = ref(false);
    const reindexingOperas = ref(false);
    const reindexingArtists = ref(false);
    const reindexingAll = ref(false);

    const operaCount = ref(0);
    const artistCount = ref(0);
    const operaIndexing = ref(false);
    const artistIndexing = ref(false);
    const operaFieldDistribution = ref<Record<string, number> | null>(null);
    const artistFieldDistribution = ref<Record<string, number> | null>(null);

    // 加载统计信息
    const loadStats = async () => {
      statsLoading.value = true;
      try {
        const response = await getSearchStats();
        // response.data is SearchStatsResponse.data which contains operas and artists objects
        
        if (response.data.operas) {
          operaCount.value = response.data.operas.numberOfDocuments;
          operaIndexing.value = response.data.operas.isIndexing;
          operaFieldDistribution.value = response.data.operas.fieldDistribution || null;
        }
        
        if (response.data.artists) {
          artistCount.value = response.data.artists.numberOfDocuments;
          artistIndexing.value = response.data.artists.isIndexing;
          artistFieldDistribution.value = response.data.artists.fieldDistribution || null;
        }
      } catch (error: any) {
        console.error('加载统计失败:', error);
        ElMessage.error(error.response?.data?.msg || '加载统计失败');
      } finally {
        statsLoading.value = false;
      }
    };

    // 重新索引作品
    const handleReindexOperas = async () => {
      try {
        await ElMessageBox.confirm(
          '确定要重新索引所有作品吗？这可能需要一些时间。',
          '确认操作',
          {
            confirmButtonText: '确定',
            cancelButtonText: '取消',
            type: 'warning',
          }
        );

        reindexingOperas.value = true;
        const response = await reindexOperas();
        ElMessage.success(response.data.message || '作品索引已开始');
        
        // 延迟刷新统计
        setTimeout(() => {
          loadStats();
        }, 2000);
      } catch (error: any) {
        if (error !== 'cancel') {
          console.error('重新索引失败:', error);
          ElMessage.error(error.response?.data?.msg || '重新索引失败');
        }
      } finally {
        reindexingOperas.value = false;
      }
    };

    // 重新索引艺术家
    const handleReindexArtists = async () => {
      try {
        await ElMessageBox.confirm(
          '确定要重新索引所有艺术家吗？这可能需要一些时间。',
          '确认操作',
          {
            confirmButtonText: '确定',
            cancelButtonText: '取消',
            type: 'warning',
          }
        );

        reindexingArtists.value = true;
        const response = await reindexArtists();
        ElMessage.success(response.data.message || '艺术家索引已开始');
        
        // 延迟刷新统计
        setTimeout(() => {
          loadStats();
        }, 2000);
      } catch (error: any) {
        if (error !== 'cancel') {
          console.error('重新索引失败:', error);
          ElMessage.error(error.response?.data?.msg || '重新索引失败');
        }
      } finally {
        reindexingArtists.value = false;
      }
    };

    // 重新索引全部
    const handleReindexAll = async () => {
      try {
        await ElMessageBox.confirm(
          '确定要重新索引所有数据吗？这将同时索引作品和艺术家，可能需要较长时间。',
          '确认操作',
          {
            confirmButtonText: '确定',
            cancelButtonText: '取消',
            type: 'warning',
          }
        );

        reindexingAll.value = true;
        
        // 并行执行两个索引任务
        await Promise.all([
          reindexOperas(),
          reindexArtists(),
        ]);
        
        ElMessage.success('全部数据索引已开始');
        
        // 延迟刷新统计
        setTimeout(() => {
          loadStats();
        }, 3000);
      } catch (error: any) {
        if (error !== 'cancel') {
          console.error('重新索引失败:', error);
          ElMessage.error(error.response?.data?.msg || '重新索引失败');
        }
      } finally {
        reindexingAll.value = false;
      }
    };

    onMounted(() => {
      loadStats();
    });

    return {
      statsLoading,
      reindexingOperas,
      reindexingArtists,
      reindexingAll,
      operaCount,
      artistCount,
      operaIndexing,
      artistIndexing,
      operaFieldDistribution,
      artistFieldDistribution,
      loadStats,
      handleReindexOperas,
      handleReindexArtists,
      handleReindexAll,
      Refresh,
      VideoPlay,
      User
    };
  }
});
