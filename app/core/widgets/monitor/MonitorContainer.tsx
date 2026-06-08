import {useState} from "react"

import {PlatformConnection} from "../../../features/node/type"
import {TitledBox} from "../../../shared/component/box/TitledBox"
import {TabsButton} from "../../../shared/component/button/TabsButton"
import {MonitorContainerList} from "./MonitorContainerList"
import {MonitorContainerSingle} from "./MonitorContainerSingle"

type Props = {
    connection: PlatformConnection,
}

export function MonitorContainer(props: Props) {
    const {connection} = props
    const [tab, setTab] = useState(0)

    return (
        <TitledBox title={"Container"} renderActions={renderActions()} island={true}>
            {tab === 0 && <MonitorContainerSingle connection={connection}/>}
            {tab === 1 && <MonitorContainerList connection={connection}/>}
        </TitledBox>
    )

    function renderActions() {
        const tabs = [{label: connection.host}, {label: "list"}]
        return (
            <TabsButton tabs={tabs} tab={tab} setTab={setTab}/>
        )
    }
}
