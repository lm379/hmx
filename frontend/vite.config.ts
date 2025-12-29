import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import svgLoader from 'vite-svg-loader'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [vue(), svgLoader()],
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
    emptyOutDir: true
  }
})
