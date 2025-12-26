<template>
  <div class="profile-view" v-if="user">
    <div class="profile-header">
      <div class="user-cover"></div>
      <div class="user-info-container">
        <div class="avatar-large">
          <img :src="user.icon || `https://ui-avatars.com/api/?name=${user.username}&background=random&size=128`" alt="Avatar" />
        </div>
        <div class="user-details">
          <h1 class="username">{{ user.username }}</h1>
          <p class="user-meta">
            <span>UID: {{ user.user_id }}</span>
            <span class="divider">|</span>
            <span>角色: {{ user.role }}</span>
          </p>
          <p class="user-last-login" v-if="user.last_login_at">
            上次登录: {{ formatTime(user.last_login_at) }}
          </p>
        </div>
      </div>
    </div>

    <div class="profile-tabs">
      <div 
        class="tab-item" 
        :class="{ active: activeTab === 'history' }" 
        @click="activeTab = 'history'"
      >
        播放历史
      </div>
      <div 
        class="tab-item" 
        :class="{ active: activeTab === 'likes' }" 
        @click="activeTab = 'likes'"
      >
        我的点赞
      </div>
      <div 
        class="tab-item" 
        :class="{ active: activeTab === 'favorites' }" 
        @click="activeTab = 'favorites'"
      >
        我的收藏
      </div>
    </div>

    <div class="tab-content">
      <div v-if="loading" class="loading-state">加载中...</div>
      <div v-else-if="list.length === 0" class="empty-state">暂无数据</div>
      <div v-else class="video-grid">
        <VideoCard 
          v-for="opera in list" 
          :key="opera.opera_id" 
          :opera="opera" 
          @click="$router.push(`/video/${opera.opera_id}`)"
        />
      </div>

      <Pagination
        v-if="!loading && totalItems > 0"
        :current-page="currentPage"
        :total-items="totalItems"
        :page-size="pageSize"
        @page-change="handlePageChange"
      />
    </div>
  </div>
  <div v-else-if="!loadingUser" class="not-logged-in">
    <p>请先登录查看个人中心</p>
    <button class="login-btn" @click="$router.push('/login')">去登录</button>
  </div>
</template>

<script lang="ts" src="../scripts/views/ProfileView.ts"></script>
<style scoped src="../styles/views/ProfileView.css"></style>
