import {getDomain, getInitialNode} from "../../shared/helper/utils"
import {api} from "../api"
import {R} from "../management/type"
import {AutoRequest, Cluster, DeployRequest, Overview} from "./type"

export const ClusterApi = {
    list: {
        key: () => ["cluster", "list"],
        fn: (tags?: string[]) => api.get<R<Cluster[]>>("/cluster", {params: {tags}})
            .then((response) => response.data.response.map(v => (
                {...v, nodesOverview: Object.fromEntries(v.nodes.map(c => {
                    const domain = getDomain(c)
                    return [domain, getInitialNode(c)]
                }))} as Cluster
            ))),
    },
    overview: {
        key: (name?: string, host?: string, port?: string) => ["cluster", "overview", name, host, port],
        fn: (name: string, host?: string, port?: string) => api
            .get<R<Overview>>(`/cluster/overview/${name}`, {params: {host, port}})
            .then((response) => response.data.response),
    },
    update: {
        key: () => ["cluster", "update"],
        fn: (cluster: Cluster) => api.put<R<Cluster>>("/cluster", cluster)
            .then((response) => response.data.response),
    },
    createAuto: {
        key: () => ["cluster", "auto", "creation"],
        fn: (cluster: AutoRequest) => api.post<R<Cluster>>("/cluster/auto", cluster)
            .then((response) => response.data.response),
    },
    deploy: {
        key: () => ["cluster", "deploy"],
        fn: (req: DeployRequest) => api.post<R<string[]>>("/cluster/deploy", req)
            .then((response) => response.data.response),
    },
    fixAuto: {
        key: () => ["cluster", "auto", "fix"],
        fn: (name: string) => api.post<R<Cluster>>(`/cluster/auto/${name}`)
            .then((response) => response.data.response),
    },
    delete: {
        key: () => ["cluster", "delete"],
        fn: (name: string) => api.delete(`/cluster/${name}`)
            .then((response) => response.data.response),
    },
}
