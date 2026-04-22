import type {Cluster, Node} from "../src/api/cluster/type"

export function createMockCluster(overrides: Partial<Cluster> = {}): Cluster {
    return {
        name: "test-cluster",
        keepers: [{host: "localhost", port: 8008}],
        keepersOverview: {},
        tls: {keeper: false, database: false},
        certs: {},
        credentials: {},
        tags: [],
        ...overrides,
    }
}

export function createMockNode(overrides: Partial<Node> = {}): Node {
    return {
        keeper: {host: "localhost", port: 8008},
        state: "running",
        role: "leader",
        lag: 0,
        pendingRestart: false,
        database: {host: "localhost", port: 5432, name: "postgres"},
        inCluster: true,
        inKeeper: true,
        ...overrides,
    }
}

// Alias for backwards compatibility
export const createMockInstance = createMockNode
export const createMockInstanceWeb = createMockNode
