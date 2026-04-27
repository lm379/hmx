<template>
  <div class="news-list-view">
    <div class="news-page-header">
      <h1 class="news-page-title">黄梅戏资讯</h1>
      <p class="news-page-subtitle">关注演出动态、艺术活动与黄梅戏文化传播</p>
    </div>

    <div v-if="loading" class="loading-state">加载中...</div>
    <div v-else-if="newsList.length === 0" class="empty-state">暂无新闻资讯</div>

    <div v-else class="news-grid">
      <article
        v-for="item in newsList"
        :key="item.news_id"
        class="news-card"
        @click="goDetail(item.news_id)"
      >
        <div class="news-cover">
          <img v-if="item.cover" :src="item.cover" :alt="item.title" loading="lazy" />
          <div v-else class="news-cover-placeholder">黄梅戏资讯</div>
        </div>
        <div class="news-card-body">
          <div class="news-meta">
            <span>{{ item.source || '平台资讯' }}</span>
            <span>{{ formatNewsDate(item.published_at || item.created_at) }}</span>
          </div>
          <h2 class="news-title">{{ item.title }}</h2>
          <div class="news-summary" v-html="renderSummary(item.summary)"></div>
        </div>
      </article>
    </div>

    <Pagination
      v-if="!loading && totalItems > 0"
      :current-page="currentPage"
      :total-items="totalItems"
      :page-size="pageSize"
      @page-change="handlePageChange"
    />
  </div>
</template>

<script lang="ts" src="../scripts/views/NewsListView.ts"></script>
<style scoped src="../styles/views/NewsListView.css"></style>
