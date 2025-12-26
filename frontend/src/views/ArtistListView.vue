<template>
  <div class="artist-list-view">
    <div v-if="loading" class="loading-state">
      加载中...
    </div>

    <div v-else class="artist-grid">
      <div class="artist-card" v-for="artist in artists" :key="artist.artist_id"
        @click="$router.push(`/artist/${artist.artist_id}`)">
        <div class="artist-avatar-wrapper">
          <img :src="artist.avatar || `https://ui-avatars.com/api/?name=${artist.name}&background=random&size=128`"
            alt="Artist Avatar" class="artist-avatar" />
        </div>
        <h3 class="artist-name">{{ artist.name }}</h3>
        <p class="artist-bio-preview">{{ artist.bio || '暂无简介' }}</p>
      </div>
    </div>

    <div v-if="!loading && artists.length === 0" class="empty-state">
      暂无艺术家信息
    </div>

    <!-- Pagination Controls -->
    <Pagination
      v-if="!loading && totalItems > 0"
      :current-page="currentPage"
      :total-items="totalItems"
      :page-size="pageSize"
      @page-change="handlePageChange"
    />
  </div>
</template>

<script lang="ts" src="../scripts/views/ArtistListView.ts"></script>

<style scoped src="../styles/views/ArtistListView.css"></style>