import {describe, expect, it} from "@rstest/core"

import {buildUserRegistrationUrl, matchUserRegistrationRoute} from "./UserUrl"

describe("matchUserRegistrationRoute", () => {
    it("returns the token when the page is opened with one", () => {
        expect(matchUserRegistrationRoute("/user/abc.def", "http://localhost/")).toEqual({token: "abc.def"})
    })

    it("returns an empty token when the page is opened without one", () => {
        expect(matchUserRegistrationRoute("/user", "http://localhost/")).toEqual({token: ""})
        expect(matchUserRegistrationRoute("/user/", "http://localhost/")).toEqual({token: ""})
    })

    it("understands a base path behind a reverse proxy", () => {
        expect(matchUserRegistrationRoute("/ivory/user/abc", "http://localhost/ivory/")).toEqual({token: "abc"})
        expect(matchUserRegistrationRoute("/ivory/user", "http://localhost/ivory/")).toEqual({token: ""})
    })

    it("decodes an escaped token", () => {
        expect(matchUserRegistrationRoute("/user/a%2Bb", "http://localhost/")).toEqual({token: "a+b"})
    })

    it("returns nothing for any other page", () => {
        expect(matchUserRegistrationRoute("/", "http://localhost/")).toBeUndefined()
        expect(matchUserRegistrationRoute("/users", "http://localhost/")).toBeUndefined()
        expect(matchUserRegistrationRoute("/user/abc", "http://localhost/ivory/")).toBeUndefined()
    })
})

describe("buildUserRegistrationUrl", () => {
    it("builds an absolute link the invited person can open", () => {
        expect(buildUserRegistrationUrl("abc.def", "http://localhost/")).toBe("http://localhost/user/abc.def")
    })

    it("keeps the base path behind a reverse proxy", () => {
        expect(buildUserRegistrationUrl("abc.def", "http://localhost/ivory/")).toBe("http://localhost/ivory/user/abc.def")
    })

    it("round trips with the route matcher", () => {
        const url = new URL(buildUserRegistrationUrl("abc.def", "http://localhost/ivory/"))
        expect(matchUserRegistrationRoute(url.pathname, "http://localhost/ivory/")).toEqual({token: "abc.def"})
    })
})
