import {vi} from "vitest"

import type {Cluster, Node} from "../../features/cluster/api/type"
import {KeeperPlugin} from "../../features/node/api/type"
import {DbPlugin} from "../../features/query/api/type"

/**
 * Mock localStorage
 */
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

/**
 * Mock ResizeObserver
 */
export class ResizeObserverMock {
    callback: ResizeObserverCallback
    constructor(callback: ResizeObserverCallback) {
        this.callback = callback
    }
    observe() {
        // Immediately call the callback to simulate resize
        this.callback([] as any, this as any)
    }
    disconnect() {}
    unobserve() {}
}

/**
 * Mock MutationObserver
 */
export class MutationObserverMock {
    callback: MutationCallback
    constructor(callback: MutationCallback) {
        this.callback = callback
    }
    observe() {
        // Immediately call the callback to simulate mutation
        this.callback([] as any, this as any)
    }
    disconnect() {}
    takeRecords() {
        return []
    }
}

/**
 * Mock matchMedia
 */
export function setupMatchMediaMock() {
    Object.defineProperty(window, "matchMedia", {
        writable: true,
        value: vi.fn().mockImplementation(query => ({
            matches: false,
            media: query,
            onchange: null,
            addListener: vi.fn(), // deprecated
            removeListener: vi.fn(), // deprecated
            addEventListener: vi.fn(),
            removeEventListener: vi.fn(),
            dispatchEvent: vi.fn(),
        })),
    })
}

/**
 * Create a mock cluster object
 */
export function createMockCluster(overrides: Partial<Cluster> = {}): Cluster {
    return {
        name: "test-cluster",
        plugins: {
            database: DbPlugin.POSTGRES,
            keeper: KeeperPlugin.PATRONI_POSTGRES,
        },
        nodes: [{host: "localhost", keeperPort: 8008}],
        tls: {keeper: false, database: false},
        certs: {},
        vaults: {sshKeyId: "00000000-0000-0000-0000-000000000000"},
        tags: [],
        ...overrides,
    }
}

/**
 * Create a mock node object
 */
export function createMockNode(overrides: Partial<Node> = {}): Node {
    return {
        config: {host: "localhost", keeperPort: 8008},
        keeper: {
            state: "running",
            role: "leader",
            lag: 0,
            pendingRestart: false,
            discoveredHost: "localhost",
            discoveredKeeperPort: 8008,
            discoveredDbPort: 5432,
        },
        warnings: [],
        ...overrides,
    }
}
