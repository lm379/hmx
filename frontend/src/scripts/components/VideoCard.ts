import { defineComponent, computed, type PropType } from 'vue';
import type { OperaListItem, SimpleOpera } from '../../types';
import { formatDate, formatDuration } from '../../utils/dateUtils';
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
      return "黄梅戏官方"; // 默认显示
    });

    return {
      props,
      coverUrl,
      uploaderName,
      formatDuration,
      formatDate
    };
  }
});
