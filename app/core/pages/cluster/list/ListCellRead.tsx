import {Box} from "@mui/material"

import {useRouterClusterDelete} from "../../../../features/cluster/hook"
import {Feature} from "../../../../features/feature"
import {DeleteIconButton, EditIconButton} from "../../../../shared/component/button/IconButtons"
import {SxPropsMap} from "../../../../shared/helper/type"
import {Access} from "../../../widgets/access/Access"

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
            <Access feature={Feature.ManageClusterUpdate}>
                <EditIconButton disabled={isPending || isSuccess} onClick={toggle}/>
            </Access>
            <Access feature={Feature.ManageClusterDelete}>
                <DeleteIconButton loading={isPending} disabled={isSuccess} onClick={handleDelete}/>
            </Access>
        </Box>
    )

    function handleDelete() {
        deleteCluster.mutate(name)
    }
}
