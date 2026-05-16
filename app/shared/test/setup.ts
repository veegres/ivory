// This file is used to set up the testing environment for vitest.
// It is a good place to add polyfills, mocks, and other setup code.

import "@testing-library/jest-dom"

import {cleanup} from "@testing-library/react"
import {afterEach} from "vitest"

import {MutationObserverMock, ResizeObserverMock, setupLocalStorageMock, setupMatchMediaMock} from "./mocks"

// Initialize global mocks
setupLocalStorageMock()
setupMatchMediaMock()

// Setup Observer mocks globally
globalThis.ResizeObserver = ResizeObserverMock as any
globalThis.MutationObserver = MutationObserverMock as any

// Cleanup after each test
afterEach(() => {
    cleanup()
})
