import {AutocompleteTags} from "../../view/././autocomplete/AutocompleteTags";
import {useMutation, useQuery} from "@tanstack/react-query";
import {clusterApi, tagApi} from "../../../app/api";
import {useMutationOptions} from "../../../hook/QueryCustom";
import {Cluster} from "../../../type/cluster";

type Props = {
    cluster: Cluster
}

export function OptionsTags(props: Props) {
    const query = useQuery(["tag/list"], tagApi.list)
    const {data, isLoading} = query
    const tags = data ?? [];

    const updateMutationOptions = useMutationOptions([["cluster/list"]], query.refetch)
    const updateCluster = useMutation(clusterApi.update, updateMutationOptions)

    return (
        <AutocompleteTags
            tags={tags}
            selected={props.cluster.tags ?? []}
            loading={isLoading || updateCluster.isLoading}
            onUpdate={handleUpdate}
        />
    )

    function handleUpdate(tags: string[]) {
        updateCluster.mutate({...props.cluster, tags})
    }
}
