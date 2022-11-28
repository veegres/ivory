import {CancelIconButton, SaveIconButton} from "../view/IconButtons";
import {useMutation, useQueryClient} from "@tanstack/react-query";
import {clusterApi} from "../../app/api";
import {ClusterMap} from "../../app/types";
import {Box} from "@mui/material";
import {useToast} from "../../app/hooks";

type Props = {
    name: string
    nodes: string[]
    toggle: () => void
    onUpdate?: () => void
    onClose?: () => void
}

export function ClustersRowUpdate(props: Props) {
    const { toggle, onUpdate, onClose, name, nodes } = props
    const { onError } = useToast()

    const queryClient = useQueryClient();
    const updateCluster = useMutation(clusterApi.update, {
        onSuccess: (data) => {
            const map = queryClient.getQueryData<ClusterMap>(["cluster/list"]) ?? {} as ClusterMap
            map[data.name] = data
            queryClient.setQueryData<ClusterMap>(["cluster/list"], map)
            toggle()
            if (onUpdate) onUpdate()
        },
        onError,
    })

    return (
        <Box display={"flex"}>
            <CancelIconButton loading={false} disabled={updateCluster.isLoading} onClick={handleClose}/>
            <SaveIconButton loading={updateCluster.isLoading} disabled={!name} onClick={handleUpdate}/>
        </Box>
    )

    function handleClose() {
        toggle()
        if (onClose) onClose()
    }

    function handleUpdate() {
        updateCluster.mutate({ name , nodes })
    }
}
