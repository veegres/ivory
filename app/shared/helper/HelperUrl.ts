export interface RouteTail {
    tail: string,
}

export function matchRouteTail(pathname: string, routePath: string, baseUri: string): RouteTail | undefined {
    const prefix = buildRoutePrefix(routePath, baseUri)
    if (pathname !== prefix && !pathname.startsWith(`${prefix}/`)) return undefined
    const tail = pathname === prefix ? "" : pathname.slice(prefix.length + 1).replace(/\/$/, "")
    return {tail: decodeURIComponent(tail)}
}

export function getRouteTail(routePath: string): RouteTail | undefined {
    return matchRouteTail(window.location.pathname, routePath, document.baseURI)
}

export function buildRouteUrl(routePath: string, tail?: string, baseUri: string = document.baseURI): string {
    const encodedTail = tail ? `/${encodeURIComponent(tail)}` : ""
    return new URL(`${normalizeRoutePath(routePath)}${encodedTail}`, baseUri).href
}

export function redirectToHome() {
    window.location.replace(document.baseURI)
}

export function isHomePath(pathname: string = window.location.pathname, baseUri: string = document.baseURI): boolean {
    const base = getBasePath(baseUri)
    return pathname === base || pathname === `${base}/`
}

function buildRoutePrefix(routePath: string, baseUri: string): string {
    const base = getBasePath(baseUri)
    return `${base}/${normalizeRoutePath(routePath)}`
}

function getBasePath(baseUri: string): string {
    return new URL(baseUri).pathname.replace(/\/$/, "")
}

function normalizeRoutePath(routePath: string): string {
    return routePath.replace(/^\/+|\/+$/g, "")
}
