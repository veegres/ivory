// This file is used to set up the testing environment for vitest.
// It is a good place to add polyfills, mocks, and other setup code.

import "@testing-library/jest-dom"

import {cleanup} from "@testing-library/react"
import {afterEach} from "vitest"

// Cleanup after each test
afterEach(() => {
    cleanup()
})
