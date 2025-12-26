import { defineComponent, ref, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import axios from 'axios';
import type { User } from '../../types';

export default defineComponent({
  name: 'NavBar',
  setup() {
    const router = useRouter();
    const isLoggedIn = ref(false);
    const user = ref<User | null>(null);

    const checkLoginStatus = async () => {
      const token = localStorage.getItem('token');
      if (token) {
        isLoggedIn.value = true;
        // Set axios default authorization header
        axios.defaults.headers.common['Authorization'] = `Bearer ${token}`;
        
        try {
          const response = await axios.get('/api/v1/users/me');
          if (response.data && response.data.data) {
            user.value = response.data.data;
          }
        } catch (error) {
          console.error("Failed to fetch user info:", error);
          handleLogout();
        }
      } else {
        isLoggedIn.value = false;
        user.value = null;
      }
    };

    const handleLogout = () => {
      localStorage.removeItem('token');
      delete axios.defaults.headers.common['Authorization'];
      isLoggedIn.value = false;
      user.value = null;
      router.push('/login');
    };

    const goHome = () => {
      router.push('/');
    };

    const goLogin = () => {
      router.push('/login');
    };

    const goProfile = () => {
      router.push('/profile');
    };

    onMounted(() => {
      checkLoginStatus();
    });

    return {
      isLoggedIn,
      user,
      goHome,
      goLogin,
      goProfile,
      handleLogout
    };
  }
});
