// AuthToken is a session token Ivory minted on request, for a client that
// cannot carry its cookie. It is shown once, when it is issued.
export interface AuthToken {
    token: string,
    expiresAt: string,
}
