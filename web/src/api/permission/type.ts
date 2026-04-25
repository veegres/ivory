import {Feature} from "../feature"

// COMMON (WEB AND SERVER)

export enum Status {
    NOT_PERMITTED = 0,
    PENDING = 1,
    GRANTED = 2,
}

export type PermissionMap = {
    [key in Feature]: Status;
};

export interface UserPermissions {
    username: string,
    permissions: PermissionMap,
}

// SPECIFIC (WEB)
