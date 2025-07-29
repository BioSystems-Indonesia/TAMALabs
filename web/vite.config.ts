import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react-swc'


// https://vite.dev/config/
export default defineConfig({
  build: {
    target: 'esnext',
    sourcemap: false, // Disable sourcemaps in production
    minify: 'esbuild', // Enable esbuild for minification
    terserOptions: {
      compress: {
        drop_console: true, // Remove console logs for production
      },
    },
  },
  server: {
    host: '0.0.0.0',
  },
  plugins: [react()],
})
