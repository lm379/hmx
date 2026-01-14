<template>
  <div class="search-manager">
    <el-card class="header-card">
      <template #header>
        <div class="card-header">
          <h2>搜索索引管理</h2>
          <el-button
            type="primary"
            :icon="Refresh"
            @click="loadStats"
            :loading="statsLoading"
          >
            刷新统计
          </el-button>
        </div>
      </template>

      <!-- 统计信息 -->
      <el-row :gutter="20">
        <el-col :span="12">
          <el-statistic title="作品索引数量" :value="operaCount">
            <template #prefix>
              <el-icon><VideoPlay /></el-icon>
            </template>
          </el-statistic>
          <div class="index-status">
            <el-tag :type="operaIndexing ? 'warning' : 'success'" size="small">
              {{ operaIndexing ? '索引中...' : '就绪' }}
            </el-tag>
          </div>
        </el-col>
        <el-col :span="12">
          <el-statistic title="艺术家索引数量" :value="artistCount">
            <template #prefix>
              <el-icon><User /></el-icon>
            </template>
          </el-statistic>
          <div class="index-status">
            <el-tag :type="artistIndexing ? 'warning' : 'success'" size="small">
              {{ artistIndexing ? '索引中...' : '就绪' }}
            </el-tag>
          </div>
        </el-col>
      </el-row>
    </el-card>

    <!-- 索引操作 -->
    <el-card class="action-card">
      <template #header>
        <h3>索引操作</h3>
      </template>

      <el-alert
        title="提示"
        type="info"
        :closable="false"
        style="margin-bottom: 20px;"
      >
        重新索引会将数据库中的所有数据同步到 Meilisearch。通常在以下情况下使用：
        <ul style="margin: 8px 0 0 20px;">
          <li>首次部署应用时</li>
          <li>数据同步出现异常时</li>
          <li>搜索结果与数据库不一致时</li>
        </ul>
      </el-alert>

      <el-row :gutter="20">
        <el-col :span="12">
          <el-card shadow="hover" class="operation-card">
            <div class="operation-content">
              <el-icon class="operation-icon" :size="48"><VideoPlay /></el-icon>
              <h4>重新索引作品</h4>
              <p>将所有作品数据同步到 Meilisearch</p>
              <el-button
                type="primary"
                :loading="reindexingOperas"
                @click="handleReindexOperas"
                style="margin-top: 16px;"
              >
                {{ reindexingOperas ? '索引中...' : '开始索引' }}
              </el-button>
            </div>
          </el-card>
        </el-col>

        <el-col :span="12">
          <el-card shadow="hover" class="operation-card">
            <div class="operation-content">
              <el-icon class="operation-icon" :size="48"><User /></el-icon>
              <h4>重新索引艺术家</h4>
              <p>将所有艺术家数据同步到 Meilisearch</p>
              <el-button
                type="primary"
                :loading="reindexingArtists"
                @click="handleReindexArtists"
                style="margin-top: 16px;"
              >
                {{ reindexingArtists ? '索引中...' : '开始索引' }}
              </el-button>
            </div>
          </el-card>
        </el-col>
      </el-row>

      <el-divider />

      <!-- 批量操作 -->
      <div class="batch-operations">
        <h4>批量操作</h4>
        <el-button
          type="warning"
          :loading="reindexingAll"
          @click="handleReindexAll"
        >
          重新索引全部数据
        </el-button>
      </div>
    </el-card>

    <!-- 索引详情 -->
    <el-card v-if="operaFieldDistribution || artistFieldDistribution">
      <template #header>
        <h3>索引字段分布</h3>
      </template>

      <el-tabs>
        <el-tab-pane label="作品索引">
          <el-descriptions :column="2" border v-if="operaFieldDistribution">
            <el-descriptions-item
              v-for="(value, key) in operaFieldDistribution"
              :key="key"
              :label="key"
            >
              {{ value }}
            </el-descriptions-item>
          </el-descriptions>
          <el-empty v-else description="暂无数据" />
        </el-tab-pane>

        <el-tab-pane label="艺术家索引">
          <el-descriptions :column="2" border v-if="artistFieldDistribution">
            <el-descriptions-item
              v-for="(value, key) in artistFieldDistribution"
              :key="key"
              :label="key"
            >
              {{ value }}
            </el-descriptions-item>
          </el-descriptions>
          <el-empty v-else description="暂无数据" />
        </el-tab-pane>
      </el-tabs>
    </el-card>
  </div>
</template>

<script lang="ts" src="../../scripts/views/admin/SearchManager.ts"></script>

<style scoped src="../../styles/views/admin/SearchManager.css"></style>