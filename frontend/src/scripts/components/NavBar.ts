import { defineComponent, onMounted, onUnmounted, ref } from 'vue';
import { storeToRefs } from 'pinia';
import { useRouter } from 'vue-router';
import { useAuthStore } from '../../stores/auth';

export default defineComponent({
  name: 'NavBar',
  setup() {
    const router = useRouter();
    const authStore = useAuthStore();
    const { isLoggedIn, user } = storeToRefs(authStore);
    const searchQuery = ref('');
    const mobileDropdownOpen = ref(false);

    const isMobile = () => window.innerWidth <= 767;

    const toggleMobileDropdown = () => {
      if (!isMobile()) return; // PC 端由 CSS hover 处理，不走此逻辑
      mobileDropdownOpen.value = !mobileDropdownOpen.value;
    };

    const closeMobileDropdown = () => {
      mobileDropdownOpen.value = false;
    };

    const goSearch = () => {
      router.push('/search');
    };

    const handleSearch = () => {
      if (searchQuery.value.trim()) {
        router.push({
          path: '/search',
          query: { q: searchQuery.value.trim() }
        });
      }
    };

    const handleLogout = async () => {
      closeMobileDropdown();
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
      closeMobileDropdown();
      router.push('/profile');
    };

    const goCollection = () => {
      router.push({ path: '/profile', query: { tab: 'favorites' } });
    };

    const goHistory = () => {
      router.push({ path: '/profile', query: { tab: 'history' } });
    };

    const goAdmin = () => {
      closeMobileDropdown();
      router.push('/admin');
    };

    // 点击页面其他区域关闭 dropdown（仅移动端）
    const handleOutsideClick = (e: MouseEvent) => {
      if (!isMobile()) return;
      const target = e.target as HTMLElement;
      if (!target.closest('.avatar-wrapper') && !target.closest('.mobile-dropdown-overlay')) {
        closeMobileDropdown();
      }
    };

    onMounted(() => {
      authStore.checkLoginStatus();
      document.addEventListener('click', handleOutsideClick);
    });

    onUnmounted(() => {
      document.removeEventListener('click', handleOutsideClick);
    });

    return {
      isLoggedIn,
      user,
      searchQuery,
      mobileDropdownOpen,
      toggleMobileDropdown,
      closeMobileDropdown,
      handleSearch,
      goHome,
      goSearch,
      goLogin,
      goProfile,
      goCollection,
      goHistory,
      goAdmin,
      handleLogout
    };
  }
});
