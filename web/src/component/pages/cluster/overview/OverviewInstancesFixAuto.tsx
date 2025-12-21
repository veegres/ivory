import {useState} from "react"

import {useRouterClusterFixAuto} from "../../../../api/cluster/hook"
import {Permission} from "../../../../api/permission/type"
import {AutoIconButton} from "../../../view/button/IconButtons"
import {AlertDialog} from "../../../view/dialog/AlertDialog"
import {Access} from "../../../widgets/access/Access"

type Props = {
    name: string,
}

export function OverviewInstancesFixAuto(props: Props) {
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
                description={"This will update cluster with current state of the cluster from the source."}
                onClose={() => setOpen(false)}
                onAgree={() => autoFix.mutate(name)}
            />
        </Access>
    )
}