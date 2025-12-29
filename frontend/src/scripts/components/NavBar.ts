import { defineComponent, onMounted } from 'vue';
import { storeToRefs } from 'pinia';
import { useRouter } from 'vue-router';
import { useAuthStore } from '../../stores/auth';

export default defineComponent({
  name: 'NavBar',
  setup() {
    const router = useRouter();
    const authStore = useAuthStore();
    const { isLoggedIn, user } = storeToRefs(authStore);

    const handleLogout = () => {
      authStore.logout();
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
      authStore.checkLoginStatus();
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
