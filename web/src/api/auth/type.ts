export enum AuthType {
    NONE,
    BASIC,
    LDAP,
    OIDC,
}

export interface AuthInfo {
    type: AuthType,
    authorised: boolean,
    error: string,
}