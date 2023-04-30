import {CancelIconButton, SaveIconButton} from "../../view/button/IconButtons";
import {useMutation} from "@tanstack/react-query";
import {clusterApi} from "../../../app/api";
import {Box} from "@mui/material";
import {useMutationOptions} from "../../../hook/QueryCustom";
import {Sidecar, SxPropsMap} from "../../../type/common";
import {Certs, Credentials} from "../../../type/cluster";

const SX: SxPropsMap = {
    box: {display: "flex", justifyContent: "flex-end"},
}

type Props = {
    name: string,
    instances: Sidecar[],
    certs: Certs,
    credentials: Credentials,
    tags?: string[],
    toggle: () => void,
    onUpdate?: () => void,
    onClose?: () => void,
}

export function ListCellUpdate(props: Props) {
    const {name, instances, tags, credentials, certs} = props
    const {toggle, onUpdate, onClose} = props

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
        updateCluster.mutate({name, instances, certs, credentials, tags})
    }
}
