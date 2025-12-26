<template>
  <div class="navbar-container">
    <div class="navbar-content">
      <!-- Left: Logo -->
      <div class="logo" @click="goHome">
        <span class="logo-text">黄梅戏数字化传播平台</span>
      </div>

      <!-- Center: Search -->
      <div class="search-bar">
        <input type="text" placeholder="搜索感兴趣的内容" />
        <button class="search-btn">
          <svg viewBox="0 0 24 24" width="18" height="18" stroke="currentColor" stroke-width="2" fill="none" stroke-linecap="round" stroke-linejoin="round"><circle cx="11" cy="11" r="8"></circle><line x1="21" y1="21" x2="16.65" y2="16.65"></line></svg>
        </button>
      </div>

      <!-- Right: User Actions -->
      <div class="user-actions">
        <div class="avatar-wrapper">
           <img src="https://ui-avatars.com/api/?name=User&background=random" alt="Avatar" class="avatar" />
        </div>
        <!-- <div class="action-item">
           <span>动态</span>
        </div> -->
        <div class="action-item">
           <span>收藏</span>
        </div>
        <div class="action-item">
           <span>历史</span>
        </div>
        <button class="upload-btn">
          <span>投稿</span>
        </button>
      </div>
    </div>
  </div>
</template>

<script lang="ts" src="../scripts/components/NavBar.ts"></script>

<style scoped src="../styles/components/NavBar.css"></style>
