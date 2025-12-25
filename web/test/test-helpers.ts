import type {Cluster, Instance} from "../src/api/cluster/type"

export function createMockCluster(overrides: Partial<Cluster> = {}): Cluster {
    return {
        name: "test-cluster",
        sidecars: [{host: "localhost", port: 8008}],
        sidecarsOverview: {},
        tls: {sidecar: false, database: false},
        certs: {},
        credentials: {},
        tags: [],
        ...overrides,
    }
}

export function createMockInstance(overrides: Partial<Instance> = {}): Instance {
    return {
        sidecar: {host: "localhost", port: 8008},
        state: "running",
        role: "leader",
        lag: 0,
        pendingRestart: false,
        database: {host: "localhost", port: 5432, name: "postgres"},
        inCluster: true,
        inSidecar: true,
        ...overrides,
    }
}

// Alias for backwards compatibility
export const createMockInstanceWeb = createMockInstance
