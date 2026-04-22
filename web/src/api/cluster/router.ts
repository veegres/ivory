import {getDomain} from "../../app/utils"
import {api} from "../api"
import {R} from "../management/type"
import {Connection} from "../node/type"
import {Cluster, ClusterAuto, ClusterOverview} from "./type"

export const ClusterApi = {
    list: {
        key: () => ["cluster", "list"],
        fn: (tags?: string[]) => api.get<R<Cluster[]>>("/cluster", {params: {tags}})
            .then((response) => response.data.response.map(v => (
                {...v, nodesOverview: Object.fromEntries(v.nodes.map(i => [getDomain(i), undefined]))} as Cluster
            ))),
    },
    overview: {
        key: (name?: string, connection?: Connection) => ["cluster", "overview", name, connection?.host, connection?.keeperPort].filter(Boolean),
        fn: (name: string, connection?: Connection) => api.get<R<ClusterOverview>>(`/cluster/overview/${name}`, {params: {keeper: JSON.stringify(connection)}})
            .then((response) => response.data.response),
    },
    update: {
        key: () => ["cluster", "update"],
        fn: (cluster: Cluster) => api.put<R<Cluster>>("/cluster", cluster)
            .then((response) => response.data.response),
    },
    createAuto: {
        key: () => ["cluster", "auto", "creation"],
        fn: (cluster: ClusterAuto) => api.post<R<Cluster>>("/cluster/auto", cluster)
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
