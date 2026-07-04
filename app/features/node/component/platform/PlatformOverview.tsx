import {TitledBox} from "../../../../shared/component/box/TitledBox"
import {TabsButton} from "../../../../shared/component/button/TabsButton"
import {useStore, useStoreAction} from "../../../../shared/provider/StoreProvider"
import {PlatformVaultConnection} from "../../api/type"
import {PlatformLogs} from "./PlatformLogs"
import {PlatformProcesses} from "./PlatformProcesses"

type Props = {
    connection: PlatformVaultConnection,
}

export function PlatformOverview(props: Props) {
    const {connection} = props
    const tab = useStore(s => s.nodeState.platformTab)
    const {setPlatformTab} = useStoreAction

    return (
        <TitledBox title={"Platform"} renderActions={renderActions()} island={true}>
            {tab === 0 && <PlatformProcesses connection={connection}/>}
            {tab === 1 && <PlatformLogs connection={connection}/>}
        </TitledBox>
    )

    function renderActions() {
        const tabs = [{label: "processes"}, {label: "logs"}]
        return (
            <TabsButton tabs={tabs} tab={tab} setTab={setPlatformTab}/>
        )
    }
}
