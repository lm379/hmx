import { computed, defineComponent, ref, onMounted, type Ref } from 'vue';
import { useRouter } from 'vue-router';
import VideoCard from '../../components/VideoCard.vue';
import Pagination from '../../components/Pagination.vue';
import axios from 'axios';
import type { OperaListItem } from '../../types';

type RecommendationChannel = 'for_you' | 'hot' | 'latest';

interface ChannelOption {
  key: RecommendationChannel;
  label: string;
  description: string;
}

export default defineComponent({
  name: 'HomeView',
  components: {
    VideoCard,
    Pagination
  },
  setup() {
    const router = useRouter();
    const operas: Ref<OperaListItem[]> = ref([]);
    const loading = ref(true);
    const currentPage = ref(1);
    const pageSize = ref(15); // Grid usually fits 12 better (4x3)
    const totalItems = ref(0);
    const activeChannel = ref<RecommendationChannel>('for_you');
    const defaultChannel: ChannelOption = {
      key: 'for_you',
      label: '为你推荐',
      description: '结合观看、点赞和收藏记录，为你挑选可能感兴趣的黄梅戏内容'
    };
    const channels: ChannelOption[] = [
      defaultChannel,
      {
        key: 'hot',
        label: '热门',
        description: '按播放、点赞、收藏等互动热度排序，发现大家正在看的作品'
      },
      {
        key: 'latest',
        label: '最新',
        description: '查看平台最新收录的黄梅戏作品'
      }
    ];

    const activeChannelInfo = computed(() => {
      return channels.find(channel => channel.key === activeChannel.value) || defaultChannel;
    });

    const fetchOperas = async (page = 1) => {
      try {
        loading.value = true;
        // 使用推荐接口（支持游客和登录用户）
        const response = await axios.get('/api/v1/recommendations/', {
          params: {
            page: page,
            page_size: pageSize.value,
            channel: activeChannel.value,
            similarity_weight: 0.5,
            interest_weight: 0.5
          }
        });

        if (response.data && response.data.data) {
          operas.value = response.data.data.operas || [];

          // 处理分页信息
          if (response.data.data.pagination) {
            totalItems.value = response.data.data.pagination.total || 0;
            currentPage.value = response.data.data.pagination.page || page;
          } else {
            totalItems.value = operas.value.length;
            currentPage.value = page;
          }
        } else {
          console.error("Unexpected API response structure:", response.data);
        }
      } catch (error) {
        console.error("Failed to fetch recommendations, falling back to normal list:", error);
        // 降级到普通列表
        try {
          const fallbackResponse = await axios.get('/api/v1/operas/', {
            params: {
              page: page,
              page_size: pageSize.value
            }
          });
          if (fallbackResponse.data && fallbackResponse.data.data) {
            operas.value = fallbackResponse.data.data.list || [];
            if (fallbackResponse.data.data.pagination) {
              totalItems.value = fallbackResponse.data.data.pagination.total;
              currentPage.value = fallbackResponse.data.data.pagination.page;
            }
          }
        } catch (fallbackError) {
          console.error("Fallback also failed:", fallbackError);
        }
      } finally {
        loading.value = false;
      }
    };

    const navigateToVideo = (id: number) => {
      router.push(`/video/${id}`);
    };

    const handlePageChange = (newPage: number) => {
      fetchOperas(newPage);
      window.scrollTo({ top: 0, behavior: 'smooth' });
    };

    const handleChannelChange = (channel: RecommendationChannel) => {
      if (activeChannel.value === channel) return;
      activeChannel.value = channel;
      fetchOperas(1);
      window.scrollTo({ top: 0, behavior: 'smooth' });
    };

    onMounted(() => {
      fetchOperas();
    });

    return {
      operas,
      loading,
      currentPage,
      pageSize,
      totalItems,
      activeChannel,
      activeChannelInfo,
      channels,
      fetchOperas,
      navigateToVideo,
      handlePageChange,
      handleChannelChange
    };
  }
});
