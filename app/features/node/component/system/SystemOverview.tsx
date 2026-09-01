import {Box} from "@mui/material"

import {TabsButton} from "../../../../shared/component/button/TabsButton"
import {SxPropsMap} from "../../../../shared/helper/HelperType"
import {useStore, useStoreAction} from "../../../../shared/provider/StoreProvider"
import {PlatformVaultConnection} from "../../api/NodeType"
import {SystemLogs} from "./SystemLogs"
import {SystemProcesses} from "./SystemProcesses"

const TABS = [{label: "processes"}, {label: "logs"}]

const SX: SxPropsMap = {
    actions: {alignSelf: "flex-start"},
}

type Props = {
    connection: PlatformVaultConnection,
}

export function SystemOverview(props: Props) {
    const {connection} = props
    const tab = useStore(s => s.nodeState.systemTab)
    const {setSystemTab} = useStoreAction

    return (
        <>
            <Box sx={SX.actions}><TabsButton tabs={TABS} tab={tab} setTab={setSystemTab}/></Box>
            {tab === 0 && <SystemProcesses connection={connection}/>}
            {tab === 1 && <SystemLogs connection={connection}/>}
        </>
    )
}
