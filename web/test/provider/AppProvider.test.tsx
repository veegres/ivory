import {render, screen, waitFor} from "@testing-library/react"
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

    it("should have correct initial state", async () => {
        vi.mocked(window.matchMedia).mockImplementation(query => ({
            matches: query.includes("dark"),
            media: query,
            onchange: null,
            addListener: vi.fn(),
            removeListener: vi.fn(),
            addEventListener: vi.fn(),
            removeEventListener: vi.fn(),
            dispatchEvent: vi.fn(),
        }))

        render(
            <AppProvider>
                <TestComponent />
            </AppProvider>
        )

        expect(screen.getByTestId("current-mode")).toHaveTextContent("system")
        await waitFor(() => {
            expect(screen.getByTestId("current-theme")).toHaveTextContent("dark")
        })
        expect(screen.getByTestId("refetch-on-focus")).toHaveTextContent("false")
    })


    it("should use system preference when mode is SYSTEM and prefers light", async () => {
        vi.mocked(window.matchMedia).mockImplementation(query => ({
            matches: false,
            media: query,
            onchange: null,
            addListener: vi.fn(),
            removeListener: vi.fn(),
            addEventListener: vi.fn(),
            removeEventListener: vi.fn(),
            dispatchEvent: vi.fn(),
        }))

        render(
            <AppProvider>
                <TestComponent />
            </AppProvider>
        )

        // Default mode is system
        expect(screen.getByTestId("current-mode")).toHaveTextContent("system")
        await waitFor(() => {
            expect(screen.getByTestId("current-theme")).toHaveTextContent("light")
        })
    })

    it("should use system preference when mode is SYSTEM and prefers dark", async () => {
        vi.mocked(window.matchMedia).mockImplementation(query => ({
            matches: true,
            media: query,
            onchange: null,
            addListener: vi.fn(),
            removeListener: vi.fn(),
            addEventListener: vi.fn(),
            removeEventListener: vi.fn(),
            dispatchEvent: vi.fn(),
        }))

        render(
            <AppProvider>
                <TestComponent />
            </AppProvider>
        )

        // Default mode is system
        expect(screen.getByTestId("current-mode")).toHaveTextContent("system")
        await waitFor(() => {
            expect(screen.getByTestId("current-theme")).toHaveTextContent("dark")
        })
    })

})
