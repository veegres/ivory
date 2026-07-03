import {useState} from "react"

import {TitledBox} from "../../../../shared/component/box/TitledBox"
import {TabsButton} from "../../../../shared/component/button/TabsButton"
import {PlatformVaultConnection} from "../../api/type"
import {ContainerOverviewList} from "./ContainerOverviewList"
import {ContainerOverviewMain} from "./ContainerOverviewMain"

type Props = {
    connection: PlatformVaultConnection,
}

export function ContainerOverview(props: Props) {
    const {connection} = props
    const [tab, setTab] = useState(0)

    return (
        <TitledBox title={"Container"} renderActions={renderActions()} island={true}>
            {tab === 0 && <ContainerOverviewMain connection={connection}/>}
            {tab === 1 && <ContainerOverviewList connection={connection}/>}
        </TitledBox>
    )

    function renderActions() {
        const tabs = [{label: "main"}, {label: "list"}]
        return (
            <TabsButton tabs={tabs} tab={tab} setTab={setTab}/>
        )
    }
}
