<template>
  <div class="home-view">
    <section class="recommend-header">
      <div class="channel-tabs" role="tablist" aria-label="首页推荐频道">
        <button
          v-for="channel in channels"
          :key="channel.key"
          class="channel-tab"
          :class="{ active: activeChannel === channel.key }"
          type="button"
          role="tab"
          :aria-selected="activeChannel === channel.key"
          @click="handleChannelChange(channel.key)"
        >
          {{ channel.label }}
        </button>
      </div>
    </section>

    <!-- Video Grid -->
    <div class="video-grid" v-if="!loading">
      <VideoCard v-for="opera in operas" :key="opera.opera_id" :opera="opera" show-recommendation-reason
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
