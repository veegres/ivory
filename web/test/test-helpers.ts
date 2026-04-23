import type {Cluster, Node} from "../src/api/cluster/type"

export function createMockCluster(overrides: Partial<Cluster> = {}): Cluster {
    return {
        name: "test-cluster",
        nodes: [{sshKeyId: "00000000-0000-0000-0000-000000000000", host: "localhost", keeperPort: 8008}],
        nodesOverview: {},
        tls: {keeper: false, database: false},
        certs: {},
        vaults: {},
        tags: [],
        ...overrides,
    }
}

export function createMockNode(overrides: Partial<Node> = {}): Node {
    return {
        connection: {sshKeyId: "00000000-0000-0000-0000-000000000000", host: "localhost", keeperPort: 8008},
        keeper: {
            state: "running",
            role: "leader",
            lag: 0,
            pendingRestart: false,
            discoveredHost: "localhost",
            discoveredKeeperPort: 8008,
            discoveredDbPort: 5432,
        },
        inCluster: true,
        inKeeper: false,
        ...overrides,
    }
}

export function setupLocalStorageMock() {
    const localStorageMock = (() => {
        let store: Record<string, string> = {}
        return {
            getItem: (key: string) => (key in store) ? store[key] : null,
            setItem: (key: string, value: string) => {
                store[key] = value.toString()
            },
            clear: () => {
                store = {}
            },
            removeItem: (key: string) => {
                delete store[key]
            },
            get length() {
                return Object.keys(store).length
            },
            key: (index: number) => Object.keys(store)[index] || null,
        }
    })()

    Object.defineProperty(window, "localStorage", {
        value: localStorageMock,
        writable: true,
    })
}
