import {Box} from "@mui/material"

import {useRouterClusterDelete} from "../../../../features/cluster/api/ClusterHook"
import {Feature} from "../../../../features/Feature"
import {ManageAccess} from "../../../../features/management/component/ManageAccess"
import {DeleteIconButton, EditIconButton} from "../../../../shared/component/button/IconButtons"
import {SxPropsMap} from "../../../../shared/helper/HelperType"

const SX: SxPropsMap = {
    box: {display: "flex", justifyContent: "flex-end"},
}

type Props = {
    name: string
    toggle: () => void
}

export function ListCellRead(props: Props) {
    const {toggle, name} = props

    const deleteCluster = useRouterClusterDelete()
    const {isPending, isSuccess} = deleteCluster

    return (
        <Box sx={SX.box}>
            <ManageAccess feature={Feature.ManageClusterUpdate}>
                <EditIconButton disabled={isPending || isSuccess} onClick={toggle}/>
            </ManageAccess>
            <ManageAccess feature={Feature.ManageClusterDelete}>
                <DeleteIconButton loading={isPending} disabled={isSuccess} onClick={handleDelete}/>
            </ManageAccess>
        </Box>
    )

    function handleDelete() {
        deleteCluster.mutate(name)
    }
}
