<template>
  <div class="navbar-container">
    <div class="navbar-content">
      <!-- Left: Logo -->
      <div class="logo" @click="goHome">
        <span class="logo-text">黄梅戏数字化传播平台</span>
      </div>

      <!-- Center: Search (PC) -->
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
        <!-- 移动端搜索图标（仅移动端显示） -->
        <button class="mobile-search-icon" @click="goSearch" aria-label="搜索">
          <svg viewBox="0 0 24 24" width="20" height="20" stroke="currentColor" stroke-width="2" fill="none"
            stroke-linecap="round" stroke-linejoin="round">
            <circle cx="11" cy="11" r="8"></circle>
            <line x1="21" y1="21" x2="16.65" y2="16.65"></line>
          </svg>
        </button>

        <template v-if="isLoggedIn && user">
          <!--
            PC:  hover 触发 dropdown（纯 CSS）
            移动: click 触发 toggleMobileDropdown，不直接跳转
          -->
          <div class="avatar-wrapper" @click.stop="toggleMobileDropdown">
            <img :src="user.icon || `https://ui-avatars.com/api/?name=${user.username}&background=random`" alt="Avatar"
              class="avatar" />
            <!-- PC hover dropdown -->
            <div class="user-dropdown pc-dropdown">
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
          <!-- 收藏/历史 仅 PC 显示 -->
          <div class="action-item pc-only" @click="goCollection">
            <span>收藏</span>
          </div>
          <div class="action-item pc-only" @click="goHistory">
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

  <!-- 移动端 dropdown（fixed 定位，脱离 navbar 层叠上下文） -->
  <Teleport to="body">
    <Transition name="dropdown-fade">
      <div v-if="isLoggedIn && user && mobileDropdownOpen" class="mobile-dropdown-panel" @click.stop>
        <div class="mobile-dropdown-user">
          <img :src="user.icon || `https://ui-avatars.com/api/?name=${user.username}&background=random`" alt="Avatar"
            class="mobile-dropdown-avatar" />
          <div>
            <p class="mobile-dropdown-name">{{ user.username }}</p>
            <p class="mobile-dropdown-role">{{ user.role }}</p>
          </div>
        </div>
        <div class="dropdown-divider"></div>
        <div class="mobile-dropdown-item" @click="goProfile">个人中心</div>
        <div class="mobile-dropdown-item" v-if="user.role === 'Administrator'" @click="goAdmin">后台管理</div>
        <div class="mobile-dropdown-item danger" @click="handleLogout">退出登录</div>
      </div>
    </Transition>
    <!-- 点击外部关闭遮罩 -->
    <div v-if="mobileDropdownOpen" class="mobile-dropdown-overlay" @click="closeMobileDropdown"></div>
  </Teleport>
</template>

<script lang="ts">
import NavBarScript from '../scripts/components/NavBar';
export default NavBarScript;
</script>

<style scoped src="../styles/components/NavBar.css"></style>
