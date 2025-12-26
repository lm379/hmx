<template>
  <div class="artist-profile-view" v-if="artist">
    <!-- Header -->
    <div class="artist-header">
       <div class="artist-avatar-large">
          <img :src="artist.avatar || `https://ui-avatars.com/api/?name=${artist.name}&background=random&size=200`" alt="Avatar" />
       </div>
       <div class="artist-info">
          <h1 class="artist-name">{{ artist.name }}</h1>
          <p class="artist-bio">{{ artist.bio || '这位艺术家很低调，暂无简介。' }}</p>
       </div>
    </div>

    <!-- Works -->
    <div class="artist-works">
      <h2 class="section-title">代表作品</h2>
      <div class="video-grid" v-if="artist.operas && artist.operas.length > 0">
        <VideoCard 
            v-for="opera in artist.operas" 
            :key="opera.opera_id" 
            :opera="opera"
            @click="$router.push(`/video/${opera.opera_id}`)"
        />
      </div>
      <div v-else class="empty-state">
        暂无收录作品
      </div>
    </div>
  </div>
  <div v-else-if="loading" class="loading-state">
      加载中...
  </div>
  <div v-else class="error-state">
      未找到艺术家信息
  </div>
</template>

<script lang="ts">
import { defineComponent, ref, onMounted, type Ref } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import axios from 'axios';
import type { Artist } from '../types';
import VideoCard from '../components/VideoCard.vue';

export default defineComponent({
  name: 'ArtistProfileView',
  components: {
      VideoCard
  },
  setup() {
    const route = useRoute();
    const router = useRouter();
    const artist: Ref<Artist | null> = ref(null);
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
</script>

<style scoped>
.artist-profile-view {
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px 24px;
}

.artist-header {
  display: flex;
  gap: 40px;
  margin-bottom: 40px;
  background: white;
  padding: 30px;
  border-radius: 8px;
  border: 1px solid #e3e5e7;
}

.artist-avatar-large {
  width: 160px;
  height: 160px;
  border-radius: 50%;
  overflow: hidden;
  flex-shrink: 0;
  border: 4px solid #f1f2f3;
}

.artist-avatar-large img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.artist-info {
  flex: 1;
}

.artist-name {
  font-size: 28px;
  font-weight: 600;
  color: #18191c;
  margin: 10px 0 20px;
}

.artist-bio {
  font-size: 15px;
  color: #61666d;
  line-height: 1.8;
  white-space: pre-wrap;
}

.section-title {
  font-size: 20px;
  font-weight: 600;
  margin-bottom: 20px;
  border-left: 4px solid #A40000;
  padding-left: 12px;
}

.video-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 20px;
}

@media (max-width: 1000px) {
  .video-grid {
    grid-template-columns: repeat(3, 1fr);
  }
}

@media (max-width: 768px) {
  .artist-header {
      flex-direction: column;
      align-items: center;
      text-align: center;
      gap: 20px;
  }
  
  .video-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}

.loading-state, .error-state, .empty-state {
  text-align: center;
  padding: 40px;
  color: #9499a0;
  font-size: 16px;
}
</style>
