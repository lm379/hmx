import { defineStore } from 'pinia';
import { ref } from 'vue';
import { askQuestion, getQAHistory, submitFeedback } from '../api/qa';
import type { AskResponse, QAHistoryItem, Source } from '../api/qa';

export interface ChatMessage {
  id: string;
  role: 'user' | 'assistant';
  content: string;
  sources?: Source[];
  loading?: boolean;
  error?: boolean;
}

const SESSION_ID_KEY = 'qa_session_id';

function getOrCreateSessionId(): string {
  let sid = localStorage.getItem(SESSION_ID_KEY);
  if (!sid) {
    sid = crypto.randomUUID();
    localStorage.setItem(SESSION_ID_KEY, sid);
  }
  return sid;
}

export const useQAStore = defineStore('qa', () => {
  const isOpen = ref(false);
  const messages = ref<ChatMessage[]>([]);
  const isLoading = ref(false);
  const sessionId = ref<string>(getOrCreateSessionId());

  // QA history (separate from chat messages)
  const history = ref<QAHistoryItem[]>([]);
  const historyTotal = ref(0);
  const historyPage = ref(1);

  function openChat() {
    isOpen.value = true;
  }

  function closeChat() {
    isOpen.value = false;
  }

  function toggleChat() {
    isOpen.value = !isOpen.value;
  }

  async function sendMessage(query: string, operaId?: number) {
    if (!query.trim() || isLoading.value) return;

    // Add user message
    messages.value.push({
      id: crypto.randomUUID(),
      role: 'user',
      content: query
    });

    // Add loading placeholder
    const loadingId = crypto.randomUUID();
    messages.value.push({
      id: loadingId,
      role: 'assistant',
      content: '',
      loading: true
    });

    isLoading.value = true;

    try {
      const res = await askQuestion(query, sessionId.value, operaId);
      const data: AskResponse = res.data;

      // Replace loading placeholder
      const idx = messages.value.findIndex(m => m.id === loadingId);
      if (idx !== -1) {
        messages.value[idx] = {
          id: data.qa_id,
          role: 'assistant',
          content: data.answer,
          sources: data.sources,
          loading: false
        };
      }
    } catch (e) {
      const idx = messages.value.findIndex(m => m.id === loadingId);
      if (idx !== -1) {
        messages.value[idx] = {
          id: loadingId,
          role: 'assistant',
          content: '抱歉，问答服务暂时不可用，请稍后再试。',
          loading: false,
          error: true
        };
      }
    } finally {
      isLoading.value = false;
    }
  }

  async function loadHistory(page = 1) {
    try {
      const res = await getQAHistory(page, 20, sessionId.value);
      history.value = res.data.list || [];
      historyTotal.value = res.data.pagination.total;
      historyPage.value = page;
    } catch (e) {
      // silently fail
    }
  }

  async function feedback(qaId: string, rating: number, comments?: string) {
    try {
      await submitFeedback(qaId, rating, comments);
    } catch (e) {
      // silently fail
    }
  }

  function clearMessages() {
    messages.value = [];
  }

  return {
    isOpen,
    messages,
    isLoading,
    sessionId,
    history,
    historyTotal,
    historyPage,
    openChat,
    closeChat,
    toggleChat,
    sendMessage,
    loadHistory,
    feedback,
    clearMessages
  };
});
