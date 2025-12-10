import {createTheme} from "@mui/material"
import {AxiosError} from "axios"
import dayjs from "dayjs"
import utc from "dayjs/plugin/utc"
import {describe, expect, it} from "vitest"

import {Sidecar} from "../../src/api/instance/type"
import {
    DateTimeFormatter,
    getDomain,
    getDomains,
    getErrorMessage,
    getSidecar,
    getSidecars,
    isSidecarEqual,
    randomUnicodeAnimal,
    shortUuid,
    SizeFormatter,
    SxPropsFormatter
} from "../../src/app/utils"

// Extend dayjs with UTC plugin for DateTimeFormatter tests
dayjs.extend(utc)

describe("shortUuid", () => {
  it("should return the first 8 characters of a UUID", () => {
    const uuid = "12345678-90ab-cdef-1234-567890abcdef"
    expect(shortUuid(uuid)).toBe("12345678")
  })
})

describe("getDomain", () => {
  it("should return the domain string from a Sidecar object", () => {
    const sidecar: Sidecar = {host: "localhost", port: 8008}
    expect(getDomain(sidecar)).toBe("localhost:8008")
  })

})

describe("getDomains", () => {
    it("should return an array of domain strings from an array of Sidecar objects", () => {
        const sidecars: Sidecar[] = [
            {host: "localhost", port: 8008},
            {host: "127.0.0.1", port: 8008},
        ]
        expect(getDomains(sidecars)).toEqual(["localhost:8008", "127.0.0.1:8008"])
    })
})

describe("getSidecar", () => {
    it("should return a Sidecar object from a domain string", () => {
        const domain = "localhost:8008"
        expect(getSidecar(domain)).toEqual({host: "localhost", port: 8008})
    })

    it("should return a Sidecar object with default port if port is not in domain string", () => {
        const domain = "localhost"
        expect(getSidecar(domain)).toEqual({host: "localhost", port: 8008})
    })
})

describe("getSidecars", () => {
    it("should return an array of Sidecar objects from an array of domain strings", () => {
        const domains = ["localhost:8008", "127.0.0.1"]
        expect(getSidecars(domains)).toEqual([
            {host: "localhost", port: 8008},
            {host: "127.0.0.1", port: 8008},
        ])
    })
})

describe("isSidecarEqual", () => {
    it("should return true if sidecars are equal", () => {
        const sidecar1: Sidecar = {host: "localhost", port: 8008}
        const sidecar2: Sidecar = {host: "localhost", port: 8008}
        expect(isSidecarEqual(sidecar1, sidecar2)).toBe(true)
    })

    it("should return false if sidecars are not equal", () => {
        const sidecar1: Sidecar = {host: "localhost", port: 8008}
        const sidecar2: Sidecar = {host: "localhost", port: 8009}
        expect(isSidecarEqual(sidecar1, sidecar2)).toBe(false)
    })
})

describe("SxPropsFormatter", () => {
    describe("merge", () => {
        it("should merge two sx props objects", () => {
            const sx1 = {color: "red"}
            const sx2 = {backgroundColor: "blue"}
            const result = SxPropsFormatter.merge(sx1, sx2)
            expect(result).toEqual([sx1, sx2])
        })

        it("should merge sx props arrays", () => {
            const sx1 = [{color: "red"}, {fontSize: 14}]
            const sx2 = [{backgroundColor: "blue"}]
            const result = SxPropsFormatter.merge(sx1, sx2)
            expect(result).toEqual([{color: "red"}, {fontSize: 14}, {backgroundColor: "blue"}])
        })

        it("should handle undefined sx props", () => {
            const sx1 = {color: "red"}
            const result = SxPropsFormatter.merge(sx1, undefined)
            expect(result).toEqual([sx1, undefined])
        })
    })

    describe("style", () => {
        it("should have paper style", () => {
            expect(SxPropsFormatter.style.paper).toEqual({
                backgroundImage: "linear-gradient(rgba(255, 255, 255, 0.09), rgba(255, 255, 255, 0.09))"
            })
        })

        it("should generate bgImageError style", () => {
            const theme = createTheme({palette: {mode: "dark"}})
            const bgImageError = SxPropsFormatter.style.bgImageError
            if (typeof bgImageError === "function") {
                const result = bgImageError(theme) as {backgroundImage?: string}
                expect(result.backgroundImage).toBeDefined()
                expect(result.backgroundImage).toContain(theme.palette.error.dark)
            }
        })

        it("should generate bgImageSelected style", () => {
            const theme = createTheme({palette: {mode: "dark"}})
            const bgImageSelected = SxPropsFormatter.style.bgImageSelected
            if (typeof bgImageSelected === "function") {
                const result = bgImageSelected(theme) as {backgroundImage?: string}
                expect(result.backgroundImage).toBeDefined()
                expect(result.backgroundImage).toContain(theme.palette.action.hover)
            }
        })
    })
})

describe("DateTimeFormatter", () => {
    describe("utc", () => {
        it("should convert UTC time to local time with timezone", () => {
            const utcTime = "2024-01-15 10:30:00"
            const result = DateTimeFormatter.utc(utcTime)
            // Result should be in format "YYYY-MM-DD HH:mm Z"
            expect(result).toMatch(/\d{4}-\d{2}-\d{2} \d{2}:\d{2} [+-]\d{2}:\d{2}/)
        })
    })
})

describe("SizeFormatter", () => {
    describe("pretty", () => {
        it("should format bytes to human readable format", () => {
            const result1 = SizeFormatter.pretty(1024)
            const result2 = SizeFormatter.pretty(1048576)
            const result3 = SizeFormatter.pretty(500)
            // Just check that it returns a string with the expected unit
            expect(result1).toContain("K")
            expect(result2).toContain("M")
            expect(result3).toContain("B")
        })

        it("should handle zero bytes", () => {
            const result = SizeFormatter.pretty(0)
            expect(result).toContain("B")
        })
    })
})

describe("randomUnicodeAnimal", () => {
    it("should return a unicode animal string", () => {
        const animal = randomUnicodeAnimal()
        expect(typeof animal).toBe("string")
        expect(animal.length).toBeGreaterThan(0)
    })
})

describe("getErrorMessage", () => {
    it("should extract message from axios error with error field", () => {
        const error = new AxiosError("Request failed")
        error.response = {
            data: {error: "Test error message"},
            status: 400,
            statusText: "Bad Request",
            headers: {},
            config: {} as any,
        }
        expect(getErrorMessage(error)).toBe("Test error message")
    })

    it("should use error message if no response", () => {
        const error = new AxiosError("Network error")
        expect(getErrorMessage(error)).toBe("Network error")
    })

    it("should use status text if response has no data", () => {
        const error = new AxiosError("Request failed")
        error.response = {
            data: null,
            status: 500,
            statusText: "Internal Server Error",
            headers: {},
            config: {} as any,
        }
        expect(getErrorMessage(error)).toBe("500 Internal Server Error")
    })

    it("should convert string errors to string", () => {
        const error = "Simple error string"
        expect(getErrorMessage(error)).toBe("Simple error string")
    })

    it("should return unknown for unhandled error types", () => {
        const error = {someField: "value"}
        expect(getErrorMessage(error)).toBe("unknown")
    })
})
