<template>
  <div class="play-view" v-if="!loading && opera">
    <div class="play-container">
      <!-- Left Column -->
      <div class="main-column">
        <!-- Title Section -->
        <div class="video-header">
          <h1 class="video-title">{{ opera.opera_title }}</h1>
          <div class="video-meta">
            <span>{{ 0 }}播放</span>
            <span>{{ 0 }}点赞</span>
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
            <div class="action-item">
              <span class="action-icon">
                <svg width="36" height="36" viewBox="0 0 36 36" xmlns="http://www.w3.org/2000/svg"
                  class="video-like-icon video-toolbar-item-icon">
                  <path fill-rule="evenodd" clip-rule="evenodd"
                    d="M9.77234 30.8573V11.7471H7.54573C5.50932 11.7471 3.85742 13.3931 3.85742 15.425V27.1794C3.85742 29.2112 5.50932 30.8573 7.54573 30.8573H9.77234ZM11.9902 30.8573V11.7054C14.9897 10.627 16.6942 7.8853 17.1055 3.33591C17.2666 1.55463 18.9633 0.814421 20.5803 1.59505C22.1847 2.36964 23.243 4.32583 23.243 6.93947C23.243 8.50265 23.0478 10.1054 22.6582 11.7471H29.7324C31.7739 11.7471 33.4289 13.402 33.4289 15.4435C33.4289 15.7416 33.3928 16.0386 33.3215 16.328L30.9883 25.7957C30.2558 28.7683 27.5894 30.8573 24.528 30.8573H11.9911H11.9902Z"
                    fill="currentColor"></path>
                </svg>
              </span>
              <span>{{ 0 }}</span>
            </div>
            <div class="action-item">
              <span class="action-icon">
                <svg width="28" height="28" viewBox="0 0 28 28" xmlns="http://www.w3.org/2000/svg"
                  class="video-fav-icon video-toolbar-item-icon" data-v-5f171eb7="">
                  <path fill-rule="evenodd" clip-rule="evenodd"
                    d="M19.8071 9.26152C18.7438 9.09915 17.7624 8.36846 17.3534 7.39421L15.4723 3.4972C14.8998 2.1982 13.1004 2.1982 12.4461 3.4972L10.6468 7.39421C10.1561 8.36846 9.25639 9.09915 8.19315 9.26152L3.94016 9.91102C2.63155 10.0734 2.05904 11.6972 3.04049 12.6714L6.23023 15.9189C6.96632 16.6496 7.29348 17.705 7.1299 18.7605L6.39381 23.307C6.14844 24.6872 7.62063 25.6614 8.84745 25.0119L12.4461 23.0634C13.4276 22.4951 14.6544 22.4951 15.6359 23.0634L19.2345 25.0119C20.4614 25.6614 21.8518 24.6872 21.6882 23.307L20.8703 18.7605C20.7051 17.705 21.0339 16.6496 21.77 15.9189L24.9597 12.6714C25.9412 11.6972 25.3687 10.0734 24.06 9.91102L19.8071 9.26152Z"
                    fill="currentColor"></path>
                </svg>
              </span>
              <span>{{ 0 }}</span>
            </div>
            <div class="action-item">
              <span class="action-icon">
                <svg data-v-12f7cbf0="" width="28" height="28" viewBox="0 0 28 28" xmlns="http://www.w3.org/2000/svg"
                  class="video-share-icon video-toolbar-item-icon">
                  <path
                    d="M12.6058 10.3326V5.44359C12.6058 4.64632 13.2718 4 14.0934 4C14.4423 4 14.78 4.11895 15.0476 4.33606L25.3847 12.7221C26.112 13.3121 26.2087 14.3626 25.6007 15.0684C25.5352 15.1443 25.463 15.2144 25.3847 15.2779L15.0476 23.6639C14.4173 24.1753 13.4791 24.094 12.9521 23.4823C12.7283 23.2226 12.6058 22.8949 12.6058 22.5564V18.053C7.59502 18.053 5.37116 19.9116 2.57197 23.5251C2.47607 23.6489 2.00031 23.7769 2.00031 23.2122C2.00031 16.2165 3.90102 10.3326 12.6058 10.3326Z"
                    fill="currentColor"></path>
                </svg>
              </span>
              <span>{{ 0 }}</span>
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