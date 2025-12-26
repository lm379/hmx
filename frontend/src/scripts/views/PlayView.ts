import { defineComponent, ref, onMounted, onUnmounted, watch, nextTick } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import axios from 'axios';
import DPlayer from 'dplayer';
import type { Opera } from '../../types';

export default defineComponent({
  name: 'PlayView',
  setup() {
    const route = useRoute();
    const router = useRouter();
    const opera = ref<Opera | null>(null);
    const loading = ref(true);
    const isPip = ref(false);
    const playerBox = ref<HTMLElement | null>(null);
    const dplayerContainer = ref<HTMLElement | null>(null);
    const recommendations = ref<Opera[]>([]);
    let dp: DPlayer | null = null;

    const fetchOpera = async () => {
      // Destroy previous player instance if exists
      if (dp) {
          dp.destroy();
          dp = null;
      }
      
      try {
        loading.value = true;
        const id = route.params.id;
        const response = await axios.get(`/api/v1/operas/${id}`);
        if (response.data && response.data.data) {
          opera.value = response.data.data;
        }

        // Fetch random recommendations
        const recResponse = await axios.get('/api/v1/operas/');
        if (recResponse.data && recResponse.data.data && recResponse.data.data.list) {
          recommendations.value = recResponse.data.data.list.slice(0, 5);
        }
      } catch (error: any) {
        console.error("Failed to fetch opera:", error);
        if (error.response && error.response.status === 404) {
            router.replace('/404');
        }
      } finally {
        loading.value = false;
        // Try to init player here after loading is false
        await nextTick();
        initDPlayer();
      }
    };

    const initDPlayer = () => {
      if (!opera.value) {
          console.warn("initDPlayer: Opera data is missing");
          return;
      }
      
      // Double check container availability
      if (!dplayerContainer.value) {
          console.warn("initDPlayer: Container element not found. Waiting for DOM update...");
          return;
      }

      try {
          const options: any = {
            container: dplayerContainer.value,
            video: {
              url: opera.value.video_path,
              pic: opera.value.avatar,
            },
            autoplay: false,
            theme: '#b7daff',
            lang: 'zh-cn',
            screenshot: false,
            hotkey: true,
            preload: 'auto',
            volume: 0.7,
            mutex: true,
          };

          if (opera.value.srt_path) {
             options.subtitle = {
                 url: opera.value.srt_path,
                 type: 'webvtt',
                 fontSize: '25px',
                 bottom: '0%',
                 color: '#b7daff',
             };
          }

          dp = new DPlayer(options);
      } catch (e) {
          console.error("Error initializing DPlayer:", e);
      }
    };

    const handleScroll = () => {
      if (!playerBox.value) return;

      const rect = playerBox.value.getBoundingClientRect();
      if (rect.bottom < 100) {
        isPip.value = true;
      } else {
        isPip.value = false;
      }
    };

    // Watch for route changes to reload video
    watch(() => route.params.id, (newId) => {
        if (newId) {
            fetchOpera();
        }
    });

    onMounted(() => {
      fetchOpera();
      window.addEventListener('scroll', handleScroll);
    });

    onUnmounted(() => {
      window.removeEventListener('scroll', handleScroll);
      if (dp) {
        dp.destroy();
      }
    });

    const formatTime = (time: string) => {
      if (!time) return '';
      const date = new Date(time);
      const year = date.getFullYear();
      const month = String(date.getMonth() + 1).padStart(2, '0');
      const day = String(date.getDate()).padStart(2, '0');
      const hours = String(date.getHours()).padStart(2, '0');
      const minutes = String(date.getMinutes()).padStart(2, '0');
      const seconds = String(date.getSeconds()).padStart(2, '0');
      return `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`;
    };

    const formatArtists = (artists: { Name: string }[]) => {
      if (!artists || artists.length === 0) return '黄梅戏官方';
      // 超过三人则只显示前三人名字
      const names = artists.slice(0, 3).map(artist => artist.Name);
      return names.length === 3 ? names.join(', ') + '等' : names.join(', ');
    };

    return {
      opera,
      loading,
      isPip,
      playerBox,
      dplayerContainer,
      recommendations,
      formatTime,
      formatArtists
    };
  }
});
