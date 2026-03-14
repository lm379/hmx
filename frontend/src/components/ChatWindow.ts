import { defineComponent, ref, computed, nextTick, watch, onMounted, onUnmounted } from 'vue';
import { useQAStore } from '../stores/qa';
import RelatedSources from './RelatedSources.vue';

// Default geometry (reset on every page load)
// NOTE: left/bottom are computed lazily in setup() so window.innerWidth is
// already correct when the component mounts.
const DEFAULT_BTN_BOTTOM = 28;   // px from bottom
const DEFAULT_WIN_WIDTH  = 360;  // px
const DEFAULT_WIN_HEIGHT = 540;  // px
const MIN_WIN_WIDTH  = 280;
const MIN_WIN_HEIGHT = 320;

/** 获取底部安全区高度（iPhone 刘海屏等），兜底为 0 */
function safeAreaBottom(): number {
  // CSS env() 不能直接用 JS 读取，通过临时元素间接获取
  const el = document.createElement('div');
  el.style.cssText = 'position:fixed;bottom:0;height:env(safe-area-inset-bottom,0px);pointer-events:none;visibility:hidden';
  document.body.appendChild(el);
  const h = el.getBoundingClientRect().height || 0;
  document.body.removeChild(el);
  return h;
}

/** Returns the default left position so the button sits 28 px from the right edge. */
function defaultBtnLeft() {
  const vw = window.innerWidth;
  return Math.max(4, Math.min(vw - 52 - 4, vw - 28 - 52));
}

/** Returns the default bottom position, accounting for safe area. */
function defaultBtnBottom() {
  return DEFAULT_BTN_BOTTOM + safeAreaBottom();
}

export default defineComponent({
  name: 'ChatWindow',
  components: { RelatedSources },
  setup() {
    const qaStore = useQAStore();
    const inputText = ref('');
    const messagesEl  = ref<HTMLElement | null>(null);
    const inputEl     = ref<HTMLTextAreaElement | null>(null);
    const widgetEl    = ref<HTMLElement | null>(null);
    const windowEl    = ref<HTMLElement | null>(null);
    const feedbackGiven = ref<Record<string, number>>({});

    // ── Position of the widget (button anchor) ────────────────────────
    // We track as { left, bottom } in viewport px.
    // left/bottom avoids the right-edge anchor problem: when resizing
    // from the left (w/nw/sw), only left+width need to change; the
    // right edge stays fixed naturally.
    const btnLeft   = ref(defaultBtnLeft());
    const btnBottom = ref(defaultBtnBottom());

    // ── Window size ────────────────────────────────────────────────────
    const winWidth  = ref(DEFAULT_WIN_WIDTH);
    const winHeight = ref(DEFAULT_WIN_HEIGHT);

    // ── Drag state (button) ────────────────────────────────────────────
    let dragActive  = false;
    let dragStartX  = 0;
    let dragStartY  = 0;
    let dragOrigLeft   = 0;
    let dragOrigBottom = 0;
    let dragMoved   = false;       // distinguish drag from click

    // widgetStyle: positions the 52×52 button anchor.
    // The chat window is absolutely positioned relative to this anchor via CSS.
    const widgetStyle = computed(() => ({
      position: 'fixed' as const,
      left:   `${btnLeft.value}px`,
      bottom: `${btnBottom.value}px`,
      zIndex: 9000,
      touchAction: 'none',
      userSelect: 'none' as const,
    }));

    // windowStyle: size only — position is handled by CSS (absolute, right:0, bottom: 64px)
    const windowStyle = computed(() => ({
      width:     `${winWidth.value}px`,
      height:    `${winHeight.value}px`,
      maxHeight: 'none',   // override CSS default so user resize works
    }));

    // ── Button drag ────────────────────────────────────────────────────
    function onBtnPointerDown(e: PointerEvent) {
      dragActive   = true;
      dragMoved    = false;
      dragStartX   = e.clientX;
      dragStartY   = e.clientY;
      dragOrigLeft   = btnLeft.value;
      dragOrigBottom = btnBottom.value;
      e.preventDefault();
    }

    function onBtnClick() {
      // Only toggle if we didn't actually drag
      if (!dragMoved) {
        qaStore.toggleChat();
      }
    }

    function onPointerMove(e: PointerEvent) {
      if (!dragActive) return;
      const dx = e.clientX - dragStartX;
      const dy = e.clientY - dragStartY;
      if (Math.abs(dx) > 4 || Math.abs(dy) > 4) dragMoved = true;
      if (!dragMoved) return;

      const vw = window.innerWidth;
      const vh = window.innerHeight;
      const btnSize = 52;

      let newLeft   = dragOrigLeft   + dx;        // dx>0 → move right → left increases
      let newBottom = dragOrigBottom - dy;         // dy>0 → move down  → bottom decreases

      // Clamp to viewport
      newLeft   = Math.max(4, Math.min(vw - btnSize - 4, newLeft));
      newBottom = Math.max(4, Math.min(vh - btnSize - 4, newBottom));

      btnLeft.value   = newLeft;
      btnBottom.value = newBottom;
    }

    function onPointerUp() {
      dragActive = false;
    }

    // ── Window header drag ─────────────────────────────────────────────
    let winDragActive = false;
    let winDragStartX = 0;
    let winDragStartY = 0;
    let winOrigLeft   = 0;
    let winOrigBottom = 0;

    function onWindowDragStart(e: PointerEvent) {
      // Don't start if target is an icon-btn
      if ((e.target as HTMLElement).closest('.icon-btn')) return;
      winDragActive = true;
      winDragStartX = e.clientX;
      winDragStartY = e.clientY;
      winOrigLeft   = btnLeft.value;
      winOrigBottom = btnBottom.value;
      e.preventDefault();
    }

    function onWindowPointerMove(e: PointerEvent) {
      if (!winDragActive) return;
      const dx = e.clientX - winDragStartX;
      const dy = e.clientY - winDragStartY;
      const vw = window.innerWidth;
      const vh = window.innerHeight;

      // winOrigLeft stores btnLeft at drag-start.
      // Window left edge = btnLeft + btnSize - winWidth  (window is right-aligned to btn anchor)
      const btnSize = 52;
      const winLeftAtStart = winOrigLeft + btnSize - winWidth.value;

      let newWinLeft = winLeftAtStart + dx;
      let newBottom  = winOrigBottom - dy;         // dy>0 → move down → bottom decreases

      // Clamp window to viewport
      newWinLeft = Math.max(4, Math.min(vw - winWidth.value  - 4, newWinLeft));
      newBottom  = Math.max(4, Math.min(vh - winHeight.value - 4, newBottom));

      // Convert window left edge back to btnLeft
      btnLeft.value   = newWinLeft - btnSize + winWidth.value;
      btnBottom.value = newBottom;
    }

    function onWindowPointerUp() {
      winDragActive = false;
    }

    // ── Resize ─────────────────────────────────────────────────────────
    // Widget is positioned with left/bottom.  The window grows upward
    // and rightward from the anchor (bottom-left of the widget div).
    //
    //   e (right edge)  : width += dx,  left unchanged
    //   w (left edge)   : width -= dx,  left += dx   (left edge follows pointer)
    //   s (bottom edge) : height += dy, bottom unchanged  (dy>0 → bottom edge moves down)
    //   n (top edge)    : height -= dy, bottom += dy  (dy<0 → top moves up, bottom shifts up)
    //   corners         : combination of two axes
    type ResizeDir = 'n'|'s'|'e'|'w'|'nw'|'ne'|'sw'|'se';
    let resizeActive  = false;
    let resizeDir: ResizeDir = 'se';
    let resizeStartX  = 0;
    let resizeStartY  = 0;
    let resizeOrigW   = 0;
    let resizeOrigH   = 0;
    let resizeOrigLeft   = 0;
    let resizeOrigBottom = 0;

    function onResizeStart(e: PointerEvent, dir: ResizeDir) {
      resizeActive     = true;
      resizeDir        = dir;
      resizeStartX     = e.clientX;
      resizeStartY     = e.clientY;
      resizeOrigW      = winWidth.value;
      resizeOrigH      = winHeight.value;
      resizeOrigLeft   = btnLeft.value;
      resizeOrigBottom = btnBottom.value;
      e.preventDefault();
    }

    function onResizeMove(e: PointerEvent) {
      if (!resizeActive) return;
      const dx = e.clientX - resizeStartX;
      const dy = e.clientY - resizeStartY;
      const dir = resizeDir;
      const btnSize = 52;

      let newW = resizeOrigW;
      let newH = resizeOrigH;
      // Work in window-left coordinates (not btnLeft) to keep viewport clamp simple.
      // winLeft = btnLeft + btnSize - winWidth
      let newWinLeft = resizeOrigLeft + btnSize - resizeOrigW;
      let newBottom  = resizeOrigBottom;

      // Horizontal
      if (dir.includes('e')) {
        // Right edge: widen rightward, window-left stays
        newW = Math.max(MIN_WIN_WIDTH, resizeOrigW + dx);
      }
      if (dir.includes('w')) {
        // Left edge: pointer moves left (dx<0) → widen, window-left follows pointer
        newW = Math.max(MIN_WIN_WIDTH, resizeOrigW - dx);
        const actualDw = newW - resizeOrigW;
        newWinLeft = newWinLeft - actualDw;   // left edge moves left as window widens
      }

      // Vertical (bottom-anchored)
      if (dir.includes('s')) {
        // Bottom edge moves down: bottom decreases, height increases
        const dh = dy;
        newH = Math.max(MIN_WIN_HEIGHT, resizeOrigH + dh);
        const actualDh = newH - resizeOrigH;
        newBottom = resizeOrigBottom - actualDh;
      }
      if (dir.includes('n')) {
        // Top edge moves up: height increases, bottom unchanged
        newH = Math.max(MIN_WIN_HEIGHT, resizeOrigH - dy);
      }

      // Viewport boundary protection (window-left coordinate)
      const vw = window.innerWidth;
      const vh = window.innerHeight;
      newWinLeft = Math.max(4, Math.min(vw - newW - 4, newWinLeft));
      newBottom  = Math.max(4, Math.min(vh - newH - 4, newBottom));

      winWidth.value  = newW;
      winHeight.value = newH;
      // Convert window-left back to btnLeft
      btnLeft.value   = newWinLeft + newW - btnSize;
      btnBottom.value = newBottom;
    }

    function onResizeEnd() {
      resizeActive = false;
    }

    // ── Global pointer listeners ───────────────────────────────────────
    function onGlobalPointerMove(e: PointerEvent) {
      onPointerMove(e);
      onWindowPointerMove(e);
      onResizeMove(e);
    }
    function onGlobalPointerUp() {
      onPointerUp();
      onWindowPointerUp();
      onResizeEnd();
    }

    // ── Clamp position on viewport resize ─────────────────────────────
    function onViewportResize() {
      const vw = window.innerWidth;
      const vh = window.innerHeight;
      const btnSize = 52;
      if (qaStore.isOpen) {
        // Clamp window (btn = window right-bottom corner)
        const winLeft = btnLeft.value + btnSize - winWidth.value;
        const clampedWinLeft = Math.max(4, Math.min(vw - winWidth.value - 4, winLeft));
        btnLeft.value   = Math.max(4, Math.min(vw - btnSize - 4, clampedWinLeft + winWidth.value - btnSize));
        btnBottom.value = Math.max(4, Math.min(vh - winHeight.value - 4, btnBottom.value));
      } else {
        btnLeft.value   = Math.max(4, Math.min(vw - btnSize - 4, btnLeft.value));
        btnBottom.value = Math.max(4, Math.min(vh - btnSize - 4, btnBottom.value));
      }
    }

    onMounted(() => {
      window.addEventListener('pointermove', onGlobalPointerMove);
      window.addEventListener('pointerup',   onGlobalPointerUp);
      window.addEventListener('resize',      onViewportResize);
    });
    onUnmounted(() => {
      window.removeEventListener('pointermove', onGlobalPointerMove);
      window.removeEventListener('pointerup',   onGlobalPointerUp);
      window.removeEventListener('resize',      onViewportResize);
    });

    // ── Chat logic (unchanged) ─────────────────────────────────────────
    const suggestedQuestions = [
      '黄梅戏有哪些著名剧目？',
      '黄梅戏的起源和发展历史',
      '黄梅戏有哪些知名艺术家？',
      '黄梅戏的唱腔特点是什么？',
    ];

    async function handleSend(text?: string) {
      const query = (text ?? inputText.value).trim();
      if (!query) return;
      inputText.value = '';
      if (inputEl.value) inputEl.value.style.height = 'auto';
      await qaStore.sendMessage(query);
    }

    function autoResize(e: Event) {
      const el = e.target as HTMLTextAreaElement;
      el.style.height = 'auto';
      el.style.height = Math.min(el.scrollHeight, 120) + 'px';
    }

    async function giveFeedback(msgId: string, rating: number) {
      if (feedbackGiven.value[msgId]) return;
      feedbackGiven.value[msgId] = rating;
      await qaStore.feedback(msgId, rating);
    }

    watch(
      () => qaStore.messages.length,
      async () => {
        await nextTick();
        if (messagesEl.value) {
          messagesEl.value.scrollTop = messagesEl.value.scrollHeight;
        }
      },
    );

    // Mobile: clamp window size to viewport on open
    watch(
      () => qaStore.isOpen,
      (open) => {
        if (!open) return;
        const vw = window.innerWidth;
        const vh = window.innerHeight;
        if (vw <= 480) {
          winWidth.value  = vw - 32;
          winHeight.value = Math.round(vh * 0.7);
          btnLeft.value   = 16;
          btnBottom.value = 16 + safeAreaBottom();
        }
      },
    );

    return {
      qaStore,
      inputText,
      messagesEl,
      inputEl,
      widgetEl,
      windowEl,
      feedbackGiven,
      suggestedQuestions,
      widgetStyle,
      windowStyle,
      handleSend,
      autoResize,
      giveFeedback,
      onBtnPointerDown,
      onBtnClick,
      onWindowDragStart,
      onResizeStart,
    };
  },
});
