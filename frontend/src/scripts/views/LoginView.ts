import { defineComponent, ref } from 'vue';
import { useRouter } from 'vue-router';
import axios from 'axios';

export default defineComponent({
  name: 'LoginView',
  setup() {
    const router = useRouter();
    const account = ref('');
    const password = ref('');
    const loading = ref(false);
    const error = ref('');

    const handleLogin = async () => {
      try {
        loading.value = true;
        error.value = '';
        const response = await axios.post('/api/v1/auth/login', {
          account: account.value,
          password: password.value
        });

        if (response.data && response.data.data && response.data.data.token) {
          const token = response.data.data.token;
          localStorage.setItem('token', token);
          // Set axios default authorization header
          axios.defaults.headers.common['Authorization'] = `Bearer ${token}`;
          
          router.push('/');
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

    return {
      account,
      password,
      loading,
      error,
      handleLogin
    };
  }
});
