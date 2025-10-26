
import {describe, expect,it} from "vitest"

import {Sidecar} from "../../src/api/instance/type"
import {getDomain, getDomains, getSidecar, getSidecars, isSidecarEqual,shortUuid} from "../../src/app/utils"

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
