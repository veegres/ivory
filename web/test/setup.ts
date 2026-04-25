// This file is used to set up the testing environment for vitest.
// It is a good place to add polyfills, mocks, and other setup code.

import "@testing-library/jest-dom"

import {cleanup} from "@testing-library/react"
import {afterEach} from "vitest"

import {setupLocalStorageMock} from "./test-helpers"

// Initialize global mocks
setupLocalStorageMock()

// Mock matchMedia
Object.defineProperty(window, "matchMedia", {
    writable: true,
    value: vi.fn().mockImplementation(query => ({
        matches: false,
        media: query,
        onchange: null,
        addListener: vi.fn(), // deprecated
        removeListener: vi.fn(), // deprecated
        addEventListener: vi.fn(),
        removeEventListener: vi.fn(),
        dispatchEvent: vi.fn(),
    })),
})

// Cleanup after each test
afterEach(() => {
    cleanup()
})
