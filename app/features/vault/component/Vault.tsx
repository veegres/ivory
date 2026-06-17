import {useState} from "react"

import {TabsButton} from "../../../shared/component/button/TabsButton"
import {LinearProgressStateful} from "../../../shared/component/progress/LinearProgressStateful"
import {LastElementScrolling} from "../../../shared/component/scrolling/LastElementScrolling"
import {Feature} from "../../feature"
import {ManageAccess} from "../../management/component/ManageAccess"
import {useRouterVault} from "../api/hook"
import {VaultTabs, VaultType} from "../api/type"
import {VaultList} from "./VaultList"
import {VaultNewPassword} from "./VaultNewPassword"
import {VaultNewSshKey} from "./VaultNewSshKey"

export const TABS: VaultTabs = {
    0: {label: "DATABASE PASS", type: VaultType.DATABASE_PASSWORD},
    1: {label: "KEEPER PASS", type: VaultType.KEEPER_PASSWORD},
    2: {label: "SSH KEY", type: VaultType.SSH_KEY},
}

export function Vault() {
    const [tab, setTab] = useState(0)
    const type = TABS[tab].type
    const query = useRouterVault(type)
    const {data, error, isFetching} = query

    return (
        <LastElementScrolling>
            <TabsButton tabs={TABS} tab={tab} setTab={setTab}/>
            <ManageAccess feature={Feature.ManageVaultCreate}>
                {type === VaultType.SSH_KEY ? (
                    <VaultNewSshKey/>
                ) : (
                    <VaultNewPassword type={type}/>
                )}
            </ManageAccess>
            <LinearProgressStateful color={"inherit"} loading={isFetching} line/>
            <VaultList vaults={data} error={error}/>
        </LastElementScrolling>
    )
}
