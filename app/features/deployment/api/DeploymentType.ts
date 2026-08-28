import {KeeperPlugin, PlatformPlugin} from "../../node/api/NodeType"

// CreationType distinguishes the templates Ivory ships from the ones you own,
// using the same vocabulary as the query builder's own templates.
export enum CreationType {
    Manual = "manual",
    System = "system",
}

// TemplateCommand is one node's deployment. It has no identity beyond its
// position: which node it lands on is chosen at deploy time.
export interface TemplateCommand {
    command: string,
    postScript?: string,
}

// Template is a saved deployment: an ordered list of commands, one per node,
// and nothing about where they run.
export interface Template {
    id: string,
    name: string,
    description?: string,
    keeper: KeeperPlugin,
    platform: PlatformPlugin,
    commands: TemplateCommand[],
    creation: CreationType,
    createdAt: number,
}

export interface TemplateRequest {
    name: string,
    description?: string,
    keeper: KeeperPlugin,
    platform: PlatformPlugin,
    commands: TemplateCommand[],
}

export interface TemplateListRequest {
    keeper?: KeeperPlugin,
    platform?: PlatformPlugin,
}
