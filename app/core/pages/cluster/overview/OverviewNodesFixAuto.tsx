import {useState} from "react"

import {useRouterClusterFixAuto} from "../../../../features/cluster/api/ClusterHook"
import {Feature} from "../../../../features/Feature"
import {ManageAccess} from "../../../../features/management/component/ManageAccess"
import {AutoIconButton} from "../../../../shared/component/button/IconButtons"
import {AlertDialog} from "../../../../shared/component/dialog/AlertDialog"

type Props = {
    name: string,
}

export function OverviewNodesFixAuto(props: Props) {
    const {name} = props
    const autoFix = useRouterClusterFixAuto(name)
    const [open, setOpen] = useState(false)

    return (
        <ManageAccess feature={Feature.ManageClusterUpdate}>
            <AutoIconButton
                tooltip={"Auto Fix"}
                onClick={() => setOpen(true)}
                loading={autoFix.isPending}
            />
            <AlertDialog
                open={open}
                title={"Auto Fix"}
                description={"Update cluster to match the keeper's current state."}
                onClose={() => setOpen(false)}
                onAgree={() => autoFix.mutate(name)}
            />
        </ManageAccess>
    )
}