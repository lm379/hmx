<template>
  <div class="app-container">
    <NavBar v-if="!$route.path.startsWith('/admin')" />

    <main :class="{ 'main-content': !$route.path.startsWith('/admin') }">
      <!-- Category/Filter Bar (Global for specific pages) -->
      <div class="category-bar-container" v-if="$route.meta.showCategoryBar">
        <div class="category-bar">
          <button class="cat-btn" :class="{ active: $route.path === '/' }" @click="$router.push('/')">
            个性推荐
          </button>
          <button class="cat-btn" :class="{ active: $route.path === '/artists' }" @click="$router.push('/artists')">
            名家风采
          </button>
          <button class="cat-btn" :class="{ active: $route.path === '/ranking' }" @click="$router.push('/ranking')">
            热门排行
          </button>
          <button class="cat-btn" :class="{ active: $route.path.startsWith('/news') }" @click="$router.push('/news')">
            戏曲资讯
          </button>
        </div>
      </div>

      <router-view />
    </main>

    <!-- Global knowledge Q&A chat widget (hidden in admin panel) -->
    <ChatWindow v-if="!$route.path.startsWith('/admin')" />
  </div>
</template>

<script lang="ts" src="./scripts/App.ts"></script>
<style src="./styles/App.css"></style>
