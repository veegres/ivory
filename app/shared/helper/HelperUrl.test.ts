import {describe, expect, it} from "@rstest/core"

import {buildRouteUrl, isHomePath, matchRouteTail} from "./HelperUrl"

describe("matchRouteTail", () => {
    it("returns the decoded tail under the base route", () => {
        expect(matchRouteTail("/ivory/user/a%2Bb", "user", "http://localhost/ivory/")).toEqual({tail: "a+b"})
    })

    it("returns the tail when the page is opened with one", () => {
        expect(matchRouteTail("/user/abc.def", "/user", "http://localhost/")).toEqual({tail: "abc.def"})
    })

    it("returns an empty tail for the route itself", () => {
        expect(matchRouteTail("/user", "/user", "http://localhost/")).toEqual({tail: ""})
        expect(matchRouteTail("/user/", "/user", "http://localhost/")).toEqual({tail: ""})
        expect(matchRouteTail("/ivory/user", "/user/", "http://localhost/ivory/")).toEqual({tail: ""})
    })

    it("understands a base path behind a reverse proxy", () => {
        expect(matchRouteTail("/ivory/user/abc", "/user", "http://localhost/ivory/")).toEqual({tail: "abc"})
        expect(matchRouteTail("/ivory/user", "/user", "http://localhost/ivory/")).toEqual({tail: ""})
    })

    it("returns nothing outside the base route", () => {
        expect(matchRouteTail("/", "/user", "http://localhost/")).toBeUndefined()
        expect(matchRouteTail("/users", "/user", "http://localhost/")).toBeUndefined()
        expect(matchRouteTail("/user/abc", "user", "http://localhost/ivory/")).toBeUndefined()
    })
})

describe("buildRouteUrl", () => {
    it("builds an absolute route under the base uri", () => {
        expect(buildRouteUrl("/user", "abc.def", "http://localhost/")).toBe("http://localhost/user/abc.def")
    })

    it("keeps the base path behind a reverse proxy", () => {
        expect(buildRouteUrl("/user", "abc.def", "http://localhost/ivory/")).toBe("http://localhost/ivory/user/abc.def")
    })

    it("round trips with the route matcher", () => {
        const url = new URL(buildRouteUrl("/user", "abc.def", "http://localhost/ivory/"))
        expect(matchRouteTail(url.pathname, "/user", "http://localhost/ivory/")).toEqual({tail: "abc.def"})
    })
})

describe("isHomePath", () => {
    it("matches the base path with and without a trailing slash", () => {
        expect(isHomePath("/ivory", "http://localhost/ivory/")).toBe(true)
        expect(isHomePath("/ivory/", "http://localhost/ivory/")).toBe(true)
    })
})
