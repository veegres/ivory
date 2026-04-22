import {getDomain} from "../../app/utils"
import {api} from "../api"
import {R} from "../management/type"
import {Keeper} from "../node/type"
import {Cluster, ClusterAuto, ClusterOverview} from "./type"

export const ClusterApi = {
    list: {
        key: () => ["cluster", "list"],
        fn: (tags?: string[]) => api.get<R<Cluster[]>>("/cluster", {params: {tags}})
            .then((response) => response.data.response.map(v => (
                {...v, keepersOverview: Object.fromEntries(v.keepers.map(i => [getDomain(i), undefined]))} as Cluster
            ))),
    },
    overview: {
        key: (name?: string, keeper?: Keeper) => ["cluster", "overview", name, keeper?.host, keeper?.port].filter(Boolean),
        fn: (name: string, keeper?: Keeper) => api.get<R<ClusterOverview>>(`/cluster/overview/${name}`, {params: {keeper: JSON.stringify(keeper)}})
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
