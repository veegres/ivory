import {render, screen} from "@testing-library/react"
import {useState} from "react"
import {beforeEach, describe, expect, it, vi} from "vitest"

import {AppProvider, Mode, useSettings} from "../../src/provider/AppProvider"

// Mock the localStorage hook
vi.mock("../../src/hook/LocalStorage", () => ({
    useLocalStorageState: (_key: string, initialValue: any) => {
        return useState(initialValue)
    },
}))

// Mock QueryClient from tanstack
vi.mock("@tanstack/react-query", async (importOriginal) => {
    const actual = await importOriginal<typeof import("@tanstack/react-query")>()
    const mockQueryClient = {
        setDefaultOptions: vi.fn(),
        clear: vi.fn(),
        mount: vi.fn(),
        unmount: vi.fn(),
        isFetching: vi.fn(() => 0),
        isMutating: vi.fn(() => 0),
        getQueryData: vi.fn(),
        setQueryData: vi.fn(),
        getQueriesData: vi.fn(),
        setQueriesData: vi.fn(),
        invalidateQueries: vi.fn(),
        refetchQueries: vi.fn(),
        cancelQueries: vi.fn(),
        removeQueries: vi.fn(),
        resetQueries: vi.fn(),
        getQueryCache: vi.fn(),
        getMutationCache: vi.fn(),
    }
    return {
        ...actual,
        QueryClient: vi.fn(function(this: any) {
            return mockQueryClient
        }),
    }
})

// Mock useMediaQuery
vi.mock("@mui/material", async (importOriginal) => {
    const actual = await importOriginal<typeof import("@mui/material")>()
    return {
        ...actual,
        useMediaQuery: vi.fn(() => false),
    }
})

// Test component that uses the settings
function TestComponent() {
    const {state, theme, setTheme, toggleRefetchOnWindowsRefocus} = useSettings()

    return (
        <div>
            <div data-testid={"current-mode"}>{state.mode}</div>
            <div data-testid={"current-theme"}>{theme}</div>
            <div data-testid={"refetch-on-focus"}>{state.refetchOnWindowsFocus.toString()}</div>
            <button onClick={() => setTheme(Mode.DARK)}>Set Dark</button>
            <button onClick={() => setTheme(Mode.LIGHT)}>Set Light</button>
            <button onClick={() => setTheme(Mode.SYSTEM)}>Set System</button>
            <button onClick={toggleRefetchOnWindowsRefocus}>Toggle Refetch</button>
        </div>
    )
}

describe("AppProvider", () => {
    beforeEach(() => {
        vi.clearAllMocks()
    })

    it("should render children correctly", () => {
        render(
            <AppProvider>
                <div>Test Child</div>
            </AppProvider>
        )

        expect(screen.getByText("Test Child")).toBeInTheDocument()
    })

    it("should have correct initial state", () => {
        render(
            <AppProvider>
                <TestComponent />
            </AppProvider>
        )

        expect(screen.getByTestId("current-mode")).toHaveTextContent("system")
        expect(screen.getByTestId("refetch-on-focus")).toHaveTextContent("false")
    })


    it("should use system preference when mode is SYSTEM and prefers light", async () => {
        const {useMediaQuery} = await import("@mui/material")
        vi.mocked(useMediaQuery).mockReturnValue(false) // System prefers light

        render(
            <AppProvider>
                <TestComponent />
            </AppProvider>
        )

        // Default mode is system
        expect(screen.getByTestId("current-mode")).toHaveTextContent("system")
        expect(screen.getByTestId("current-theme")).toHaveTextContent("light")
    })

    it("should use system preference when mode is SYSTEM and prefers dark", async () => {
        const {useMediaQuery} = await import("@mui/material")
        vi.mocked(useMediaQuery).mockReturnValue(true) // System prefers dark

        render(
            <AppProvider>
                <TestComponent />
            </AppProvider>
        )

        // Default mode is system
        expect(screen.getByTestId("current-mode")).toHaveTextContent("system")
        expect(screen.getByTestId("current-theme")).toHaveTextContent("dark")
    })

})
