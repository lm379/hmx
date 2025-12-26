import { defineComponent, computed, type PropType } from 'vue';
import type { Opera } from '../../types';

export default defineComponent({
  name: 'VideoCard',
  props: {
    opera: {
      type: Object as PropType<Opera>,
      required: true
    }
  },
  setup(props) {
    // Fallback image if avatar is null
    const coverUrl = computed(() => {
      return props.opera.avatar || 'https://via.placeholder.com/320x180?text=No+Cover';
    });

    // Mock uploader name since API doesn't provide it yet
    const uploaderName = computed(() => {
      // 超过三人则只显示前三人名字
      const names = props.opera.artists.slice(0, 3).map(artist => artist.Name);
      return names.length === 0 ? "黄梅戏官方" : (names.length === 3 ? names.join(", ") + "等" : names.join(", "));
    });

    const formatDuration = (durationStr?: string) => {
      if (!durationStr) return '';
      // Assume duration comes as "HH:MM:SS" or similar from SQL time type.
      return durationStr;
    };

    const formatDate = (dateStr: string) => {
      if (!dateStr) return '';
      const date = new Date(dateStr);
      const now = new Date();
      const diff = now.getTime() - date.getTime();
      const days = Math.floor(diff / (1000 * 60 * 60 * 24));

      if (days === 0) return '今天';
      if (days === 1) return '昨天';
      if (days < 7) return `${days}天前`;

      return `${date.getMonth() + 1}-${date.getDate()}`;
    };

    return {
      props,
      coverUrl,
      uploaderName,
      formatDuration,
      formatDate
    };
  }
});
