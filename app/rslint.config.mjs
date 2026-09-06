import {defineConfig, globalIgnores, globals, js, reactHooksPlugin, reactPlugin, ts} from "@rslint/core"
import simpleImportSortPlugin from "eslint-plugin-simple-import-sort"

export default defineConfig([
    globalIgnores(["node_modules/**", "build/**", "coverage/**"]),
    js.configs.recommended,
    ts.configs.recommended,
    reactPlugin.configs.recommended,
    reactHooksPlugin.configs.recommended,
    {
        files: ["**/*.{js,jsx,ts,tsx}"],
        languageOptions: {
            globals: {...globals.browser, ...globals.node},
        },
        settings: {
            react: {version: "detect"},
        },
        rules: {
            "quotes": ["error", "double"],
            "semi": ["error", "never"],
            "react/react-in-jsx-scope": "off",
            "react/no-unescaped-entities": "off",
            "@typescript-eslint/no-explicit-any": "off",
            "no-restricted-imports": "error",
            "react/jsx-curly-brace-presence": ["error", {"props": "always"}],
            "object-curly-spacing": ["error", "never"],
            "react/function-component-definition": ["error", {"namedComponents": "function-declaration"}],
        },
    },
    {
        plugins: {
            "simple-import-sort": simpleImportSortPlugin,
        },
        rules: {
            "simple-import-sort/imports": "error",
            "simple-import-sort/exports": "error",
        },
    },
])
