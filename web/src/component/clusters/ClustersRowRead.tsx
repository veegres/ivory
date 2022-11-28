import {DeleteIconButton, EditIconButton} from "../view/IconButtons";
import {useMutation} from "@tanstack/react-query";
import {clusterApi} from "../../app/api";
import {Box} from "@mui/material";
import {useMutationOptions} from "../../app/hooks";

type Props = {
    name: string
    toggle: () => void
    onDelete?: () => void
}

export function ClustersRowRead(props: Props) {
    const { toggle, onDelete, name } = props

    const deleteMutationOptions = useMutationOptions(["cluster/list"], onDelete)
    const deleteCluster = useMutation(clusterApi.delete, deleteMutationOptions)

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
