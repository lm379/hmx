import { defineComponent, computed, type PropType } from 'vue';
import type { OperaListItem, SimpleOpera } from '../../types';
type VideoCardProps = OperaListItem | SimpleOpera;

export default defineComponent({
  name: 'VideoCard',
  props: {
    opera: {
      type: Object as PropType<VideoCardProps>,
      required: true
    }
  },
  setup(props) {
    // Fallback image if avatar is null
    const coverUrl = computed(() => {
      return props.opera.avatar || 'https://via.placeholder.com/320x180?text=No+Cover';
    });

    // 显示艺术家名字
    const uploaderName = computed(() => {
      // SimpleOpera 不包含 artists 字段，或者检查 artists 是否存在且有内容
      if ('artists' in props.opera && props.opera.artists && props.opera.artists.length > 0) {
        const names = props.opera.artists.slice(0, 3).map(artist => artist.name);
        return names.length === 3 ? names.join(", ") + "等" : names.join(", ");
      }
      return null; // 返回 null 表示不显示艺术家信息
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
