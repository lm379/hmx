import { defineComponent, ref, onMounted } from 'vue';
import axios from 'axios';

export default defineComponent({
  name: 'DashboardView',
  setup() {
    const stats = ref({
      total_users: 0,
      total_operas: 0,
      total_artists: 0,
      total_likes: 0,
      total_plays: 0,
      total_favorites: 0,
      total_shares: 0,
    });

    const fetchStats = async () => {
      try {
        const res = await axios.get('/api/v1/admin/stats');
        if (res.data && res.data.data) {
          stats.value = res.data.data;
        }
      } catch (err) {
        console.error(err);
      }
    };

    onMounted(() => {
      fetchStats();
    });

    return {
      stats
    };
  }
});
