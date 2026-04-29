// COMMON (WEB AND SERVER)

export enum Plugin {
    POSTGRES = "postgres",
    ETCD = "etcd",
}

export interface Config {
    plugin: Plugin,
    host: string,
    port: number,
    name?: string,
    schema?: string,
}
