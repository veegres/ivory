import {pluginReact} from "@rsbuild/plugin-react"
import {defineConfig} from "@rstest/core"

export default defineConfig({
    plugins: [pluginReact()],
    globals: true,
    testEnvironment: "jsdom",
    setupFiles: "shared/test/TestSetup.ts",
})
