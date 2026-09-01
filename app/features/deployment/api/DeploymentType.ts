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
    // NOTE: each step runs as its own execution, so nothing may assume a shell
    // - the images these run in are increasingly distroless and have no "&&"
    postScripts?: string[],
    defaults?: CommandDefaults,
}

// CommandDefaults is what this command fills its node card in with. A command
// and the endpoints it answers on are one fact rather than two: a single-host
// template writes a distinct peer port into each of its commands, and only that
// command knows which client port has to match it - sshPort included, since a
// single-host node can still be forwarded on its own port. Host is empty in a
// multi-host template - each node is a distinct real machine a template can
// never know ahead of time - but a single-host template can and does name one,
// since all of its commands land on the same machine by definition.
export interface CommandDefaults {
    name?: string,
    host?: string,
    sshPort?: number,
    keeperPort?: number,
    dbPort?: number,
}

// TemplateDefaults is what the whole template fills the deploy screen's
// credential fields with. Credentials are one answer for the whole cluster, so
// they sit here rather than on a command, where three copies could only ever
// disagree. A username means the deployment ends up with that account - spilo
// names its superuser postgres, etcd can only enable auth through root - so the
// screen opens on it; where it names none, the screen opens with that credential
// switched off. Passwords are never here: a template is stored, read and copied.
export interface TemplateDefaults {
    keeperUser?: string,
    dbUser?: string,
}

// Template is a saved deployment: an ordered list of commands, one per node,
// and nothing about where they run.
export interface Template {
    id: string,
    name: string,
    description?: string,
    keeper: KeeperPlugin,
    platform: PlatformPlugin,
    defaults?: TemplateDefaults,
    commands: TemplateCommand[],
    creation: CreationType,
    createdAt: number,
}

export interface TemplateRequest {
    name: string,
    description?: string,
    keeper: KeeperPlugin,
    platform: PlatformPlugin,
    defaults?: TemplateDefaults,
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
// runs a template: the template itself, and the logs of the run that already
// happened, if there was one. Everything a deploy needs to know about the
// engine is in the template - that is what replaced the keeper plugin's own
// deploy spec.
export interface DeployScreenProps {
    template: Template,
    logs?: string[],
    onDeployed: (logs: string[]) => void,
}
