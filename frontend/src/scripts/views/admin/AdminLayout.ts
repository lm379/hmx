import { defineComponent, computed } from 'vue';
import { useAuthStore } from '../../../stores/auth';
import { useRouter } from 'vue-router';

export default defineComponent({
  name: 'AdminLayout',
  setup() {
    const authStore = useAuthStore();
    const router = useRouter();
    const user = computed(() => authStore.user);

    const logout = async () => {
      await authStore.logout();
      router.push('/login');
    };

    return {
      user,
      logout
    };
  }
});
