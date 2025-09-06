export enum AuthType {
    NONE,
    BASIC,
}

export interface AuthConfig {
    type: AuthType,
    body: any,
}

export interface AppConfig {
    company: string,
    availability: Availability,
    auth: AuthConfig,
}

export interface Availability {
    manualQuery: boolean,
}