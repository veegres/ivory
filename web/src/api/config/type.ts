export interface NewAppConfig {
    secret?: string,
    appConfig: AppConfig,
}

export interface AppConfig {
    company: string,
    auth: AuthConfig,
}

export interface AuthConfig {
    superusers: string[],
    basic?: BasicConfig,
    ldap?: LdapConfig,
    oidc?: OidcConfig,
}

export type AuthConfigObject = object

export interface BasicConfig extends AuthConfigObject {
    username: string,
    password: string,
}

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
