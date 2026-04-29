// COMMON (WEB AND SERVER)

export enum Type {
    POSTGRES,
    ETCD,
}

export interface Config {
    type: Type,
    host: string,
    port: number,
    name?: string,
    schema?: string,
}
