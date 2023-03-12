export interface SecretStatus {
    key: boolean,
    ref: boolean,
}

export interface SecretSetRequest {
    key: string,
    ref: string,
}

export interface SecretUpdateRequest {
    previousKey: string,
    newKey: string,
}
