// COMMON (WEB AND SERVER)

export interface Vault {
    username: string,
    secret: string,
    type: VaultType,
    metadata?: string,
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

export interface VaultTabs {
    [key: number]: { label: string, type: VaultType }
}