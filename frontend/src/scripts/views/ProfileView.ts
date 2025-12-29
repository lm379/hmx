import { defineComponent, ref, onMounted, watch, reactive } from 'vue';
import { storeToRefs } from 'pinia';
import axios from 'axios';
import VideoCard from '../../components/VideoCard.vue';
import Pagination from '../../components/Pagination.vue';
import { useAuthStore } from '../../stores/auth';
import type { Opera, UpdateUserProfileRequest } from '../../types';

export default defineComponent({
  name: 'ProfileView',
  components: {
    VideoCard,
    Pagination
  },
  setup() {
    const authStore = useAuthStore();
    const { user } = storeToRefs(authStore);
    const loadingUser = ref(true);
    const activeTab = ref('history');
    const list = ref<Opera[]>([]);
    const loading = ref(false);
    
    const currentPage = ref(1);
    const pageSize = ref(10);
    const totalItems = ref(0);

    // Edit Modal State
    const showEditModal = ref(false);
    const editLoading = ref(false);
    const editError = ref('');
    const editSuccess = ref('');
    const codeCountdown = ref(0);
    
    const editForm = reactive<UpdateUserProfileRequest>({
        username: '',
        phone: '',
        sex: 'Other',
        email: '',
        code: ''
    });

    const fetchUser = async () => {
      try {
        loadingUser.value = true;
        await authStore.checkLoginStatus();
        if (user.value) {
          fetchTabData();
        }
      } catch (error) {
        console.error("Failed to fetch user:", error);
      } finally {
        loadingUser.value = false;
      }
    };

    const fetchTabData = async (page = 1) => {
      try {
        loading.value = true;
        let endpoint = '';
        switch (activeTab.value) {
          case 'history': endpoint = '/api/v1/users/me/history'; break;
          case 'likes': endpoint = '/api/v1/users/me/likes'; break;
          case 'favorites': endpoint = '/api/v1/users/me/favorites'; break;
        }

        const response = await axios.get(endpoint, {
          params: {
            page,
            page_size: pageSize.value
          }
        });

        if (response.data && response.data.data) {
          list.value = response.data.data.list || [];
          if (response.data.data.pagination) {
            totalItems.value = response.data.data.pagination.total;
            currentPage.value = response.data.data.pagination.page;
          }
        }
      } catch (error) {
        console.error("Failed to fetch tab data:", error);
        list.value = [];
        totalItems.value = 0;
      } finally {
        loading.value = false;
      }
    };

    const handlePageChange = (newPage: number) => {
      fetchTabData(newPage);
      window.scrollTo({ top: 0, behavior: 'smooth' });
    };

    watch(activeTab, () => {
      currentPage.value = 1;
      fetchTabData();
    });

    onMounted(() => {
      fetchUser();
    });

    const formatTime = (time: string) => {
      if (!time) return '';
      const date = new Date(time);
      const year = date.getFullYear();
      const month = String(date.getMonth() + 1).padStart(2, '0');
      const day = String(date.getDate()).padStart(2, '0');
      const hours = String(date.getHours()).padStart(2, '0');
      const minutes = String(date.getMinutes()).padStart(2, '0');
      const seconds = String(date.getSeconds()).padStart(2, '0');
      return `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`;
    };

    // --- Edit Profile Logic ---

    const openEditModal = () => {
        if (!user.value) return;
        editForm.username = user.value.username;
        editForm.phone = user.value.phone;
        editForm.sex = user.value.sex;
        editForm.email = user.value.email || '';
        editForm.code = '';
        editError.value = '';
        editSuccess.value = '';
        showEditModal.value = true;
    };

    const closeEditModal = () => {
        showEditModal.value = false;
    };

    const sendVerifyCode = async () => {
        if (codeCountdown.value > 0) return;
        
        try {
            editError.value = '';
            // Send to CURRENT bound email
            await axios.post('/api/v1/users/me/email-code');
            codeCountdown.value = 60;
            const timer = setInterval(() => {
                codeCountdown.value--;
                if (codeCountdown.value <= 0) clearInterval(timer);
            }, 1000);
        } catch (e: any) {
            editError.value = e.response?.data?.error || '发送验证码失败';
        }
    };

    const handleUpdateProfile = async () => {
        editLoading.value = true;
        editError.value = '';
        editSuccess.value = '';
        
        try {
            // Only send code if email is changing
            const payload: any = { ...editForm };
            if (user.value?.email === editForm.email) {
                delete payload.code;
            }

            await axios.put('/api/v1/users/me', payload);
            editSuccess.value = '资料更新成功';
            
            // Refresh user data
            await authStore.checkLoginStatus();
            
            setTimeout(() => {
                closeEditModal();
            }, 1500);
        } catch (e: any) {
            editError.value = e.response?.data?.error || '更新失败';
        } finally {
            editLoading.value = false;
        }
    };

    // --- Change Password Logic ---
    const showPasswordModal = ref(false);
    const passwordLoading = ref(false);
    const passwordError = ref('');
    const passwordSuccess = ref('');
    const passwordForm = reactive({
        old_password: '',
        new_password: '',
        confirm_password: ''
    });

    const openPasswordModal = () => {
        passwordForm.old_password = '';
        passwordForm.new_password = '';
        passwordForm.confirm_password = '';
        passwordError.value = '';
        passwordSuccess.value = '';
        showPasswordModal.value = true;
    };

    const closePasswordModal = () => {
        showPasswordModal.value = false;
    };

    const handleUpdatePassword = async () => {
        if (passwordForm.new_password !== passwordForm.confirm_password) {
            passwordError.value = '两次输入的密码不一致';
            return;
        }

        passwordLoading.value = true;
        passwordError.value = '';
        passwordSuccess.value = '';

        try {
            await axios.post('/api/v1/users/me/password', {
                old_password: passwordForm.old_password,
                new_password: passwordForm.new_password
            });

            passwordSuccess.value = '密码修改成功，请重新登录';
            setTimeout(() => {
                authStore.logout();
                window.location.href = '/login'; // Force reload to clear state
            }, 1500);
        } catch (e: any) {
            passwordError.value = e.response?.data?.error || '修改失败';
        } finally {
            passwordLoading.value = false;
        }
    };

    return {
      user,
      loadingUser,
      activeTab,
      list,
      loading,
      currentPage,
      pageSize,
      totalItems,
      handlePageChange,
      formatTime,
      // Edit Profile
      showEditModal,
      editLoading,
      editError,
      editSuccess,
      editForm,
      codeCountdown,
      openEditModal,
      closeEditModal,
      sendVerifyCode,
      handleUpdateProfile,
      // Change Password
      showPasswordModal,
      passwordLoading,
      passwordError,
      passwordSuccess,
      passwordForm,
      openPasswordModal,
      closePasswordModal,
      handleUpdatePassword
    };
  }
});
