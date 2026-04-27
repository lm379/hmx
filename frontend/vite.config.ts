import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import svgLoader from 'vite-svg-loader'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [vue(), svgLoader()],
  optimizeDeps: {
    // PDF.js 体积大且带 worker，开发环境预构建容易卡在 /node_modules/.vite/deps/pdfjs-dist.js。
    // 让浏览器直接按 ESM 加载具体 build 文件，避免点击 PDF 时 504。
    exclude: ['pdfjs-dist']
  },
  server: {
    proxy: {
      '/api': {
        target: 'http://localhost:15379',
        changeOrigin: true,
        // rewrite: (path) => path.replace(/^\/api/, '') // backend route starts with /api/v1 so we don't need to rewrite if frontend calls /api/v1
      }
    }
  },
  build: {
    outDir: '../cmd/server/static', 
    emptyOutDir: true,
    rollupOptions: {
      output: {
        manualChunks(id) {
          if (id.includes('node_modules')) {
            return 'vendor';
          }
        },
        entryFileNames: 'js/[name]-[hash].js',
        chunkFileNames: 'js/[name]-[hash].js',
        assetFileNames: (assetInfo) => {
          if (assetInfo.name && assetInfo.name.endsWith('.css')) {
            return 'css/[name]-[hash][extname]';
          }
          return 'assets/[name]-[hash][extname]';
        }
      }
    }
  }
})
