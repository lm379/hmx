<template>
  <div class="not-found-container">
    <div class="not-found-content">
      <h1 class="error-code">404</h1>
      <p class="error-message">抱歉，您访问的页面不存在</p>
      <button class="home-btn" @click="$router.push('/')">返回首页</button>
    </div>
  </div>
</template>

<script lang="ts" src="../scripts/views/NotFound.ts"></script>

<style scoped src="../styles/views/NotFound.css"></style>