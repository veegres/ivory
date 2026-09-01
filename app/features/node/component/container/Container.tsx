import {Box} from "@mui/material"

import {ErrorSshMissing} from "../../../../shared/component/box/ErrorManual"
import {TabsButton} from "../../../../shared/component/button/TabsButton"
import {SxPropsMap} from "../../../../shared/helper/HelperType"
import {useStore, useStoreAction} from "../../../../shared/provider/StoreProvider"
import {Feature} from "../../../Feature"
import {ManageAccessBox} from "../../../management/component/ManageAccess"
import {PlatformVaultConnection} from "../../api/NodeType"
import {ContainerList} from "./ContainerList"
import {ContainerOverview} from "./ContainerOverview"

const TABS = [{label: "overview"}, {label: "list"}]

const SX: SxPropsMap = {
    box: {display: "flex", flexDirection: "column", justifyContent: "center", gap: 1},
    actions: {alignSelf: "flex-start"},
}

type Props = {
    connection?: PlatformVaultConnection,
    name: string,
}

export function Container(props: Props) {
    const {connection, name} = props
    const tab = useStore(s => s.nodeState.containerTab)
    const {setContainerTab} = useStoreAction

    if (!connection) return <ErrorSshMissing/>
    return (
        <ManageAccessBox sx={SX.box} feature={Feature.ViewNodePlatformContainer} error={true}>
            <Box sx={SX.actions}><TabsButton tabs={TABS} tab={tab} setTab={setContainerTab}/></Box>
            {tab === 0 && <ContainerOverview connection={connection} name={name}/>}
            {tab === 1 && <ContainerList connection={connection}/>}
        </ManageAccessBox>
    )
}
