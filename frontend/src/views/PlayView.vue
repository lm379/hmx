<template>
  <div class="play-view" v-if="!loading && opera">
    <div class="play-container">
      <!-- Left Column -->
      <div class="main-column">
        <!-- Title Section -->
        <div class="video-header">
          <h1 class="video-title">{{ opera.opera_title }}</h1>
          <div class="video-meta">
            <span><PlayCountIcon />{{ opera.play_count || 0 }}</span>
            <span class="divider"></span>
            <span>{{ formatTime(opera.created_at) }}</span>
          </div>
        </div>
        <!-- Video Player Section -->
        <!-- Container to hold space when player goes fixed -->
        <div class="player-placeholder" ref="playerBox" style="width: 100%; aspect-ratio: 16/9; margin-bottom: 16px;">
          <div class="video-player-wrapper" :class="{ 'pip-mode': isPip }">
            <div ref="dplayerContainer" class="dplayer-container"></div>
            <!-- PiP Close Button (Optional, only visible in PiP) -->
            <div v-if="isPip" class="pip-close-hint"
              style="position:absolute; top:5px; right:5px; background:rgba(0,0,0,0.5); color:white; padding:2px 5px; font-size:12px; border-radius:4px; pointer-events:none; z-index: 10001;">
              上滑关闭画中画
            </div>
          </div>
        </div>

        <!-- Toolbar -->
        <div class="video-toolbar">
          <div class="toolbar-left">
            <div class="action-item" :class="{ active: opera.liked }" @click="onToggleLike">
              <span class="action-icon">
                <LikeIcon />
              </span>
              <span>{{ (opera.like_count || 0) }}</span>
            </div>
            <div class="action-item" :class="{ active: opera.favorited }" @click="onToggleFavorite">
              <span class="action-icon">
                <FavIcon />
              </span>
              <span>{{ (opera.favorite_count || 0) }}</span>
            </div>
            <div class="action-item" @click="onShare">
              <span class="action-icon">
                <ShareIcon />
              </span>
              <span>{{ (opera.share_count || 0) }}</span>
            </div>
          </div>
        </div>

        <!-- Description -->
        <div class="video-desc">
          {{ opera.description || '暂无简介' }}
        </div>

        <!-- Comments -->
        <div class="comments-section">
          <h3 class="section-title">评论</h3>
          <!-- Mock Comments -->
          <div class="comment-item" style="display:flex; gap:12px; margin-bottom:20px;">
            <img src="https://ui-avatars.com/api/?name=游客" style="width:40px; height:40px; border-radius:50%;">
            <div>
              <div style="font-weight:500; font-size:13px; color:#61666d; margin-bottom:4px;">游客</div>
              <div style="font-size:14px;">哎呦，不错哦，发条评论吧</div>
            </div>
          </div>
        </div>
      </div>

      <!-- Right Column -->
      <div class="sidebar-column">
        <h3 class="section-title" style="font-size:16px; margin-bottom:12px;">接下来播放</h3>
        <div class="rec-list">
          <div class="rec-card" v-for="rec in recommendations" :key="rec.opera_id"
            @click="$router.push(`/video/${rec.opera_id}`)">
            <img :src="rec.avatar" class="rec-cover" loading="lazy" />
            <div class="rec-info">
              <div class="rec-title">{{ rec.opera_title }}</div>
              <div class="rec-meta">{{ formatArtists(rec.artists) }}</div>
              <div class="rec-meta">0播放</div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
  <div v-else class="loading-state">
    加载中...
  </div>
</template>

<script lang="ts" src="../scripts/views/PlayView.ts"></script>

<style scoped src="../styles/views/PlayView.css"></style>