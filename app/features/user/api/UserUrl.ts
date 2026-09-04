// UserWebPath mirrors the server's `user.WebPath`: the sub-path a registration
// link points at, with the token as the segment after it.
export const UserWebPath = "user"

export interface UserRegistrationRoute {
    token: string,
}

// matchUserRegistrationRoute answers whether a path is the registration page
// and, if it is, what token it carries - an empty token meaning the page was
// opened without one, which is what sends the visitor home.
export function matchUserRegistrationRoute(pathname: string, baseUri: string): UserRegistrationRoute | undefined {
    const base = new URL(baseUri).pathname.replace(/\/$/, "")
    const prefix = `${base}/${UserWebPath}`
    if (pathname !== prefix && !pathname.startsWith(`${prefix}/`)) return undefined
    const token = pathname.slice(prefix.length + 1).replace(/\/$/, "")
    return {token: decodeURIComponent(token)}
}

export function getUserRegistrationRoute(): UserRegistrationRoute | undefined {
    return matchUserRegistrationRoute(window.location.pathname, document.baseURI)
}

export function buildUserRegistrationUrl(token: string, baseUri: string = document.baseURI): string {
    return new URL(`${UserWebPath}/${encodeURIComponent(token)}`, baseUri).href
}

export function redirectToHome() {
    window.location.replace(document.baseURI)
}
