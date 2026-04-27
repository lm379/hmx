<template>
  <article class="news-detail-view">
    <div v-if="loading" class="loading-state">加载中...</div>

    <template v-else-if="news">
      <button class="back-btn" type="button" @click="$router.push('/news')">返回资讯列表</button>

      <header class="detail-header">
        <h1 class="detail-title">{{ news.title }}</h1>
        <div class="detail-meta">
          <span>{{ news.source || '平台资讯' }}</span>
          <span v-if="news.author">{{ news.author }}</span>
          <span>{{ formatNewsDate(news.published_at || news.created_at) }}</span>
        </div>
        <div v-if="news.summary" class="detail-summary" v-html="summaryHtml"></div>
      </header>

      <img v-if="news.cover" class="detail-cover" :src="news.cover" :alt="news.title" />

      <section class="detail-content" v-html="contentHtml"></section>
    </template>

    <div v-else class="empty-state">新闻不存在或尚未发布</div>
  </article>
</template>

<script lang="ts" src="../scripts/views/NewsDetailView.ts"></script>
<style scoped src="../styles/views/NewsDetailView.css"></style>
