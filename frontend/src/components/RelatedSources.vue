<template>
  <div v-if="sources && sources.length" class="related-sources">
    <div class="sources-header" @click="expanded = !expanded">
      <span class="sources-title">
        <svg viewBox="0 0 24 24" width="14" height="14" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/>
          <polyline points="14 2 14 8 20 8"/>
        </svg>
        参考来源 ({{ sources.length }})
      </span>
      <svg
        class="chevron"
        :class="{ rotated: expanded }"
        viewBox="0 0 24 24" width="14" height="14" fill="none"
        stroke="currentColor" stroke-width="2"
      >
        <polyline points="6 9 12 15 18 9"/>
      </svg>
    </div>
    <div v-if="expanded" class="sources-list">
      <div
        v-for="(s, i) in sources"
        :key="s.chunk_id"
        class="source-item"
      >
        <span class="source-index">{{ i + 1 }}</span>
        <div class="source-body">
          <div class="source-doc-title">{{ s.doc_title }}</div>
          <div class="source-snippet">{{ truncate(s.chunk_text, 120) }}</div>
        </div>
        <span class="source-score">{{ (s.score * 100).toFixed(0) }}%</span>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent, ref } from 'vue';
import type { Source } from '../api/qa';

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
</script>

<style scoped>
.related-sources {
  margin-top: 8px;
  border: 1px solid rgba(255,255,255,0.12);
  border-radius: 8px;
  overflow: hidden;
  font-size: 12px;
}

.sources-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 6px 10px;
  cursor: pointer;
  background: rgba(255,255,255,0.06);
  user-select: none;
}

.sources-title {
  display: flex;
  align-items: center;
  gap: 5px;
  color: rgba(255,255,255,0.6);
}

.chevron {
  transition: transform 0.2s;
  color: rgba(255,255,255,0.4);
}
.chevron.rotated {
  transform: rotate(180deg);
}

.sources-list {
  padding: 6px 8px;
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.source-item {
  display: flex;
  align-items: flex-start;
  gap: 8px;
}

.source-index {
  flex-shrink: 0;
  width: 18px;
  height: 18px;
  border-radius: 50%;
  background: rgba(255,255,255,0.15);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 10px;
  color: rgba(255,255,255,0.7);
}

.source-body {
  flex: 1;
  overflow: hidden;
}

.source-doc-title {
  font-weight: 600;
  color: rgba(255,255,255,0.8);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.source-snippet {
  color: rgba(255,255,255,0.45);
  line-height: 1.4;
  margin-top: 2px;
}

.source-score {
  flex-shrink: 0;
  color: rgba(255,255,255,0.35);
  font-size: 11px;
}
</style>
