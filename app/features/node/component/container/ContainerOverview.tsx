import {Box} from "@mui/material"

import {TabsButton} from "../../../../shared/component/button/TabsButton"
import {SxPropsMap} from "../../../../shared/helper/HelperType"
import {useStore, useStoreAction} from "../../../../shared/provider/StoreProvider"
import {PlatformVaultConnection} from "../../api/NodeType"
import {ContainerList} from "./ContainerList"
import {ContainerLogs} from "./ContainerLogs"

const TABS = [{label: "logs"}, {label: "list"}]

const SX: SxPropsMap = {
    actions: {alignSelf: "flex-start"},
}

type Props = {
    connection: PlatformVaultConnection,
    name: string,
}

export function ContainerOverview(props: Props) {
    const {connection, name} = props
    const tab = useStore(s => s.nodeState.containerTab)
    const {setContainerTab} = useStoreAction

    return (
        <>
            <Box sx={SX.actions}><TabsButton tabs={TABS} tab={tab} setTab={setContainerTab}/></Box>
            {tab === 0 && <ContainerLogs connection={connection} name={name}/>}
            {tab === 1 && <ContainerList connection={connection}/>}
        </>
    )
}
