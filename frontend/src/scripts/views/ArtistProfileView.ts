import { defineComponent, ref, onMounted, type Ref } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import axios from 'axios';
import type { ArtistDetail } from '../../types';
import VideoCard from '../../components/VideoCard.vue';

export default defineComponent({
  name: 'ArtistProfileView',
  components: {
      VideoCard
  },
  setup() {
    const route = useRoute();
    const router = useRouter();
    const artist: Ref<ArtistDetail | null> = ref(null);
    const loading = ref(true);

    const fetchArtist = async () => {
      try {
        loading.value = true;
        const id = route.params.id;
        const response = await axios.get(`/api/v1/artists/${id}`);
        if (response.data && response.data.data) {
           artist.value = response.data.data;
        }
      } catch (error: any) {
        console.error("Failed to fetch artist:", error);
        if (error.response && error.response.status === 404) {
             router.replace('/404');
        }
      } finally {
        loading.value = false;
      }
    };

    onMounted(() => {
        fetchArtist();
    });

    return {
      artist,
      loading
    };
  }
});
