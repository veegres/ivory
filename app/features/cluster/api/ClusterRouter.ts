import {getDomain, getInitialNode} from "../../../shared/helper/HelperUtils"
import {api} from "../../Api"
import {R} from "../../management/api/ManagementType"
import {KeeperPlugin} from "../../node/api/NodeType"
import {DbPlugin} from "../../query/api/QueryType"
import {AutoRequest, Cluster, DeployRequest, Overview} from "./ClusterType"

export const ClusterApi = {
    list: {
        key: (tags: string[] = [], keeper?: KeeperPlugin, database?: DbPlugin) => ["cluster", "list", keeper, database, ...tags],
        keyCommon: () => ["cluster", "list"],
        fn: (tags?: string[], keeper?: KeeperPlugin, database?: DbPlugin) => api
            .get<R<Cluster[]>>("/cluster", {params: {tags, keeper, database}})
            .then((response) => response.data.response.map(v => (
                {...v, nodesOverview: Object.fromEntries(v.nodes.map(c => {
                    const domain = getDomain(c, true)
                    return [domain, getInitialNode(c)]
                }))} as Cluster
            ))),
    },
    overview: {
        key: (name: string, host?: string, port?: string) => ["cluster", "overview", name, host, port],
        // NOTE: omit name to match every active overview query regardless of cluster
        keyCommon: (name?: string) => name ? ["cluster", "overview", name] : ["cluster", "overview"],
        fn: (name: string, host?: string, port?: string) => api
            .get<R<Overview>>(`/cluster/overview/${name}`, {params: {host, port}})
            .then((response) => response.data.response),
    },
    update: {
        key: () => ["cluster", "update"],
        fn: (cluster: Cluster) => api.put<R<Cluster>>("/cluster", cluster)
            .then((response) => response.data.response),
    },
    detect: {
        key: () => ["cluster", "detect"],
        fn: (cluster: AutoRequest) => api.post<R<Cluster>>("/cluster/detect", cluster)
            .then((response) => response.data.response),
    },
    deploy: {
        key: () => ["cluster", "deploy"],
        fn: (req: DeployRequest) => api.post<R<string[]>>("/cluster/deploy", req)
            .then((response) => response.data.response),
    },
    fix: {
        key: () => ["cluster", "fix"],
        fn: (name: string) => api.post<R<Cluster>>(`/cluster/fix/${name}`)
            .then((response) => response.data.response),
    },
    delete: {
        key: () => ["cluster", "delete"],
        fn: (name: string) => api.delete(`/cluster/${name}`)
            .then((response) => response.data.response),
    },
}
