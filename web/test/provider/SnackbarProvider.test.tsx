import {render, screen, waitFor} from "@testing-library/react"
import {describe, expect, it} from "vitest"

import {SnackbarProvide, useSnackbar} from "../../src/provider/SnackbarProvider"

// Test component that uses the snackbar
function TestComponent() {
    const snackbar = useSnackbar()

    return (
        <div>
            <button onClick={() => snackbar("Success message", "success")}>
                Show Success
            </button>
            <button onClick={() => snackbar("Error message", "error")}>
                Show Error
            </button>
            <button onClick={() => snackbar("Warning message", "warning")}>
                Show Warning
            </button>
            <button onClick={() => snackbar("Info message", "info")}>
                Show Info
            </button>
        </div>
    )
}

describe("SnackbarProvider", () => {

    it("should render children correctly", () => {
        render(
            <SnackbarProvide>
                <div>Test Child</div>
            </SnackbarProvide>
        )

        expect(screen.getByText("Test Child")).toBeInTheDocument()
    })

    it("should show success message when snackbar is called", async () => {
        render(
            <SnackbarProvide>
                <TestComponent />
            </SnackbarProvide>
        )

        const button = screen.getByText("Show Success")
        button.click()

        await waitFor(
            () => {
                expect(screen.getByText("Success message")).toBeInTheDocument()
            },
            {timeout: 1000}
        )
    })

    it("should show error message with correct severity", async () => {
        render(
            <SnackbarProvide>
                <TestComponent />
            </SnackbarProvide>
        )

        const button = screen.getByText("Show Error")
        button.click()

        await waitFor(
            () => {
                expect(screen.getByText("Error message")).toBeInTheDocument()
            },
            {timeout: 1000}
        )

        // Check that the alert has the error severity
        const alert = screen.getByRole("alert")
        expect(alert).toHaveClass("MuiAlert-filledError")
    })

    it("should show warning message with correct severity", async () => {
        render(
            <SnackbarProvide>
                <TestComponent />
            </SnackbarProvide>
        )

        const button = screen.getByText("Show Warning")
        button.click()

        await waitFor(
            () => {
                expect(screen.getByText("Warning message")).toBeInTheDocument()
            },
            {timeout: 1000}
        )

        const alert = screen.getByRole("alert")
        expect(alert).toHaveClass("MuiAlert-filledWarning")
    })

    it("should show info message with correct severity", async () => {
        render(
            <SnackbarProvide>
                <TestComponent />
            </SnackbarProvide>
        )

        const button = screen.getByText("Show Info")
        button.click()

        await waitFor(
            () => {
                expect(screen.getByText("Info message")).toBeInTheDocument()
            },
            {timeout: 1000}
        )

        const alert = screen.getByRole("alert")
        expect(alert).toHaveClass("MuiAlert-filledInfo")
    })
})
