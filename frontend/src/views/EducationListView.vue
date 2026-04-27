<template>
  <div class="education-list-view">
    <div v-if="loading" class="loading-state">加载中...</div>
    <div v-else-if="books.length === 0" class="empty-state">暂无教育资源</div>

    <div v-else class="book-grid">
      <article v-for="book in books" :key="book.book_id" class="book-card" @click="goReader(book.book_id)">
        <div class="book-cover">
          <img v-if="book.cover_url" :src="book.cover_url" :alt="book.title" loading="lazy" />
          <div v-else class="book-icon">PDF</div>
        </div>
        <div class="book-info">
          <h2 class="book-title">{{ book.title }}</h2>
        </div>
      </article>
    </div>

    <Pagination
      v-if="!loading && totalItems > 0"
      :current-page="currentPage"
      :total-items="totalItems"
      :page-size="pageSize"
      @page-change="handlePageChange"
    />
  </div>
</template>

<script lang="ts" src="../scripts/views/EducationListView.ts"></script>
<style scoped src="../styles/views/EducationListView.css"></style>
