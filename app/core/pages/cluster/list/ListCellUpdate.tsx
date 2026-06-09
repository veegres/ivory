import {Box} from "@mui/material"

import {useRouterClusterUpdate} from "../../../../features/cluster/hook"
import {Cluster} from "../../../../features/cluster/type"
import {CancelIconButton, SaveIconButton} from "../../../../shared/component/button/IconButtons"
import {SxPropsMap} from "../../../../shared/helper/type"

const SX: SxPropsMap = {
    box: {display: "flex", justifyContent: "flex-end"},
}

type Props = {
    cluster: Cluster,
    toggle: () => void,
    onUpdate?: () => void,
    onClose?: () => void,
}

export function ListCellUpdate(props: Props) {
    const {cluster} = props
    const {toggle, onUpdate, onClose} = props

    const updateCluster = useRouterClusterUpdate(cluster.name, handleSuccess)

    return (
        <Box sx={SX.box}>
            <CancelIconButton loading={false} disabled={updateCluster.isPending} onClick={handleClose}/>
            <SaveIconButton loading={updateCluster.isPending} disabled={!cluster.name} onClick={handleUpdate}/>
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
        updateCluster.mutate(cluster)
    }
}
