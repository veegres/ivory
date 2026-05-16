# Frontend development

## Build

Use `npm install` to download all libraries.

In the project directory, you can run:

- `npm start` - Runs the app in the development mode
- `npm build` - Builds the app for production to the `build` folder
- `npm test` - Launches the test runner in the interactive watch mode
- `npm run lint` - Checks code style and quality

You can use `yarn` if you want to.

## Testing Guidelines

### What to Test

Focus on **unit tests for business logic**:

- **Providers** - Test state management and context providers (e.g., `StoreProvider`, `SnackbarProvider`)
- **Custom Hooks** - Test hook logic, state updates, and side effects (e.g., `useInstanceDetection`, `useDebounce`)
- **Utility Functions** - Test pure functions and helpers (e.g., `utils.ts`)

### What NOT to Test

- **Component Tests** - UI components should not be tested at the unit level. Focus on business logic instead.
- **API Router Tests** - API endpoint tests belong in the backend, not frontend.
- **Integration Tests** - Complex integration scenarios are better suited for E2E tests.

### Test Structure

- All tests are located in the `test/` directory
- Directory structure mirrors `src/` structure (e.g., `src/hook/Debounce.ts` → `test/hook/Debounce.test.ts`)
- Use mock factories from `test/test-helpers.ts` for creating test data
- Follow existing test patterns for consistency

### Testing Tools

- **Vitest** - Test framework with fake timers and mocking
- **React Testing Library** - For testing hooks with `renderHook`
- **@testing-library/jest-dom** - For DOM assertions like `toBeInTheDocument`

## Learn More

This app is built with:

- [React](https://reactjs.org) - A JavaScript library for building user interfaces
- [React Query](https://react-query.tanstack.com) - Provide data global state for requests with fetch, cache and update data
- [MUI](https://mui.com) - The React Material Design library
- [Vitest](https://vitest.dev) - Testing framework
