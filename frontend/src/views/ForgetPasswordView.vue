<template>
  <div class="login-view">
    <div class="login-container">
      <h2 class="login-title">重置密码</h2>
      <form @submit.prevent="handleForgetPassword" class="login-form">
        <div class="form-group">
          <label for="email">邮箱</label>
          <div class="input-with-btn">
            <input type="email" id="email" v-model="email" placeholder="请输入绑定的邮箱" required />
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
          <label for="new_password">新密码</label>
          <input type="password" id="new_password" v-model="new_password" placeholder="至少6位新密码" required minlength="6" />
        </div>
        <div v-if="error" class="error-message">
          {{ error }}
        </div>
        <div v-if="success" class="success-message">
          {{ success }}
        </div>
        <button type="submit" class="login-submit-btn" :disabled="loading">
          {{ loading ? '重置中...' : '重置密码' }}
        </button>
      </form>
      <div class="login-footer">
        <p><router-link to="/login">返回登录</router-link></p>
      </div>
    </div>
  </div>
</template>

<script lang="ts" src="../scripts/views/ForgetPasswordView.ts"></script>
<style scoped src="../styles/views/LoginView.css"></style>
<style scoped src="../styles/views/ForgetPasswordView.css"></style>
