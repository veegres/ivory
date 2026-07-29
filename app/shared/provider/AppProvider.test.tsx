import {beforeEach, describe, expect, it, rs} from "@rstest/core"
import * as actualReactQuery from "@tanstack/react-query" with {rstest: "importActual"}
import {render, screen, waitFor} from "@testing-library/react"
import {useState} from "react"

import {AppProvider, Mode, useSettings} from "./AppProvider"

// Mock the localStorage hook
rs.mock("../hook/LocalStorage", () => ({
    useLocalStorageState: (_key: string, initialValue: any) => {
        return useState(initialValue)
    },
}))

// Mock QueryClient from tanstack
rs.mock("@tanstack/react-query", () => {
    const mockQueryClient = {
        setDefaultOptions: rs.fn(),
        clear: rs.fn(),
        mount: rs.fn(),
        unmount: rs.fn(),
        isFetching: rs.fn(() => 0),
        isMutating: rs.fn(() => 0),
        getQueryData: rs.fn(),
        setQueryData: rs.fn(),
        getQueriesData: rs.fn(),
        setQueriesData: rs.fn(),
        invalidateQueries: rs.fn(),
        refetchQueries: rs.fn(),
        cancelQueries: rs.fn(),
        removeQueries: rs.fn(),
        resetQueries: rs.fn(),
        getQueryCache: rs.fn(),
        getMutationCache: rs.fn(),
    }
    return {
        ...actualReactQuery,
        QueryClient: rs.fn(function(this: any) {
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
        rs.clearAllMocks()
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
        rs.mocked(window.matchMedia).mockImplementation(query => ({
            matches: query.includes("dark"),
            media: query,
            onchange: null,
            addListener: rs.fn(),
            removeListener: rs.fn(),
            addEventListener: rs.fn(),
            removeEventListener: rs.fn(),
            dispatchEvent: rs.fn(),
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
        rs.mocked(window.matchMedia).mockImplementation(query => ({
            matches: false,
            media: query,
            onchange: null,
            addListener: rs.fn(),
            removeListener: rs.fn(),
            addEventListener: rs.fn(),
            removeEventListener: rs.fn(),
            dispatchEvent: rs.fn(),
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
        rs.mocked(window.matchMedia).mockImplementation(query => ({
            matches: true,
            media: query,
            onchange: null,
            addListener: rs.fn(),
            removeListener: rs.fn(),
            addEventListener: rs.fn(),
            removeEventListener: rs.fn(),
            dispatchEvent: rs.fn(),
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
