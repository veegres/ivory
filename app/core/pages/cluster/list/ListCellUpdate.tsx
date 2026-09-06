import {Box} from "@mui/material"

import {useRouterClusterCreate, useRouterClusterUpdate} from "../../../../features/cluster/api/ClusterHook"
import {Cluster} from "../../../../features/cluster/api/ClusterType"
import {CancelIconButton, SaveIconButton} from "../../../../shared/component/button/IconButtons"
import {SxPropsMap} from "../../../../shared/helper/HelperType"

const SX: SxPropsMap = {
    box: {display: "flex", justifyContent: "flex-end"},
}

type Props = {
    cluster: Cluster,
    toggle: () => void,
    onUpdate?: () => void,
    onClose?: () => void,
    create?: boolean,
}

export function ListCellUpdate(props: Props) {
    const {cluster} = props
    const {toggle, onUpdate, onClose, create = false} = props

    const createCluster = useRouterClusterCreate(handleSuccess)
    const updateCluster = useRouterClusterUpdate(cluster.name, handleSuccess)
    const isPending = createCluster.isPending || updateCluster.isPending

    return (
        <Box sx={SX.box}>
            <CancelIconButton loading={false} disabled={isPending} onClick={handleClose}/>
            <SaveIconButton loading={isPending} disabled={!cluster.name} onClick={handleUpdate}/>
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
        if (create) {
            createCluster.mutate(cluster)
            return
        }
        updateCluster.mutate(cluster)
    }
}
