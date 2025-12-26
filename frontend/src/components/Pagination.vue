<template>
  <div v-if="totalItems > 0" class="pagination">
    <button 
      class="page-btn" 
      :disabled="currentPage === 1" 
      @click="handlePageChange(currentPage - 1)"
    >
      上一页
    </button>
    <span class="page-info">第 {{ currentPage }} 页 / 共 {{ totalPages }} 页</span>
    <button 
      class="page-btn" 
      :disabled="currentPage === totalPages" 
      @click="handlePageChange(currentPage + 1)"
    >
      下一页
    </button>
  </div>
</template>

<script lang="ts" src="../scripts/components/Pagination.ts"></script>

<style scoped src="../styles/components/Pagination.css"></style>
