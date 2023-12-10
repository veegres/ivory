import {defineConfig} from 'vite'
import react from '@vitejs/plugin-react-swc'

// https://vitejs.dev/config/
// Note that Vite only performs transpilation on .ts files and does NOT perform type checking.
// It assumes type checking is taken care of by your IDE and build process. If it is going to
// be annoying we can add vite-plugin-checker for ts and eslint checking
export default defineConfig({
    plugins: [react()],
    build: {
        outDir: 'build'
    },
    server: {
        open: true,
        proxy: {
            "/api": {
                target: "http://localhost:8080",
                changeOrigin: true,
                secure: false,
            },
        },
    },
})
