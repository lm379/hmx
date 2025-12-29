<template>
  <div class="login-view">
    <div class="login-container">
      <h2 class="login-title">登录</h2>
      <form @submit.prevent="handleLogin" class="login-form">
        <div class="form-group">
          <label for="account">账号</label>
          <input 
            type="text" 
            id="account" 
            v-model="account" 
            placeholder="请输入用户名或邮箱" 
            required 
          />
        </div>
        <div class="form-group">
          <label for="password">密码</label>
          <input 
            type="password" 
            id="password" 
            v-model="password" 
            placeholder="请输入密码" 
            required 
          />
        </div>
        <div class="form-footer-actions">
          <label class="remember-me">
            <input type="checkbox" v-model="rememberMe" />
            <span>记住密码</span>
          </label>
          <router-link to="/forget-password" class="forget-pwd-link">忘记密码？</router-link>
        </div>
        <div v-if="error" class="error-message">
          {{ error }}
        </div>
        <button type="submit" class="login-submit-btn" :disabled="loading">
          {{ loading ? '登录中...' : '登录' }}
        </button>
      </form>
      <div class="login-footer">
        <p>还没有账号？ <router-link to="/register">立即注册</router-link></p>
      </div>
    </div>
  </div>
</template>

<script lang="ts" src="../scripts/views/LoginView.ts"></script>
<style scoped src="../styles/views/LoginView.css"></style>
