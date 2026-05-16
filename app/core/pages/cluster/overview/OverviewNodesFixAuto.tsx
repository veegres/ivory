import {useState} from "react"

import {useRouterClusterFixAuto} from "../../../../features/cluster/hook"
import {Feature} from "../../../../features/feature"
import {AutoIconButton} from "../../../../shared/component/button/IconButtons"
import {AlertDialog} from "../../../../shared/component/dialog/AlertDialog"
import {Access} from "../../../widgets/access/Access"

type Props = {
    name: string,
}

export function OverviewNodesFixAuto(props: Props) {
    const {name} = props
    const autoFix = useRouterClusterFixAuto(name)
    const [open, setOpen] = useState(false)

    return (
        <Access feature={Feature.ManageClusterUpdate}>
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