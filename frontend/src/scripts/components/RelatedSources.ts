import { defineComponent, ref } from 'vue';
import type { Source } from '../../api/qa';

export default defineComponent({
  name: 'RelatedSources',
  props: {
    sources: {
      type: Array as () => Source[],
      default: () => []
    }
  },
  setup() {
    const expanded = ref(false);

    function truncate(text: string, max: number): string {
      if (!text) return '';
      return text.length > max ? text.slice(0, max) + '…' : text;
    }

    return { expanded, truncate };
  }
});
