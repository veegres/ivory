import {api} from "../api";
import {Cluster, ClusterAuto, ClusterMap} from "./type";
import {R} from "../management/type";

export const ClusterApi = {
    list: {
        key: () => ["cluster", "list"],
        fn: (tags?: string[]) => api.get<R<Cluster[]>>(`/cluster`, {params: {tags}})
            .then((response) => response.data.response.reduce(
                (map, cluster) => {
                    map[cluster.name] = cluster
                    return map
                },
                {} as ClusterMap
            )),
    },
    update: {
        key: () => ["cluster", "update"],
        fn: (cluster: Cluster) => api.put<R<Cluster>>(`/cluster`, cluster)
            .then((response) => response.data.response),
    },
    createAuto: {
        key: () => ["cluster", "auto", "creation"],
        fn: (cluster: ClusterAuto) => api.post<R<Cluster>>(`/cluster/auto`, cluster)
            .then((response) => response.data.response),
    },
    delete: {
        key: () => ["cluster", "delete"],
        fn: (name: string) => api.delete(`/cluster/${name}`)
            .then((response) => response.data.response),
    },
}