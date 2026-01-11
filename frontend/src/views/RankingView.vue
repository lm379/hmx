<template>
  <div class="ranking-view">
    <div v-if="loading" class="loading-state">
      加载中...
    </div>

    <ul v-else class="ranking-list">
      <li v-for="(item, index) in rankings" :key="item.opera_id" class="ranking-item" @click="navigateToVideo(item.opera_id)">
        <div class="rank-number" :class="{ 'top-3': (index + (currentPage - 1) * pageSize) < 3 }">
          {{ index + 1 + (currentPage - 1) * pageSize }}
        </div>
        
        <div class="cover-wrapper">
           <img :src="item.avatar" alt="" class="cover-img" loading="lazy">
        </div>

        <div class="video-info">
          <div class="title" :title="item.opera_title">{{ item.opera_title }}</div>
          <div class="meta-info">
             <div class="meta-item">
                <svg viewBox="0 0 24 24" width="16" height="16" fill="#9499a0" class="icon"><path d="M12 12c2.21 0 4-1.79 4-4s-1.79-4-4-4-4 1.79-4 4 1.79 4 4 4zm0 2c-2.67 0-8 1.34-8 4v2h16v-2c0-2.66-5.33-4-8-4z"></path></svg>
                <span v-if="item.artists && item.artists.length > 0">{{ item.artists.map(a => a.name).join(' / ') }}</span>
                <span v-else>未知艺术家</span>
             </div>
             <div class="meta-item">
               <svg viewBox="0 0 24 24" width="16" height="16" fill="#9499a0" class="icon"><path d="M8 5v14l11-7z"></path></svg>
               {{ item.play_count }}
             </div>
             <div class="meta-item">
               发布于 {{ formatDate(item.created_at) }}
             </div>
          </div>
        </div>
      </li>
    </ul>

    <div v-if="!loading && rankings.length === 0" class="empty-state">
      暂无排行数据
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

<script lang="ts" src="../scripts/views/RankingView.ts"></script>

<style scoped src="../styles/views/RankingView.css"></style>
