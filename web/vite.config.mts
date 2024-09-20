import {defineConfig} from 'vite'
import react from '@vitejs/plugin-react-swc'
import {createHash} from "node:crypto";

// https://vitejs.dev/config/
// Note that Vite only performs transpilation on .ts files and does NOT perform type checking.
// It assumes type checking is taken care of by your IDE and build process. If it is going to
// be annoying we can add vite-plugin-checker for ts and eslint checking
export default defineConfig({
    // NOTE: this changes url in `.html` from static `/` to the relative `./` path
    //  to be able to work it with proxies
    base: "",
    plugins: [react()],
    build: {
        outDir: "build",
        sourcemap: true,
    },
    css: {
        modules: {
            localsConvention: "camelCase",
            generateScopedName: (name: string) => {
                // disable scoped names for codemirror classes
                if (name.indexOf("cm") !== -1) return name
                const hash = createHash("shake256", {outputLength: 3}).update(name).digest("hex")
                return `css-${name}-${hash}`
            },
        },
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
});
