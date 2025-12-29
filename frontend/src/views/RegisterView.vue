<template>
  <div class="login-view">
    <div class="login-container">
      <h2 class="login-title">注册新账号</h2>
      <form @submit.prevent="handleRegister" class="login-form">
        <div class="form-group">
          <label for="username">用户名</label>
          <input type="text" id="username" v-model="username" placeholder="请输入用户名" required />
        </div>
        <div class="form-group">
          <label for="phone">手机号</label>
          <input type="tel" id="phone" v-model="phone" placeholder="请输入手机号" required />
        </div>
        <div class="form-group">
          <label for="email">邮箱</label>
          <div class="input-with-btn">
            <input type="email" id="email" v-model="email" placeholder="请输入邮箱" required />
            <button type="button" class="send-code-btn" :disabled="sending || !!countdown" @click="sendCode">
              {{ countdown ? `${countdown}s` : '获取验证码' }}
            </button>
          </div>
        </div>
        <div class="form-group">
          <label for="code">验证码</label>
          <input type="text" id="code" v-model="code" placeholder="请输入验证码" required />
        </div>
        <div class="form-group">
          <label for="password">密码</label>
          <input type="password" id="password" v-model="password" placeholder="至少6位密码" required minlength="6" />
        </div>
        <div v-if="error" class="error-message">
          {{ error }}
        </div>
        <div v-if="success" class="success-message">
          {{ success }}
        </div>
        <button type="submit" class="login-submit-btn" :disabled="loading">
          {{ loading ? '注册中...' : '注册' }}
        </button>
      </form>
      <div class="login-footer">
        <p>已有账号？ <router-link to="/login">立即登录</router-link></p>
      </div>
    </div>
  </div>
</template>

<script lang="ts" src="../scripts/views/RegisterView.ts"></script>
<style scoped src="../styles/views/LoginView.css"></style>
<style scoped src="../styles/views/RegisterView.css"></style>
