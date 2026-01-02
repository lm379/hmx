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

    const handleLogout = async () => {
      await authStore.logout();
      router.push('/login');
    };

    const goHome = () => {
      router.push('/');
    };

    const goLogin = () => {
      // 保存当前路径，登录后跳转回来
      const currentPath = router.currentRoute.value.fullPath;
      if (currentPath !== '/login' && currentPath !== '/register') {
        router.push({ path: '/login', query: { redirect: currentPath } });
      } else {
        router.push('/login');
      }
    };

    const goProfile = () => {
      router.push('/profile');
    };

    const goCollection = () => {
      router.push({ path: '/profile', query: { tab: 'favorites' } });
    };

    const goHistory = () => {
      router.push({ path: '/profile', query: { tab: 'history' } });
    };

    const goAdmin = () => {
      router.push('/admin');
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
      goCollection,
      goHistory,
      goAdmin,
      handleLogout
    };
  }
});
