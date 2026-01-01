import { defineComponent, ref, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import axios from 'axios';
import { useAuthStore } from '../../stores/auth';

const REMEMBER_ME_KEY = 'remember_me';
const SAVED_ACCOUNT_KEY = 'saved_account';
const SAVED_PASSWORD_KEY = 'saved_password';

export default defineComponent({
  name: 'LoginView',
  setup() {
    const router = useRouter();
    const authStore = useAuthStore();
    const account = ref('');
    const password = ref('');
    const rememberMe = ref(false);
    const loading = ref(false);
    const error = ref('');

    // 加载已保存的凭据
    const loadSavedCredentials = () => {
      const remembered = localStorage.getItem(REMEMBER_ME_KEY) === 'true';
      if (remembered) {
        const savedAccount = localStorage.getItem(SAVED_ACCOUNT_KEY);
        const savedPassword = localStorage.getItem(SAVED_PASSWORD_KEY);
        if (savedAccount) account.value = savedAccount;
        if (savedPassword) password.value = savedPassword;
        rememberMe.value = true;
      }
    };

    // 保存或清除凭据
    const handleCredentials = () => {
      if (rememberMe.value) {
        localStorage.setItem(REMEMBER_ME_KEY, 'true');
        localStorage.setItem(SAVED_ACCOUNT_KEY, account.value);
        localStorage.setItem(SAVED_PASSWORD_KEY, password.value);
      } else {
        localStorage.removeItem(REMEMBER_ME_KEY);
        localStorage.removeItem(SAVED_ACCOUNT_KEY);
        localStorage.removeItem(SAVED_PASSWORD_KEY);
      }
    };

    const handleLogin = async () => {
      try {
        loading.value = true;
        error.value = '';
        const response = await axios.post('/api/v1/auth/login', {
          account: account.value,
          password: password.value
        });

        if (response.data && response.data.data && response.data.data.access_token) {
          const { access_token, refresh_token } = response.data.data;
          
          // 保存或清除凭据
          handleCredentials();
          
          await authStore.login(access_token, refresh_token);
          
          // 获取redirect参数，如果有则跳转回原页面，否则跳转到首页
          const redirect = router.currentRoute.value.query.redirect as string;
          router.push(redirect || '/');
        } else {
          error.value = '登录失败，请重试';
        }
      } catch (err: any) {
        console.error("Login failed:", err);
        if (err.response && err.response.data && err.response.data.message) {
            error.value = err.response.data.message;
        } else {
            error.value = '邮箱或密码错误';
        }
      } finally {
        loading.value = false;
      }
    };

    onMounted(() => {
      loadSavedCredentials();
    });

    return {
      account,
      password,
      rememberMe,
      loading,
      error,
      handleLogin
    };
  }
});
