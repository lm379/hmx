<template>
  <div class="dashboard-view">
    <el-row :gutter="20">
      <el-col :span="6">
        <el-card shadow="hover">
          <template #header>用户总数</template>
          <div class="stat-value">{{ stats.total_users }}</div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover">
          <template #header>视频总数</template>
          <div class="stat-value">{{ stats.total_operas }}</div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover">
          <template #header>艺术家总数</template>
          <div class="stat-value">{{ stats.total_artists }}</div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover">
          <template #header>总播放量</template>
          <div class="stat-value">{{ stats.total_plays }}</div>
        </el-card>
      </el-col>
    </el-row>
    <el-row :gutter="20" style="margin-top: 20px;">
      <el-col :span="8">
         <el-card shadow="hover">
          <template #header>总点赞数</template>
          <div class="stat-value">{{ stats.total_likes }}</div>
        </el-card>
      </el-col>
      <el-col :span="8">
         <el-card shadow="hover">
          <template #header>总收藏数</template>
          <div class="stat-value">{{ stats.total_favorites }}</div>
        </el-card>
      </el-col>
      <el-col :span="8">
         <el-card shadow="hover">
          <template #header>总分享数</template>
          <div class="stat-value">{{ stats.total_shares }}</div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script lang="ts">
import DashboardViewScript from '../../scripts/views/admin/DashboardView';
export default DashboardViewScript;
</script>

<style scoped src="../../styles/views/admin/DashboardView.css"></style>
