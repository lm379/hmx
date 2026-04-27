<template>
  <div class="education-reader-view">
    <div v-if="loading" class="loading-state">加载中...</div>

    <template v-else-if="book">
      <div class="reader-header">
        <button class="back-btn" type="button" @click="$router.push('/education')">返回教育资源</button>
        <div class="reader-title-group">
          <h1 class="reader-title">{{ book.title }}</h1>
          <p v-if="book.description" class="reader-desc">{{ book.description }}</p>
        </div>
        <div class="reader-actions">
          <button class="zoom-btn" type="button" @click="zoomOut">缩小</button>
          <span class="zoom-label">{{ zoomPercent }}</span>
          <button class="zoom-btn" type="button" @click="zoomIn">放大</button>
          <a class="open-link" :href="book.pdf_url" target="_blank" rel="noopener">新窗口打开</a>
        </div>
      </div>

      <div class="pdf-status" v-if="pdfLoading">正在加载 PDF...</div>
      <div class="pdf-status error" v-else-if="pdfError">{{ pdfError }}</div>
      <div class="pdf-status" v-else-if="totalPages > 0">
        已渲染 {{ renderedCount }} / {{ totalPages }} 页
      </div>

      <div class="pdf-pages" v-if="totalPages > 0">
        <section
          v-for="page in pages"
          :key="page.pageNumber"
          :ref="(el) => setPageShellRef(el, page.pageNumber)"
          class="pdf-page-shell"
        >
          <div class="page-label">第 {{ page.pageNumber }} 页</div>
          <div class="page-canvas-wrap">
            <canvas :ref="(el) => setCanvasRef(el, page.pageNumber)" class="pdf-canvas"></canvas>
            <div v-if="page.status === 'pending'" class="page-placeholder">滚动到此处后加载</div>
            <div v-else-if="page.status === 'loading'" class="page-placeholder">正在渲染...</div>
            <div v-else-if="page.status === 'error'" class="page-placeholder error">本页加载失败</div>
          </div>
        </section>
      </div>
    </template>

    <div v-else class="empty-state">资源不存在或尚未发布</div>
  </div>
</template>

<script lang="ts" src="../scripts/views/EducationReaderView.ts"></script>
<style scoped src="../styles/views/EducationReaderView.css"></style>
