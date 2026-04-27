import { computed, defineComponent, nextTick, onMounted, onUnmounted, ref } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import * as pdfjsLib from 'pdfjs-dist/build/pdf.mjs';
import pdfWorkerUrl from 'pdfjs-dist/build/pdf.worker.mjs?url';
import { getEducationBookById } from '../../api/education';
import type { EducationBookDetail } from '../../types';

pdfjsLib.GlobalWorkerOptions.workerSrc = pdfWorkerUrl;

type PageStatus = 'pending' | 'loading' | 'rendered' | 'error';

interface PdfPageState {
  pageNumber: number;
  status: PageStatus;
}

export default defineComponent({
  name: 'EducationReaderView',
  setup() {
    const route = useRoute();
    const router = useRouter();
    const book = ref<EducationBookDetail | null>(null);
    const loading = ref(false);
    const pdfLoading = ref(false);
    const pdfError = ref('');
    const totalPages = ref(0);
    const pages = ref<PdfPageState[]>([]);
    const zoom = ref(1.15);
    const renderedCount = computed(() => pages.value.filter(page => page.status === 'rendered').length);
    const zoomPercent = computed(() => `${Math.round(zoom.value * 100)}%`);

    let pdfDocument: any = null;
    let observer: IntersectionObserver | null = null;
    const canvasRefs = new Map<number, HTMLCanvasElement>();
    const shellRefs = new Map<number, Element>();
    const visiblePages = new Set<number>();

    const fetchDetail = async () => {
      const id = Number(route.params.id);
      if (!id) {
        router.replace('/404');
        return;
      }

      loading.value = true;
      try {
        const res = await getEducationBookById(id);
        book.value = res.data;
        loading.value = false;
        // 页面元信息先展示，PDF 文档加载和后续按页渲染独立进行。
        loadPdf(res.data.pdf_url);
      } catch {
        book.value = null;
        loading.value = false;
      } finally {
        loading.value = false;
      }
    };

    const loadPdf = async (pdfUrl: string) => {
      pdfLoading.value = true;
      pdfError.value = '';
      totalPages.value = 0;
      pages.value = [];
      pdfDocument = null;

      try {
        await openPdfDocument(pdfUrl);
      } catch (error) {
        console.error('Failed to load PDF through original URL', error);
        pdfError.value = 'PDF 加载失败，请检查 PDF 文件是否可访问';
      } finally {
        pdfLoading.value = false;
      }
    };

    const openPdfDocument = async (pdfUrl: string) => {
        // PDF.js 会配合支持 Range 请求的对象存储按需拉取数据；前端只按可见页渲染 canvas。
        const task = pdfjsLib.getDocument({
          url: pdfUrl,
          rangeChunkSize: 65536,
          disableAutoFetch: false,
          disableStream: false
        });
        task.onProgress = (progress: { loaded: number; total: number }) => {
          if (progress.total > 0) {
            pdfError.value = '';
          }
        };
        pdfDocument = await task.promise;
        totalPages.value = pdfDocument.numPages;
        pages.value = Array.from({ length: pdfDocument.numPages }, (_, index) => ({
          pageNumber: index + 1,
          status: 'pending'
        }));
        await nextTick();
        setupObserver();
    };

    const setupObserver = () => {
      observer?.disconnect();
      observer = new IntersectionObserver(
        entries => {
          entries.forEach(entry => {
            const pageNumber = Number((entry.target as HTMLElement).dataset.pageNumber);
            if (!pageNumber) return;
            if (entry.isIntersecting) {
              visiblePages.add(pageNumber);
              renderPage(pageNumber);
            } else {
              visiblePages.delete(pageNumber);
            }
          });
        },
        {
          root: null,
          rootMargin: '500px 0px',
          threshold: 0.01
        }
      );

      shellRefs.forEach((el, pageNumber) => {
        (el as HTMLElement).dataset.pageNumber = String(pageNumber);
        observer?.observe(el);
      });
    };

    const setCanvasRef = (el: unknown, pageNumber: number) => {
      if (el instanceof HTMLCanvasElement) {
        canvasRefs.set(pageNumber, el);
      } else {
        canvasRefs.delete(pageNumber);
      }
    };

    const setPageShellRef = (el: unknown, pageNumber: number) => {
      if (el instanceof Element) {
        shellRefs.set(pageNumber, el);
        (el as HTMLElement).dataset.pageNumber = String(pageNumber);
        observer?.observe(el);
      } else {
        const current = shellRefs.get(pageNumber);
        if (current) observer?.unobserve(current);
        shellRefs.delete(pageNumber);
      }
    };

    const updatePageStatus = (pageNumber: number, status: PageStatus) => {
      const page = pages.value.find(item => item.pageNumber === pageNumber);
      if (page) page.status = status;
    };

    const renderPage = async (pageNumber: number) => {
      if (!pdfDocument) return;
      const state = pages.value.find(page => page.pageNumber === pageNumber);
      if (!state || state.status === 'loading' || state.status === 'rendered') return;

      const canvas = canvasRefs.get(pageNumber);
      if (!canvas) return;

      updatePageStatus(pageNumber, 'loading');
      try {
        const page = await pdfDocument.getPage(pageNumber);
        const viewport = page.getViewport({ scale: zoom.value });
        const outputScale = window.devicePixelRatio || 1;
        const context = canvas.getContext('2d');
        if (!context) throw new Error('Canvas context unavailable');

        canvas.width = Math.floor(viewport.width * outputScale);
        canvas.height = Math.floor(viewport.height * outputScale);
        canvas.style.width = `${Math.floor(viewport.width)}px`;
        canvas.style.height = `${Math.floor(viewport.height)}px`;

        context.setTransform(outputScale, 0, 0, outputScale, 0, 0);
        await page.render({ canvasContext: context, viewport }).promise;
        updatePageStatus(pageNumber, 'rendered');
      } catch (error) {
        console.error(`Failed to render PDF page ${pageNumber}`, error);
        updatePageStatus(pageNumber, 'error');
      }
    };

    const rerenderVisiblePages = async () => {
      pages.value.forEach(page => {
        if (visiblePages.has(page.pageNumber) || page.status === 'rendered') {
          page.status = 'pending';
          const canvas = canvasRefs.get(page.pageNumber);
          const context = canvas?.getContext('2d');
          if (canvas && context) {
            context.clearRect(0, 0, canvas.width, canvas.height);
          }
        }
      });
      await nextTick();
      visiblePages.forEach(pageNumber => renderPage(pageNumber));
    };

    const zoomIn = () => {
      zoom.value = Math.min(1.8, Number((zoom.value + 0.15).toFixed(2)));
      rerenderVisiblePages();
    };

    const zoomOut = () => {
      zoom.value = Math.max(0.7, Number((zoom.value - 0.15).toFixed(2)));
      rerenderVisiblePages();
    };

    onMounted(() => {
      fetchDetail();
    });

    onUnmounted(() => {
      observer?.disconnect();
    });

    return {
      book,
      loading,
      pdfLoading,
      pdfError,
      totalPages,
      pages,
      renderedCount,
      zoomPercent,
      setCanvasRef,
      setPageShellRef,
      zoomIn,
      zoomOut
    };
  }
});
