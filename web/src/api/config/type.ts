import {AuthType} from "../auth/type";

export interface ConfigInfo {
    company: string,
    configured: boolean,
    availability: Availability,
    error: string,
}

export interface AuthConfig {
    type: AuthType,
    basic?: BasicConfig,
    ldap?: LdapConfig,
    oidc?: OidcConfig,
}

export interface BasicConfig {
    username: string,
    password: string,
}

export interface LdapConfig {
    url: string,
    bindDN: string,
    bindPass: string,
    baseDN: string,
    filter: string,
}

export interface OidcConfig {
    issuerUrl: string,
    clientId: string,
    clientSecret: string,
    redirectUrl: string,
}

export interface AppConfig {
    company: string,
    availability: Availability,
    auth: AuthConfig,
}

export interface Availability {
    manualQuery: boolean,
}