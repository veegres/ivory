// COMMON (WEB AND SERVER)

export interface Password {
    username: string,
    password: string,
    type: PasswordType,
}

export enum PasswordType {
    POSTGRES,
    PATRONI,
}

export interface PasswordMap {
    [uuid: string]: Password,
}

// SPECIFIC (WEB)
