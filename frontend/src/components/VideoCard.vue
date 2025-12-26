<template>
  <div class="video-card">
    <div class="card-cover">
      <img :src="coverUrl" alt="Cover" loading="lazy" />
      <div class="duration" v-if="opera.duration">{{ formatDuration(opera.duration) }}</div>
    </div>
    <div class="card-info">
      <h3 class="title" :title="opera.opera_title">{{ opera.opera_title }}</h3>
      <div class="uploader-info">
        <span class="uploader-name">
           <svg viewBox="0 0 24 24" width="16" height="16" fill="#9499a0" style="vertical-align: text-bottom; margin-right: 4px;"><path d="M12 12c2.21 0 4-1.79 4-4s-1.79-4-4-4-4 1.79-4 4 1.79 4 4 4zm0 2c-2.67 0-8 1.34-8 4v2h16v-2c0-2.66-5.33-4-8-4z"></path></svg>
           {{ uploaderName }}
        </span>
        <span class="date">{{ formatDate(opera.created_at) }}</span>
      </div>
    </div>
  </div>
</template>

<script lang="ts" src="../scripts/components/VideoCard.ts"></script>

<style scoped src="../styles/components/VideoCard.css"></style>
