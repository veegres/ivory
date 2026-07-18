import {AutoFixHigh} from "@mui/icons-material"
import {Box, Tooltip} from "@mui/material"
import {useState} from "react"

import {useRouterClusterFixAuto} from "../../../../features/cluster/api/ClusterHook"
import {Feature} from "../../../../features/Feature"
import {ManageAccess} from "../../../../features/management/component/ManageAccess"
import {SimpleButton} from "../../../../shared/component/button/SimpleButton"
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
            <Tooltip title={"Auto Fix"} placement={"top"} arrow={true}>
                <Box component={"span"}>
                    <SimpleButton loading={autoFix.isPending} onClick={() => setOpen(true)}>
                        <AutoFixHigh/>
                    </SimpleButton>
                </Box>
            </Tooltip>
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
