
import globals from "globals";
import tseslint from "typescript-eslint";
import pluginReact from "eslint-plugin-react";
import pluginReactHooks from "eslint-plugin-react-hooks";
import pluginReactRefresh from "eslint-plugin-react-refresh";
import pluginTanstackQuery from "@tanstack/eslint-plugin-query";
import pluginSimpleImportSort from "eslint-plugin-simple-import-sort";

export default [
  {
    ignores: ["dist", "node_modules"],
  },
  {
    files: ["**/*.{js,jsx,ts,tsx}"],
    plugins: {
      "@typescript-eslint": tseslint.plugin,
      "react": pluginReact,
      "react-hooks": pluginReactHooks,
      "react-refresh": pluginReactRefresh,
      "@tanstack/query": pluginTanstackQuery,
      "simple-import-sort": pluginSimpleImportSort,
    },
    languageOptions: {
      parser: tseslint.parser,
      parserOptions: {
        ecmaFeatures: {
          jsx: true,
        },
        ecmaVersion: "latest",
        sourceType: "module",
      },
      globals: {
        ...globals.browser,
        ...globals.node,
      },
    },
    settings: {
      react: {
        version: "detect",
      },
    },
    rules: {
      ...tseslint.configs.recommended.rules,
      ...pluginReact.configs.recommended.rules,
      // eslint-plugin-react-hooks v7's "recommended" bundles the full React Compiler rule
      // set (react-hooks/use-memo, set-state-in-effect, refs, etc.), which conflicts with
      // this codebase's handle***/render*** hook convention. Keep only the two rules that
      // were enabled pre-v7.
      "react-hooks/rules-of-hooks": "error",
      "react-hooks/exhaustive-deps": "warn",
      ...pluginTanstackQuery.configs.recommended.rules,
      "simple-import-sort/imports": "error",
      "simple-import-sort/exports": "error",
      "quotes": ["error", "double"],
      "semi": ["error", "never"],
      "react/react-in-jsx-scope": "off",
      "react/no-unescaped-entities": "off",
      "@typescript-eslint/no-explicit-any": "off",
      "no-restricted-imports": "error",
      "react/jsx-curly-brace-presence": ["error", { "props": "always" }],
      "object-curly-spacing": ["error", "never"],
      "react/function-component-definition": ["error", { "namedComponents": "function-declaration" }],
    },
  },
];
