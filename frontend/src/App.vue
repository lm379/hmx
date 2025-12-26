<template>
  <div class="app-container">
    <NavBar />
    
    <main class="main-content">
      <!-- Category/Filter Bar (Global for specific pages) -->
      <div class="category-bar-container" v-if="$route.meta.showCategoryBar">
        <div class="category-bar">
          <button 
            class="cat-btn" 
            :class="{ active: $route.path === '/' }"
            @click="$router.push('/')"
          >
            个性推荐
          </button>
          <button 
            class="cat-btn" 
            :class="{ active: $route.path === '/artists' }"
            @click="$router.push('/artists')"
          >
            名家风采
          </button>
          <button class="cat-btn">热门视频</button>
          <button class="cat-btn">排行</button>
        </div>
      </div>

      <router-view />
    </main>
  </div>
</template>

<script lang="ts" src="./scripts/App.ts"></script>
<style src="./styles/App.css"></style>