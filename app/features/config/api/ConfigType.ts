import {UserSetupRequest} from "../../user/api/UserType"

export interface NewAppConfig {
    secret?: string,
    appConfig: AppConfig,
    user?: UserSetupRequest,
}

export interface AppConfig {
    company: string,
    auth: AuthConfig,
}

export interface AuthConfig {
    basic?: BasicConfig,
    ldap?: LdapConfig,
    oidc?: OidcConfig,
}

export type AuthConfigObject = object

// BasicConfig carries no credentials: basic auth authenticates against the
// Ivory users, so all the config states is that it is on.
export type BasicConfig = AuthConfigObject

export interface LdapConfig extends AuthConfigObject {
    url: string,
    bindDN: string,
    bindPass: string,
    baseDN: string,
    filter: string,
}

export interface OidcConfig extends AuthConfigObject {
    issuerUrl: string,
    clientId: string,
    clientSecret: string,
    redirectUrl: string,
}
