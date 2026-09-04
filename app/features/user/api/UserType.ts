// COMMON (WEB AND SERVER)

// UserAuthType is a way of signing in a user was registered for. It is a
// different question from AuthType, which names how the session in front of you
// was made - this one says what a person is allowed to use at all.
export enum UserAuthType {
    BASIC = "basic",
    LDAP = "ldap",
    OIDC = "oidc",
}

// UserRegistrationStatus is where a user stands with the one thing Ivory holds
// for them, their password. The server states it so every screen reads the same
// answer.
export enum UserRegistrationStatus {
    ACTIVE = "active",
    PENDING = "pending",
    EXPIRED = "expired",
    MISSING = "missing",
}

export interface UserRegistrationState {
    status: UserRegistrationStatus,
    expiresAt?: string,
}

export interface User {
    username: string,
    authTypes: UserAuthType[],
    superuser: boolean,
    registration?: UserRegistrationState,
}

// UserRegistration is the only shape carrying the token: a registration is
// shown once, when it is issued, and is never stored to be shown again.
export interface UserRegistration {
    token: string,
    username: string,
    expiresAt: string,
}

export interface UserRegistered {
    user: User,
    registration?: UserRegistration,
}

// UserCreateRequest carries no password: a user is created by name and by the
// ways they may sign in, and sets their own password on the page their
// registration link opens.
export interface UserCreateRequest {
    username: string,
    authTypes: UserAuthType[],
    superuser: boolean,
}

export interface UserUpdateRequest {
    authTypes: UserAuthType[],
}

// UserSetupRequest is the first user Ivory gets, registered a superuser. Setup
// is the one place a password is typed on somebody else's behalf, because there
// is nobody yet who could be sent a registration.
export interface UserSetupRequest {
    username: string,
    password: string,
    authTypes: UserAuthType[],
}

export interface UserRegistrationPayload {
    username: string,
    expiresAt: string,
}

export interface UserRegistrationPasswordRequest {
    token: string,
    password: string,
}

export interface UserPasswordUpdateRequest {
    previousPassword: string,
    newPassword: string,
}
