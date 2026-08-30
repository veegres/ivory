import {SubContentBox} from "../../../../shared/component/box/SubContentBox"
import {TabsButton} from "../../../../shared/component/button/TabsButton"
import {useStore, useStoreAction} from "../../../../shared/provider/StoreProvider"
import {PlatformVaultConnection} from "../../api/NodeType"
import {SystemLogs} from "./SystemLogs"
import {SystemProcesses} from "./SystemProcesses"

type Props = {
    connection: PlatformVaultConnection,
}

export function SystemOverview(props: Props) {
    const {connection} = props
    const tab = useStore(s => s.nodeState.systemTab)
    const {setSystemTab} = useStoreAction

    return (
        <SubContentBox label={"System"} renderActions={renderActions()} island={true} collapsible={false}>
            {tab === 0 && <SystemProcesses connection={connection}/>}
            {tab === 1 && <SystemLogs connection={connection}/>}
        </SubContentBox>
    )

    function renderActions() {
        const tabs = [{label: "processes"}, {label: "logs"}]
        return (
            <TabsButton tabs={tabs} tab={tab} setTab={setSystemTab}/>
        )
    }
}
