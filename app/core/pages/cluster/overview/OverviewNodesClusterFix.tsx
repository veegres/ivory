import {AutoFixHigh} from "@mui/icons-material"
import {useState} from "react"

import {useRouterClusterFix} from "../../../../features/cluster/api/ClusterHook"
import {Feature} from "../../../../features/Feature"
import {ManageAccess} from "../../../../features/management/component/ManageAccess"
import {SimpleButton} from "../../../../shared/component/button/SimpleButton"
import {AlertDialog} from "../../../../shared/component/dialog/AlertDialog"

type Props = {
    name: string,
}

export function OverviewNodesClusterFix(props: Props) {
    const {name} = props
    const autoFix = useRouterClusterFix(name)
    const [open, setOpen] = useState(false)

    return (
        <ManageAccess feature={Feature.ManageClusterUpdate}>
            <SimpleButton tooltip={"Fix Cluster"} loading={autoFix.isPending} onClick={() => setOpen(true)}>
                <AutoFixHigh fontSize={"small"}/>
            </SimpleButton>
            <AlertDialog
                open={open}
                title={"Fix Cluster"}
                description={"Update cluster to match the keeper's current state."}
                onClose={() => setOpen(false)}
                onAgree={() => autoFix.mutate(name)}
            />
        </ManageAccess>
    )
}
