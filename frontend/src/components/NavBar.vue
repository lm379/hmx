<template>
  <div class="navbar-container">
    <div class="navbar-content">
      <!-- Left: Logo -->
      <div class="logo" @click="goHome">
        <span class="logo-text">黄梅戏数字化传播平台</span>
      </div>

      <!-- Center: Search -->
      <div class="search-bar">
        <input 
          type="text" 
          v-model="searchQuery"
          @keyup.enter="handleSearch"
          placeholder="搜索你感兴趣的内容" 
        />
        <button class="search-btn" @click="handleSearch">
          <svg viewBox="0 0 24 24" width="18" height="18" stroke="currentColor" stroke-width="2" fill="none"
            stroke-linecap="round" stroke-linejoin="round">
            <circle cx="11" cy="11" r="8"></circle>
            <line x1="21" y1="21" x2="16.65" y2="16.65"></line>
          </svg>
        </button>
      </div>

      <!-- Right: User Actions -->
      <div class="user-actions">
        <template v-if="isLoggedIn && user">
          <div class="avatar-wrapper" @click="goProfile">
            <img :src="user.icon || `https://ui-avatars.com/api/?name=${user.username}&background=random`" alt="Avatar"
              class="avatar" />
            <div class="user-dropdown">
              <div class="user-info-brief" @click.stop="goProfile">
                <p class="username">{{ user.username }}</p>
                <p class="user-role">{{ user.role }}</p>
              </div>
              <div class="dropdown-divider"></div>
              <div class="dropdown-item" @click.stop="goProfile">个人中心</div>
              <div class="dropdown-item" @click.stop="goAdmin" v-if="user.role === 'Administrator'">后台管理</div>
              <div class="dropdown-item" @click.stop="handleLogout">退出登录</div>
            </div>
          </div>
          <div class="action-item" @click="goCollection">
            <span>收藏</span>
          </div>
          <div class="action-item" @click="goHistory">
            <span>历史</span>
          </div>
        </template>
        <template v-else>
          <div class="login-trigger" @click="goLogin">
            <div class="avatar-placeholder">
              <svg viewBox="0 0 24 24" width="20" height="20" fill="#61666d">
                <path
                  d="M12 12c2.21 0 4-1.79 4-4s-1.79-4-4-4-4 1.79-4 4 1.79 4 4 4zm0 2c-2.67 0-8 1.34-8 4v2h16v-2c0-2.66-5.33-4-8-4z">
                </path>
              </svg>
            </div>
            <span class="login-text">登录</span>
          </div>
        </template>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import NavBarScript from '../scripts/components/NavBar';
export default NavBarScript;
</script>

<style scoped src="../styles/components/NavBar.css"></style>
