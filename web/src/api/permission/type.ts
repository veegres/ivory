// COMMON (WEB AND SERVER)

export enum PermissionStatus {
    NOT_PERMITTED = 0,
    PENDING = 1,
    GRANTED = 2,
}

export interface UserPermissions {
    username: string,
    permissions: { [permission: string]: PermissionStatus },
}

// SPECIFIC (WEB)