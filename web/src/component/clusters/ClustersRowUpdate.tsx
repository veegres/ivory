import {CancelIconButton, SaveIconButton} from "../view/IconButtons";
import {useMutation} from "@tanstack/react-query";
import {clusterApi} from "../../app/api";
import {Box} from "@mui/material";
import {useMutationOptions} from "../../hook/QueryCustom";

type Props = {
    name: string
    nodes: string[]
    toggle: () => void
    onUpdate?: () => void
    onClose?: () => void
}

export function ClustersRowUpdate(props: Props) {
    const { toggle, onUpdate, onClose, name, nodes } = props

    const updateMutationOptions = useMutationOptions(["cluster/list"], handleSuccess)
    const updateCluster = useMutation(clusterApi.update, updateMutationOptions)

    return (
        <Box display={"flex"}>
            <CancelIconButton loading={false} disabled={updateCluster.isLoading} onClick={handleClose}/>
            <SaveIconButton loading={updateCluster.isLoading} disabled={!name} onClick={handleUpdate}/>
        </Box>
    )

    function handleSuccess() {
        toggle()
        if (onUpdate) onUpdate()
    }

    function handleClose() {
        toggle()
        if (onClose) onClose()
    }

    function handleUpdate() {
        updateCluster.mutate({ name , nodes })
    }
}
