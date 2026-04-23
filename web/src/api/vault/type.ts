// COMMON (WEB AND SERVER)

export interface Vault {
    username: string,
    secret: string,
    type: VaultType,
}

export enum VaultType {
    DATABASE_PASSWORD,
    KEEPER_PASSWORD,
    SSH_PASSWORD,
    SSH_KEY,
}

export interface VaultMap {
    [uuid: string]: Vault,
}

// SPECIFIC (WEB)
