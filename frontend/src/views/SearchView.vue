<template>
  <div class="search-view">
    <!-- 搜索结果信息 -->
    <div v-if="hasSearched" class="search-info">
      <span class="result-count">找到 {{ totalResults }} 个结果</span>
      <span v-if="processingTime" class="processing-time">耗时: {{ processingTime }}ms</span>
    </div>

    <!-- 加载状态 -->
    <div v-if="loading" class="loading-container">
      <el-icon class="is-loading" :size="40"><Loading /></el-icon>
      <p>搜索中...</p>
    </div>

    <!-- 作品搜索结果 -->
    <div v-else-if="operaResults.length > 0" class="results-grid">
      <VideoCard
        v-for="opera in operaResults"
        :key="opera.opera_id"
        :opera="opera"
      />
    </div>

    <!-- 空状态 -->
    <el-empty
      v-if="hasSearched && totalResults === 0"
      description="没有找到相关结果"
      :image-size="200"
    />

    <!-- 分页 -->
    <Pagination
      v-if="totalResults > 0"
      :current-page="currentPage"
      :total-items="totalResults"
      :page-size="pageSize"
      @page-change="handlePageChange"
    />
  </div>
</template>

<script lang="ts" src="../scripts/views/SearchView.ts"></script>

<style scoped src="../styles/views/SearchView.css"></style>
