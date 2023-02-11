import {CancelIconButton, SaveIconButton} from "../../view/IconButtons";
import {useMutation} from "@tanstack/react-query";
import {clusterApi} from "../../../app/api";
import {Box} from "@mui/material";
import {useMutationOptions} from "../../../hook/QueryCustom";
import {Sidecar, SxPropsMap} from "../../../app/types";
import {getHostAndPort} from "../../../app/utils";

const SX: SxPropsMap = {
    box: {display: "flex", justifyContent: "flex-end"},
}

type Props = {
    name: string
    nodes: string[]
    toggle: () => void
    onUpdate?: () => void
    onClose?: () => void
}

export function ListRowUpdate(props: Props) {
    const {toggle, onUpdate, onClose, name, nodes} = props

    const updateMutationOptions = useMutationOptions([["cluster/list"]], handleSuccess)
    const updateCluster = useMutation(clusterApi.update, updateMutationOptions)

    return (
        <Box sx={SX.box}>
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
        const instances: Sidecar[] = nodes.map(value => getHostAndPort(value))
        updateCluster.mutate({name, instances, certs: {}, credentials: {}, tags: []})
    }
}
