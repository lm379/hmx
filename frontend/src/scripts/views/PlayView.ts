import { defineComponent, ref, onMounted, onUnmounted, watch, nextTick, computed } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { useAuthStore } from '../../stores/auth';
import { ElMessage, ElMessageBox } from 'element-plus';
import axios from 'axios';
import DPlayer from 'dplayer';
import { getSimilarOperas } from '../../api/recommendation';
import type { OperaDetail, OperaListItem, Comment } from '../../types';
import EmojiIcon from '../../assets/emoji.svg?component';
import PlayCountIcon from '../../assets/play_count.svg?component';
import LikeIcon from '../../assets/like.svg?component';
import FavIcon from '../../assets/fav.svg?component';
import ShareIcon from '../../assets/share.svg?component';

export default defineComponent({
  name: 'PlayView',
  components: {
    EmojiIcon,
    PlayCountIcon,
    LikeIcon,
    FavIcon,
    ShareIcon
  },
  setup() {
    const route = useRoute();
    const router = useRouter();
    const authStore = useAuthStore();

    const opera = ref<OperaDetail | null>(null);
    const loading = ref(true);
    const isPip = ref(false);
    const playerBox = ref<HTMLElement | null>(null);
    const dplayerContainer = ref<HTMLElement | null>(null);
    const recommendations = ref<OperaListItem[]>([]);
    const newCommentText = ref('');
    const submittingComment = ref(false);
    const showEmojiPicker = ref(false);

    const comments = ref<Comment[]>([]);
    const emojiList = [
      "😀", "😃", "😄", "😁", "😆", "😅", "🤣", "😂", "🙂", "🙃", "😉", "😊", "😇", "🥰", "😍", "🤩", "😘", "😗", "☺️", "😚", "😙", "🥲", "😋", "😛", "😜", "🤪", "😝", "🤑", "🤗", "🤭", "🤫", "🤔", "🤐", "🤨", "😐", "😑", "😶", "😏", "😒", "🙄", "😬", "🤥", "😌", "😔", "😪", "🤤", "😴", "😷", "🤒", "🤕", "🤢", "🤮", "🤧", "🥵", "🥶", "🥴", "😵", "🤯", "🤠", "🥳", "😎", "🤓", "🧐", "😕", "😟", "🙁", "☹️", "😮", "😯", "😲", "😳", "🥺", "😦", "😧", "😨", "😰", "😥", "😢", "😭", "😱", "😖", "😣", "😞", "😓", "😩", "😫", "🥱", "😤", "😡", "😠", "🤬", "😈", "👿", "💀", "☠️", "💩", "🤡", "👹", "👺", "👻", "👽", "👾", "🤖", "😺", "😸", "😹", "😻", "😼", "😽", "🙀", "😿", "😾", "🙈", "🙉", "🙊", "👍", "👎", "👊", "✊", "🤛", "🤜", "🤞", "✌️", "🤟", "🤘", "👌", "🤏", "👈", "👉", "👆", "👇", "☝️", "✋", "🤚", "🖐", "🖖", "👋", "🤙", "💪", "🖕", "✍️", "🙏", "🦶", "🦵", "👂", "🦻", "👃", "🧠", "🦷", "🦴", "👀", "👁", "👅", "👄", "💋"
    ];

    let dp: DPlayer | null = null;

    const isLoggedIn = computed(() => authStore.isLoggedIn);
    const currentUser = computed(() => authStore.user);

    const toggleEmojiPicker = () => {
      showEmojiPicker.value = !showEmojiPicker.value;
    };

    const addEmoji = (emoji: string) => {
      newCommentText.value += emoji;
      showEmojiPicker.value = false; // Close after picking
    };

    const toggleReplyEmojiPicker = (comment: Comment) => {
      comment.showReplyEmojiPicker = !comment.showReplyEmojiPicker;
    };

    const addEmojiToReply = (comment: Comment, emoji: string) => {
      if (!comment.replyText) {
        comment.replyText = '';
      }
      comment.replyText += emoji;
      comment.showReplyEmojiPicker = false; // Close after picking
    };

    // Close emoji picker when clicking outside (simple implementation using event listener on window)
    const closeEmojiPicker = (e: MouseEvent) => {
      const target = e.target as HTMLElement;
      if (!target.closest('.emoji-trigger') && !target.closest('.emoji-picker')) {
        showEmojiPicker.value = false;
        // 关闭所有回复框的表情选择器
        comments.value.forEach(comment => {
          comment.showReplyEmojiPicker = false;
          if (comment.replies) {
            comment.replies.forEach(reply => {
              reply.showReplyEmojiPicker = false;
            });
          }
        });
      }
    };

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
          // Record watch history
          axios.post(`/api/v1/operas/${id}/history`).catch(err => console.error("Failed to record history:", err));
        }

        // Fetch similar operas recommendations
        try {
          const similarResponse = await getSimilarOperas(Number(id), { limit: 5 });
          if (similarResponse.data && similarResponse.data.operas) {
            recommendations.value = similarResponse.data.operas;
          }
        } catch (error) {
          console.error("Failed to fetch similar operas, falling back to random:", error);
        }
        // Fallback to random recommendations
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

    const fetchComments = async () => {
      const id = route.params.id;
      try {
        const res = await axios.get(`/api/v1/operas/${id}/comments`);
        if (res.data && res.data.data) {
          // 组织评论为树形结构
          const allComments = res.data.data as Comment[];
          const commentMap = new Map<number, Comment>();
          const rootComments: Comment[] = [];

          // 初始化所有评论
          allComments.forEach(comment => {
            comment.replies = [];
            comment.showReplyInput = false;
            comment.replyText = '';
            comment.showReplyEmojiPicker = false;
            commentMap.set(comment.comment_id, comment);
          });

          // 构建树形结构
          allComments.forEach(comment => {
            if (comment.parent_comment_id) {
              const parent = commentMap.get(comment.parent_comment_id);
              if (parent) {
                parent.replies!.push(comment);
              }
            } else {
              rootComments.push(comment);
            }
          });

          comments.value = rootComments;
        }
      } catch (e) {
        console.error("Failed to fetch comments", e);
      }
    };

    const postComment = async () => {
      if (!newCommentText.value.trim()) return;
      const id = route.params.id;
      submittingComment.value = true;
      try {
        await axios.post(`/api/v1/operas/${id}/comments`, {
          comment_text: newCommentText.value
        });
        newCommentText.value = '';
        await fetchComments();
      } catch (e: any) {
        console.error("Failed to post comment", e);
        ElMessage.error(e.response?.data?.error || "发表评论失败");
      } finally {
        submittingComment.value = false;
      }
    };

    const deleteComment = async (commentId: number) => {
      try {
        await ElMessageBox.confirm(
          '确定要删除这条评论吗？',
          '提示',
          {
            confirmButtonText: '确定',
            cancelButtonText: '取消',
            type: 'warning'
          }
        );
        await axios.delete(`/api/v1/comments/${commentId}`);
        ElMessage.success("删除成功");
        await fetchComments();
      } catch (e: any) {
        if (e !== 'cancel') {
          console.error("Failed to delete comment", e);
          ElMessage.error(e.response?.data?.error || "删除失败");
        }
      }
    };

    const canDelete = (comment: Comment) => {
      if (!isLoggedIn.value || !currentUser.value) return false;
      // Check if owner or admin
      return currentUser.value.user_id === comment.user_id || currentUser.value.role === 'Administrator';
    };

    const onToggleCommentLike = async (comment: Comment) => {
      if (!isLoggedIn.value) {
        ElMessage.warning("请先登录");
        return;
      }
      try {
        const res = await axios.post(`/api/v1/comments/${comment.comment_id}/like`);
        const data = res.data?.data;
        if (data) {
          comment.liked = data.liked;
          comment.like_count = data.like_count;
        }
      } catch (e) {
        console.error("Failed to toggle comment like", e);
      }
    };

    const toggleReplyInput = (comment: Comment) => {
      if (!isLoggedIn.value) {
        ElMessage.warning("请先登录");
        return;
      }
      comment.showReplyInput = !comment.showReplyInput;
      if (!comment.showReplyInput) {
        comment.replyText = '';
      }
    };

    const postReply = async (parentComment: Comment) => {
      if (!parentComment.replyText || !parentComment.replyText.trim()) return;
      const id = route.params.id;
      try {
        await axios.post(`/api/v1/operas/${id}/comments`, {
          comment_text: parentComment.replyText,
          parent_comment_id: parentComment.comment_id
        });
        parentComment.replyText = '';
        parentComment.showReplyInput = false;
        await fetchComments();
      } catch (e: any) {
        console.error("Failed to post reply", e);
        ElMessage.error(e.response?.data?.error || "回复失败");
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
          theme: '#af000e',
          lang: 'zh-cn',
          screenshot: false,
          hotkey: true,
          preload: 'auto',
          volume: 0.5,
          mutex: false,
        };

        if (opera.value.srt_path) {
          options.subtitle = {
            url: opera.value.srt_path,
            type: 'webvtt',
            fontSize: '25px',
            bottom: '5%',
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
      window.addEventListener('click', closeEmojiPicker);
    });

    onUnmounted(() => {
      window.removeEventListener('scroll', handleScroll);
      window.removeEventListener('click', closeEmojiPicker);
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

    const formatArtists = (artists: { name: string }[]) => {
      if (!artists || artists.length === 0) return '黄梅戏官方';
      // 超过三人则只显示前三人名字
      const names = artists.slice(0, 3).map(artist => artist.name);
      return names.length === 3 ? names.join(', ') + '等' : names.join(', ');
    };

    const onToggleLike = async () => {
      if (!opera.value) return;
      const id = route.params.id;
      try {
        const res = await axios.post(`/api/v1/operas/${id}/like`);
        const data = res.data?.data || {};
        if (opera.value) {
          opera.value.like_count = data.like_count ?? opera.value.like_count ?? 0;
          opera.value.favorite_count = data.favorite_count ?? opera.value.favorite_count ?? 0;
          opera.value.share_count = data.share_count ?? opera.value.share_count ?? 0;
          opera.value.liked = data.liked ?? opera.value.liked ?? false;
          opera.value.favorited = data.favorited ?? opera.value.favorited ?? false;
        }
      } catch (e) {
        console.error('Failed to toggle like', e);
      }
    };

    const onToggleFavorite = async () => {
      if (!opera.value) return;
      const id = route.params.id;
      try {
        const res = await axios.post(`/api/v1/operas/${id}/favorite`);
        const data = res.data?.data || {};
        if (opera.value) {
          opera.value.like_count = data.like_count ?? opera.value.like_count ?? 0;
          opera.value.favorite_count = data.favorite_count ?? opera.value.favorite_count ?? 0;
          opera.value.share_count = data.share_count ?? opera.value.share_count ?? 0;
          opera.value.liked = data.liked ?? opera.value.liked ?? false;
          opera.value.favorited = data.favorited ?? opera.value.favorited ?? false;
        }
      } catch (e) {
        console.error('Failed to toggle favorite', e);
      }
    };

    const onShare = async () => {
      if (!opera.value) return;
      const id = route.params.id;
      try {
        const res = await axios.post(`/api/v1/operas/${id}/share`);
        const data = res.data?.data || {};
        if (opera.value) {
          opera.value.share_count = data.share_count ?? opera.value.share_count ?? 0;
          // 分享不改变 like/favorite 状态
        }
      } catch (e) {
        console.error('Failed to record share', e);
      }
    };

    return {
      opera,
      loading,
      isPip,
      playerBox,
      dplayerContainer,
      recommendations,
      comments,
      newCommentText,
      submittingComment,
      isLoggedIn,
      currentUser,
      showEmojiPicker,
      emojiList,
      formatTime,
      formatArtists,
      onToggleLike,
      onToggleFavorite,
      onShare,
      postComment,
      canDelete,
      toggleEmojiPicker,
      addEmoji,
      toggleReplyEmojiPicker,
      addEmojiToReply,
      onToggleCommentLike,
      toggleReplyInput,
      postReply,
      deleteComment
    };
  }
});
