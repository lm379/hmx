<template>
  <!-- Teleport to body: avoids any stacking-context issues on mobile -->
  <Teleport to="body">
    <div class="chat-widget" :style="widgetStyle" ref="widgetEl">
      <!-- Floating toggle button (draggable handle) — hidden while window is open -->
      <button v-show="!qaStore.isOpen" class="chat-toggle-btn" @pointerdown.stop="onBtnPointerDown"
        @click.stop="onBtnClick" title="知识问答">
        <svg viewBox="0 0 24 24" width="24" height="24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z" />
        </svg>
      </button>

      <!-- Chat window -->
      <transition name="chat-slide">
        <div v-if="qaStore.isOpen" class="chat-window" :style="windowStyle" ref="windowEl">
          <!-- Header (drag handle for window) -->
          <div class="chat-header" @pointerdown.stop="onWindowDragStart">
            <div class="chat-header-left">
              <svg viewBox="0 0 24 24" width="18" height="18" fill="none" stroke="currentColor" stroke-width="2">
                <circle cx="12" cy="12" r="10" />
                <line x1="12" y1="8" x2="12" y2="12" />
                <line x1="12" y1="16" x2="12.01" y2="16" />
              </svg>
              <span>黄梅戏知识问答</span>
            </div>
            <div class="chat-header-actions">
              <button class="icon-btn" title="清除对话" @click.stop="qaStore.clearMessages()">
                <svg viewBox="0 0 24 24" width="16" height="16" fill="none" stroke="currentColor" stroke-width="2">
                  <polyline points="3 6 5 6 21 6" />
                  <path d="M19 6l-1 14H6L5 6" />
                  <path d="M10 11v6M14 11v6" />
                </svg>
              </button>
              <button class="icon-btn" title="关闭" @click.stop="qaStore.closeChat()">
                <svg viewBox="0 0 24 24" width="16" height="16" fill="none" stroke="currentColor" stroke-width="2">
                  <line x1="18" y1="6" x2="6" y2="18" />
                  <line x1="6" y1="6" x2="18" y2="18" />
                </svg>
              </button>
            </div>
          </div>

          <!-- Messages -->
          <div class="chat-messages" ref="messagesEl">
            <div v-if="qaStore.messages.length === 0" class="chat-welcome">
              <div class="welcome-icon">
                <svg viewBox="0 0 24 24" width="32" height="32" fill="none" stroke="currentColor" stroke-width="1.5">
                  <path d="M12 2a10 10 0 1 0 10 10A10 10 0 0 0 12 2zm0 18a8 8 0 1 1 8-8 8 8 0 0 1-8 8z" />
                  <path d="M12 6v6l4 2" />
                </svg>
              </div>
              <p class="welcome-title">你好！有什么想了解黄梅戏的？</p>
              <p class="welcome-sub">可以问我关于剧目、唱腔、历史、艺术家等内容</p>
              <div class="suggested-questions">
                <button v-for="q in suggestedQuestions" :key="q" class="suggestion-btn" @click="handleSend(q)">{{ q
                  }}</button>
              </div>
            </div>

            <div v-for="msg in qaStore.messages" :key="msg.id" class="message-row" :class="msg.role">
              <div class="message-bubble">
                <div v-if="msg.loading" class="loading-dots">
                  <span></span><span></span><span></span>
                </div>
                <div v-else class="message-text">{{ msg.content }}</div>
                <RelatedSources v-if="!msg.loading && msg.role === 'assistant' && msg.sources" :sources="msg.sources" />
                <div v-if="!msg.loading && msg.role === 'assistant' && !msg.error" class="feedback-row">
                  <button class="feedback-btn" :class="{ active: feedbackGiven[msg.id] === 1 }"
                    @click="giveFeedback(msg.id, 1)" title="有帮助">
                    <svg viewBox="0 0 24 24" width="13" height="13" fill="none" stroke="currentColor" stroke-width="2">
                      <path d="M14 9V5a3 3 0 0 0-3-3l-4 9v11h11.28a2 2 0 0 0 2-1.7l1.38-9a2 2 0 0 0-2-2.3H14z" />
                      <path d="M7 22H4a2 2 0 0 1-2-2v-7a2 2 0 0 1 2-2h3" />
                    </svg>
                  </button>
                  <button class="feedback-btn" :class="{ active: feedbackGiven[msg.id] === -1 }"
                    @click="giveFeedback(msg.id, -1)" title="没有帮助">
                    <svg viewBox="0 0 24 24" width="13" height="13" fill="none" stroke="currentColor" stroke-width="2">
                      <path d="M10 15v4a3 3 0 0 0 3 3l4-9V2H5.72a2 2 0 0 0-2 1.7l-1.38 9a2 2 0 0 0 2 2.3H10z" />
                      <path d="M17 2h2.67A2.31 2.31 0 0 1 22 4v7a2.31 2.31 0 0 1-2.33 2H17" />
                    </svg>
                  </button>
                </div>
              </div>
            </div>
          </div>

          <!-- Input area -->
          <div class="chat-input-area">
            <textarea v-model="inputText" ref="inputEl" placeholder="输入你的问题…" rows="1"
              @keydown.enter.exact.prevent="handleSend()" @input="autoResize" />
            <button class="send-btn" :disabled="!inputText.trim() || qaStore.isLoading" @click="handleSend()">
              <svg viewBox="0 0 24 24" width="18" height="18" fill="none" stroke="currentColor" stroke-width="2">
                <line x1="22" y1="2" x2="11" y2="13" />
                <polygon points="22 2 15 22 11 13 2 9 22 2" />
              </svg>
            </button>
          </div>

          <!-- 8-direction resize handles -->
          <div class="rh rh-n" data-dir="n" @pointerdown.stop="onResizeStart($event, 'n')" />
          <div class="rh rh-s" data-dir="s" @pointerdown.stop="onResizeStart($event, 's')" />
          <div class="rh rh-w" data-dir="w" @pointerdown.stop="onResizeStart($event, 'w')" />
          <div class="rh rh-e" data-dir="e" @pointerdown.stop="onResizeStart($event, 'e')" />
          <div class="rh rh-nw" data-dir="nw" @pointerdown.stop="onResizeStart($event, 'nw')" />
          <div class="rh rh-ne" data-dir="ne" @pointerdown.stop="onResizeStart($event, 'ne')" />
          <div class="rh rh-sw" data-dir="sw" @pointerdown.stop="onResizeStart($event, 'sw')" />
          <div class="rh rh-se" data-dir="se" @pointerdown.stop="onResizeStart($event, 'se')" />
        </div>
      </transition>
    </div>
  </Teleport>
</template>

<script lang="ts" src="./ChatWindow.ts"></script>

<style scoped src="./ChatWindow.css"></style>
