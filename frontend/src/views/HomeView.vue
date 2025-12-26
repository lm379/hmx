<template>
  <div class="home-view">
    <!-- Video Grid -->
    <div class="video-grid" v-if="!loading">
      <VideoCard v-for="opera in operas" :key="opera.opera_id" :opera="opera"
        @click="navigateToVideo(opera.opera_id)" />
    </div>

    <div v-else class="loading-state">
      加载中...
    </div>

    <div v-if="!loading && operas.length === 0" class="empty-state">
      暂无推荐内容
    </div>

    <Pagination
      v-if="!loading"
      :current-page="currentPage"
      :total-items="totalItems"
      :page-size="pageSize"
      @page-change="handlePageChange"
    />
  </div>
</template>

<script lang="ts" src="../scripts/views/HomeView.ts"></script>

<style scoped src="../styles/views/HomeView.css"></style>
