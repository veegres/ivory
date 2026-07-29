import {createHash} from "node:crypto"

import {defineConfig} from "@rsbuild/core"
import {pluginReact} from "@rsbuild/plugin-react"

export default defineConfig({
    plugins: [pluginReact()],
    source: {
        entry: {index: "./core/index.tsx"},
    },
    html: {
        template: "./index.html",
    },
    server: {
        open: true,
        publicDir: {name: "shared/assets"},
        proxy: {
            "/api": {
                target: "http://localhost:8080",
                changeOrigin: true,
                secure: false,
            },
        },
    },
    output: {
        distPath: {root: "build"},
        // NOTE: this makes asset urls relative so the build works behind the `IVORY_URL_PATH` proxy
        assetPrefix: "auto",
        sourceMap: {js: "source-map"},
    },
    tools: {
        cssLoader: {
            modules: {
                exportLocalsConvention: "camelCase",
                getLocalIdent: (_context, _localIdentName, localName) => {
                    // disable scoped names for codemirror classes
                    if (localName.indexOf("cm") !== -1) return localName
                    const hash = createHash("shake256", {outputLength: 3}).update(localName).digest("hex")
                    return `css-${localName}-${hash}`
                },
            },
        },
    },
})
