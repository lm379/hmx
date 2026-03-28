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
import RelatedSourcesScript from '../scripts/components/RelatedSources';
export default RelatedSourcesScript;
</script>

<style scoped src="../styles/components/RelatedSources.css"></style>
