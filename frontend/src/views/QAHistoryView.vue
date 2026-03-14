<template>
  <div class="qa-history-page">
    <div class="page-header">
      <h2>问答历史</h2>
      <p class="page-sub">你的黄梅戏知识问答记录</p>
    </div>

    <div v-if="qaStore.history.length === 0 && !loading" class="empty-state">
      <svg viewBox="0 0 24 24" width="48" height="48" fill="none" stroke="currentColor" stroke-width="1.2">
        <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"/>
      </svg>
      <p>还没有问答记录</p>
      <button class="start-btn" @click="startChat">开始提问</button>
    </div>

    <div v-else class="history-list">
      <div
        v-for="item in qaStore.history"
        :key="item.qa_id"
        class="history-card"
      >
        <div class="history-q">
          <span class="tag user-tag">问</span>
          <span class="text">{{ item.query }}</span>
        </div>
        <div class="history-a">
          <span class="tag ai-tag">答</span>
          <span class="text">{{ truncate(item.answer, 200) }}</span>
        </div>
        <div class="history-meta">
          <span>{{ formatDate(item.created_at) }}</span>
          <span v-if="item.sources && item.sources.length" class="source-count">
            {{ item.sources.length }} 个来源
          </span>
        </div>
      </div>
    </div>

    <div v-if="qaStore.historyTotal > 20" class="pagination">
      <button
        :disabled="qaStore.historyPage <= 1"
        @click="loadPage(qaStore.historyPage - 1)"
      >上一页</button>
      <span>第 {{ qaStore.historyPage }} 页 / 共 {{ Math.ceil(qaStore.historyTotal / 20) }} 页</span>
      <button
        :disabled="qaStore.historyPage >= Math.ceil(qaStore.historyTotal / 20)"
        @click="loadPage(qaStore.historyPage + 1)"
      >下一页</button>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent, ref, onMounted } from 'vue';
import { useQAStore } from '../stores/qa';

export default defineComponent({
  name: 'QAHistoryView',
  setup() {
    const qaStore = useQAStore();
    const loading = ref(true);

    onMounted(async () => {
      await qaStore.loadHistory(1);
      loading.value = false;
    });

    async function loadPage(page: number) {
      loading.value = true;
      await qaStore.loadHistory(page);
      loading.value = false;
    }

    function startChat() {
      qaStore.openChat();
    }

    function truncate(text: string, max: number): string {
      if (!text) return '';
      return text.length > max ? text.slice(0, max) + '…' : text;
    }

    function formatDate(iso: string): string {
      const d = new Date(iso);
      return d.toLocaleString('zh-CN', { hour12: false });
    }

    return { qaStore, loading, loadPage, startChat, truncate, formatDate };
  }
});
</script>

<style scoped>
.qa-history-page {
  max-width: 800px;
  margin: 0 auto;
  padding: 32px 20px;
}

.page-header {
  margin-bottom: 28px;
}
.page-header h2 {
  color: #e0e0e0;
  font-size: 22px;
  margin: 0 0 4px;
}
.page-sub {
  color: rgba(255,255,255,0.4);
  font-size: 13px;
  margin: 0;
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
  padding: 64px 0;
  color: rgba(255,255,255,0.3);
}
.empty-state p { margin: 0; font-size: 14px; }
.start-btn {
  background: #8b0000;
  color: #fff;
  border: none;
  border-radius: 8px;
  padding: 8px 20px;
  font-size: 14px;
  cursor: pointer;
  margin-top: 8px;
}
.start-btn:hover { background: #c0392b; }

.history-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.history-card {
  background: rgba(255,255,255,0.05);
  border: 1px solid rgba(255,255,255,0.08);
  border-radius: 12px;
  padding: 16px;
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.history-q, .history-a {
  display: flex;
  gap: 10px;
  align-items: flex-start;
}

.tag {
  flex-shrink: 0;
  width: 22px;
  height: 22px;
  border-radius: 6px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 11px;
  font-weight: 700;
}
.user-tag { background: #8b0000; color: #fff; }
.ai-tag { background: rgba(255,255,255,0.12); color: rgba(255,255,255,0.7); }

.text {
  flex: 1;
  color: #d0d0d0;
  font-size: 13px;
  line-height: 1.6;
  word-break: break-word;
}
.history-q .text { font-weight: 500; }

.history-meta {
  display: flex;
  gap: 12px;
  font-size: 11px;
  color: rgba(255,255,255,0.3);
  padding-top: 4px;
  border-top: 1px solid rgba(255,255,255,0.05);
}
.source-count { color: rgba(255,255,255,0.4); }

.pagination {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 16px;
  margin-top: 24px;
  color: rgba(255,255,255,0.5);
  font-size: 13px;
}
.pagination button {
  background: rgba(255,255,255,0.07);
  border: 1px solid rgba(255,255,255,0.1);
  color: rgba(255,255,255,0.7);
  border-radius: 8px;
  padding: 6px 14px;
  cursor: pointer;
  transition: background 0.15s;
}
.pagination button:hover:not(:disabled) { background: rgba(255,255,255,0.14); }
.pagination button:disabled { opacity: 0.3; cursor: not-allowed; }
</style>
