<template>
  <div class="not-found-container">
    <div class="not-found-content">
      <h1 class="error-code">404</h1>
      <p class="error-message">抱歉，您访问的页面不存在</p>
      <button class="home-btn" @click="$router.push('/')">返回首页</button>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue';

export default defineComponent({
  name: 'NotFound'
});
</script>

<style scoped>
.not-found-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 60vh;
  text-align: center;
  color: #18191c;
}

.error-code {
  font-size: 80px;
  font-weight: bold;
  color: #A40000;
  margin: 0;
  line-height: 1;
}

.error-message {
  font-size: 18px;
  color: #61666d;
  margin: 20px 0 30px;
}

.home-btn {
  background-color: #A40000;
  color: white;
  border: none;
  padding: 10px 24px;
  border-radius: 4px;
  font-size: 14px;
  cursor: pointer;
  transition: background-color 0.2s;
}

.home-btn:hover {
  background-color: #c90000;
}
</style>
