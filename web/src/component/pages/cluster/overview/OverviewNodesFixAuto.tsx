import {useState} from "react"

import {useRouterClusterFixAuto} from "../../../../api/cluster/hook"
import {Permission} from "../../../../api/permission/type"
import {AutoIconButton} from "../../../view/button/IconButtons"
import {AlertDialog} from "../../../view/dialog/AlertDialog"
import {Access} from "../../../widgets/access/Access"

type Props = {
    name: string,
}

export function OverviewNodesFixAuto(props: Props) {
    const {name} = props
    const autoFix = useRouterClusterFixAuto(name)
    const [open, setOpen] = useState(false)

    return (
        <Access permission={Permission.ManageClusterUpdate}>
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
        </Access>
    )
}