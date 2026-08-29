import {KeeperDeploySpecResponse, KeeperPlugin, PlatformPlugin} from "../../node/api/NodeType"

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
    // NOTE: each step runs as its own execution, so nothing may assume a shell
    // - the images these run in are increasingly distroless and have no "&&"
    postScripts?: string[],
    defaults?: TemplateDefaults,
}

// TemplateDefaults is what this command fills its node card in with. A command
// and the endpoints it answers on are one fact rather than two: a single-host
// template writes a distinct peer port into each of its commands, and only that
// command knows which client port has to match it. Host and ssh port are absent
// on purpose - both describe the machine, which is what a template never knows.
export interface TemplateDefaults {
    name?: string,
    keeperPort?: number,
    dbPort?: number,
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

// DeployCredentials is what a deploy preview is allowed to show of one
// credential pair: the username as it really is, and the password only ever as
// its mask - the real one is substituted on the server.
export interface DeployCredentials {
    user?: string,
    pass?: string,
}

// DeployScreenProps is what DeploymentTemplateDialog hands the screen that
// runs a template: the template itself, the keeper's requirements, and the
// logs of the run that already happened, if there was one.
export interface DeployScreenProps {
    template: Template,
    spec: KeeperDeploySpecResponse,
    logs?: string[],
    onDeployed: (logs: string[]) => void,
}
