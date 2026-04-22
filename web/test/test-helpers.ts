import type {Cluster, Node} from "../src/api/cluster/type"

export function createMockCluster(overrides: Partial<Cluster> = {}): Cluster {
    return {
        name: "test-cluster",
        nodes: [{host: "localhost", keeperPort: 8008}],
        nodesOverview: {},
        tls: {keeper: false, database: false},
        certs: {},
        credentials: {},
        tags: [],
        ...overrides,
    }
}

export function createMockNode(overrides: Partial<Node> = {}): Node {
    return {
        connection: {host: "localhost", keeperPort: 8008},
        response: {
            state: "running",
            role: "leader",
            lag: 0,
            pendingRestart: false,
            discoveredHost: "localhost",
            discoveredKeeperPort: 8008,
            discoveredDbPort: 5432,
        },
        inCluster: true,
        inKeeper: true,
        ...overrides,
    }
}

// Alias for backwards compatibility
export const createMockInstance = createMockNode
export const createMockInstanceWeb = createMockNode
