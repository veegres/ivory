import {CancelIconButton, SaveIconButton} from "../../../view/button/IconButtons";
import {useMutation} from "@tanstack/react-query";
import {ClusterApi} from "../../../../app/api";
import {Box} from "@mui/material";
import {useMutationOptions} from "../../../../hook/QueryCustom";
import {Sidecar, SxPropsMap} from "../../../../type/general";
import {Certs, Credentials} from "../../../../type/cluster";

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

    const updateMutationOptions = useMutationOptions([["cluster", "list"], ["tag", "list"]], handleSuccess)
    const updateCluster = useMutation({mutationFn: ClusterApi.update, ...updateMutationOptions})

    return (
        <Box sx={SX.box}>
            <CancelIconButton loading={false} disabled={updateCluster.isPending} onClick={handleClose}/>
            <SaveIconButton loading={updateCluster.isPending} disabled={!name} onClick={handleUpdate}/>
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
        const tmp = tags?.filter(t => t !== "ALL")
        updateCluster.mutate({name, instances, certs, credentials, tags: tmp})
    }
}
