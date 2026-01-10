<template>
  <div class="admin-layout">
    <el-container>
      <el-aside width="200px">
        <div class="logo">管理后台</div>
        <el-menu router :default-active="$route.path" class="el-menu-vertical" background-color="#304156"
          text-color="#bfcbd9" active-text-color="#409EFF">
          <el-menu-item index="/admin/dashboard">
            <el-icon>
              <DataLine />
            </el-icon>
            <span>仪表盘</span>
          </el-menu-item>
          <el-menu-item index="/admin/videos">
            <el-icon>
              <VideoCamera />
            </el-icon>
            <span>视频管理</span>
          </el-menu-item>
          <el-menu-item index="/admin/tasks">
            <el-icon>
              <List />
            </el-icon>
            <span>任务队列</span>
          </el-menu-item>
          <el-menu-item index="/admin/artists">
            <el-icon>
              <User />
            </el-icon>
            <span>艺术家管理</span>
          </el-menu-item>
          <el-menu-item index="/admin/users">
            <el-icon>
              <UserFilled />
            </el-icon>
            <span>用户管理</span>
          </el-menu-item>
          <el-menu-item index="/">
            <el-icon>
              <HomeFilled />
            </el-icon>
            <span>返回前台</span>
          </el-menu-item>
        </el-menu>
      </el-aside>
      <el-container>
        <el-header>
          <div class="header-content">
            <el-breadcrumb separator="/">
              <el-breadcrumb-item :to="{ path: '/admin/dashboard' }">首页</el-breadcrumb-item>
              <el-breadcrumb-item>{{ $route.name }}</el-breadcrumb-item>
            </el-breadcrumb>
            <div class="user-info">
              <el-dropdown trigger="click">
                <div class="avatar-container">
                  <img :src="user?.icon || `https://ui-avatars.com/api/?name=${user?.username || 'Admin'}&background=random`" class="user-avatar" />
                  <span class="username">{{ user?.username }}</span>
                  <el-icon class="el-icon--right"><CaretBottom /></el-icon>
                </div>
                <template #dropdown>
                  <el-dropdown-menu>
                    <el-dropdown-item @click="logout">退出登录</el-dropdown-item>
                  </el-dropdown-menu>
                </template>
              </el-dropdown>
            </div>
          </div>
        </el-header>
        <el-main>
          <router-view />
        </el-main>
      </el-container>
    </el-container>
  </div>
</template>

<script lang="ts">
import AdminLayoutScript from '../../scripts/views/admin/AdminLayout';
import { DataLine, VideoCamera, User, UserFilled, HomeFilled, CaretBottom, List } from '@element-plus/icons-vue';

export default {
  ...AdminLayoutScript,
  components: {
    DataLine,
    VideoCamera,
    User,
    UserFilled,
    HomeFilled,
    CaretBottom,
    List
  }
};
</script>

<style scoped src="../../styles/views/admin/AdminLayout.css"></style>
