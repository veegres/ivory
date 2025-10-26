import type {Cluster} from "../src/api/cluster/type"
import type {InstanceWeb} from "../src/api/instance/type"

export function createMockCluster(overrides: Partial<Cluster> = {}): Cluster {
    return {
        name: "test-cluster",
        instances: [{host: "localhost", port: 8008}],
        tls: {sidecar: false, database: false},
        certs: {},
        credentials: {},
        tags: [],
        ...overrides,
    }
}

export function createMockInstanceWeb(overrides: Partial<InstanceWeb> = {}): InstanceWeb {
    return {
        sidecar: {host: "localhost", port: 8008},
        state: "running",
        role: "leader",
        lag: 0,
        pendingRestart: false,
        database: {host: "localhost", port: 5432, name: "postgres"},
        leader: true,
        inCluster: true,
        inInstances: true,
        ...overrides,
    }
}
