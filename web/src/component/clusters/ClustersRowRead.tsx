import {DeleteIconButton, EditIconButton} from "../view/IconButtons";
import {useMutation, useQueryClient} from "@tanstack/react-query";
import {clusterApi} from "../../app/api";
import {ClusterMap} from "../../app/types";
import {Box} from "@mui/material";
import {useMutationOptions} from "../../app/hooks";

type Props = {
    name: string
    toggle: () => void
    onDelete?: () => void
}

export function ClustersRowRead(props: Props) {
    const { toggle, onDelete, name } = props
    const { onError } = useMutationOptions()
    const queryClient = useQueryClient();
    const deleteCluster = useMutation(clusterApi.delete, {
        onSuccess: (_, newName) => {
            const map = queryClient.getQueryData<ClusterMap>(["cluster/list"]) ?? {}
            delete map[newName]
            queryClient.setQueryData<ClusterMap>(["cluster/list"], map)
            if (onDelete) onDelete()
        },
        onError,
    })

    return (
        <Box display={"flex"}>
            <EditIconButton loading={false} disabled={deleteCluster.isLoading} onClick={toggle}/>
            <DeleteIconButton loading={deleteCluster.isLoading} onClick={handleDelete}/>
        </Box>
    )

    function handleDelete() {
        deleteCluster.mutate(name)
    }
}
