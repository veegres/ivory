import {useState} from "react"

import {Feature} from "../../../../api/feature"
import {useRouterVault} from "../../../../api/vault/hook"
import {VaultTabs, VaultType} from "../../../../api/vault/type"
import {TabsButton} from "../../../view/button/TabsButton"
import {LinearProgressStateful} from "../../../view/progress/LinearProgressStateful"
import {Access} from "../../access/Access"
import {MenuWrapper} from "../menu/MenuWrapper"
import {VaultItemNew} from "./VaultItemNew"
import {VaultList} from "./VaultList"

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
        <MenuWrapper>
            <TabsButton tabs={TABS} tab={tab} setTab={setTab}/>
            <Access feature={Feature.ManageVaultCreate}>
                <VaultItemNew type={type}/>
            </Access>
            <LinearProgressStateful color={"inherit"} loading={isFetching} line/>
            <VaultList vaults={data} error={error}/>
        </MenuWrapper>
    )
}
