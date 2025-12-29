<template>
  <div class="profile-view" v-if="user">
    <div class="profile-header">
      <div class="user-cover"></div>
      <div class="user-info-container">
        <div class="avatar-large">
          <img :src="user.icon || `https://ui-avatars.com/api/?name=${user.username}&background=random&size=128`"
            alt="Avatar" />
        </div>
        <div class="user-details">
          <div class="user-details-top">
            <h1 class="username">{{ user.username }}</h1>
            <div class="user-actions">
              <button class="edit-btn" @click="openEditModal">编辑资料</button>
              <button class="edit-btn secondary" @click="openPasswordModal">修改密码</button>
            </div>
          </div>
          <p class="user-meta">
            <span>UID: {{ user.user_id }}</span>
            <span class="divider">|</span>
            <span>{{ user.sex === 'Male' ? '男' : user.sex === 'Female' ? '女' : '保密' }}</span>
          </p>
          <p class="user-last-login" v-if="user.last_login_at">
            上次登录: {{ formatTime(user.last_login_at) }}
          </p>
          <p class="user-last-ip" v-if="user.last_login_ip">
            IP: {{ user.last_login_ip }}
          </p>
        </div>
      </div>
    </div>

    <div class="profile-tabs">
      <div class="tab-item" :class="{ active: activeTab === 'history' }" @click="activeTab = 'history'">
        播放历史
      </div>
      <div class="tab-item" :class="{ active: activeTab === 'likes' }" @click="activeTab = 'likes'">
        我的点赞
      </div>
      <div class="tab-item" :class="{ active: activeTab === 'favorites' }" @click="activeTab = 'favorites'">
        我的收藏
      </div>
    </div>

    <div class="tab-content">
      <div v-if="loading" class="loading-state">加载中...</div>
      <div v-else-if="list.length === 0" class="empty-state">暂无数据</div>
      <div v-else class="video-grid">
        <VideoCard v-for="opera in list" :key="opera.opera_id" :opera="opera"
          @click="$router.push(`/video/${opera.opera_id}`)" />
      </div>

      <Pagination v-if="!loading && totalItems > 0" :current-page="currentPage" :total-items="totalItems"
        :page-size="pageSize" @page-change="handlePageChange" />
    </div>

    <!-- Edit Profile Modal -->
    <div class="modal-overlay" v-if="showEditModal" @click.self="closeEditModal">
      <div class="modal-content">
        <div class="modal-header">
          <h3>编辑资料</h3>
          <span class="close-btn" @click="closeEditModal">&times;</span>
        </div>
        <div class="modal-body">
          <form @submit.prevent="handleUpdateProfile">
            <div class="form-group">
              <label>用户名</label>
              <input type="text" v-model="editForm.username" class="input-field" required />
            </div>
            <div class="form-group">
              <label>手机号</label>
              <input type="text" v-model="editForm.phone" class="input-field" required />
            </div>
            <div class="form-group">
              <label>性别</label>
              <select v-model="editForm.sex" class="input-field">
                <option value="Male">男</option>
                <option value="Female">女</option>
                <option value="Other">保密</option>
              </select>
            </div>
            
            <div class="form-group">
              <label>邮箱</label>
              <div class="input-with-button">
                 <input type="email" v-model="editForm.email" class="input-field" required />
                 <!-- Send code only if email changed -->
                 <button type="button" 
                         v-if="editForm.email !== user.email"
                         @click="sendVerifyCode" 
                         :disabled="codeCountdown > 0" 
                         class="send-code-btn">
                    {{ codeCountdown > 0 ? `${codeCountdown}s` : '发送验证码' }}
                 </button>
              </div>
              <small v-if="editForm.email !== user.email" style="color:#666; margin-top:4px; display:block;">
                修改邮箱需要验证。验证码将发送到<b>当前绑定</b>的邮箱。
              </small>
            </div>

            <div class="form-group" v-if="editForm.email !== user.email">
              <label>验证码</label>
              <input type="text" v-model="editForm.code" class="input-field" placeholder="请输入验证码" required />
            </div>

            <div class="error-msg" v-if="editError">{{ editError }}</div>
            <div class="success-msg" v-if="editSuccess">{{ editSuccess }}</div>

            <div class="modal-actions">
              <button type="button" class="cancel-btn" @click="closeEditModal">取消</button>
              <button type="submit" class="confirm-btn" :disabled="editLoading">保存</button>
            </div>
          </form>
        </div>
      </div>
    </div>
    <!-- Password Change Modal -->
    <div class="modal-overlay" v-if="showPasswordModal" @click.self="closePasswordModal">
      <div class="modal-content">
        <div class="modal-header">
          <h3>修改密码</h3>
          <span class="close-btn" @click="closePasswordModal">&times;</span>
        </div>
        <div class="modal-body">
          <form @submit.prevent="handleUpdatePassword">
            <div class="form-group">
              <label>原密码</label>
              <input type="password" v-model="passwordForm.old_password" class="input-field" placeholder="请输入原密码" required />
            </div>
            <div class="form-group">
              <label>新密码</label>
              <input type="password" v-model="passwordForm.new_password" class="input-field" placeholder="请输入新密码（至少6位）" minlength="6" required />
            </div>
            <div class="form-group">
              <label>确认新密码</label>
              <input type="password" v-model="passwordForm.confirm_password" class="input-field" placeholder="请再次输入新密码" minlength="6" required />
            </div>

            <div class="error-msg" v-if="passwordError">{{ passwordError }}</div>
            <div class="success-msg" v-if="passwordSuccess">{{ passwordSuccess }}</div>

            <div class="modal-actions">
              <button type="button" class="cancel-btn" @click="closePasswordModal">取消</button>
              <button type="submit" class="confirm-btn" :disabled="passwordLoading">修改</button>
            </div>
          </form>
        </div>
      </div>
    </div>
  </div>
  <div v-else-if="!loadingUser" class="not-logged-in">
    <p>请先登录查看个人中心</p>
    <button class="login-btn" @click="$router.push('/login')">去登录</button>
  </div>
</template>

<script lang="ts" src="../scripts/views/ProfileView.ts"></script>
<style scoped src="../styles/views/ProfileView.css"></style>
