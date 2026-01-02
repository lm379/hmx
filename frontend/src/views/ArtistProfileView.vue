<template>
  <div class="artist-profile-view" v-if="artist">
    <!-- Header -->
    <div class="artist-header">
      <div class="artist-avatar-large">
        <img :src="artist.avatar || `https://ui-avatars.com/api/?name=${artist.name}&background=random&size=200`"
          alt="Avatar" />
      </div>
      <div class="artist-info">
        <h1 class="artist-name">{{ artist.name }}</h1>
        <p class="artist-bio">{{ artist.bio || '这位艺术家很低调，暂无简介。' }}</p>
      </div>
    </div>

    <!-- Works -->
    <div class="artist-works">
      <h2 class="section-title">代表作品</h2>
      <div class="video-grid" v-if="artist.operas && artist.operas.length > 0">
        <VideoCard v-for="opera in artist.operas" :key="opera.opera_id" :opera="opera"
          @click="$router.push(`/video/${opera.opera_id}`)" />
      </div>
      <div v-else class="empty-state">
        暂无收录作品
      </div>
    </div>
  </div>
  <div v-else-if="loading" class="loading-state">
    加载中...
  </div>
  <div v-else class="error-state">
    未找到艺术家信息
  </div>
</template>

<script lang="ts" src="../scripts/views/ArtistProfileView.ts"></script>

<style scoped src="../styles/views/ArtistProfileView.css"></style>
