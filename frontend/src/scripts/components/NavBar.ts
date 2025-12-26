import { defineComponent } from 'vue';
import { useRouter } from 'vue-router';

export default defineComponent({
  name: 'NavBar',
  setup() {
    const router = useRouter();

    const goHome = () => {
      router.push('/');
    };

    return {
      goHome
    };
  }
});
