import {Permission} from "../../../../api/permission/type"
import {useRouterVault} from "../../../../api/vault/hook"
import {LinearProgressStateful} from "../../../view/progress/LinearProgressStateful"
import {Access} from "../../access/Access"
import {MenuWrapper} from "../menu/MenuWrapper"
import {VaultList} from "./VaultList"
import {VaultNew} from "./VaultNew"

export function Vault() {
    const query = useRouterVault()
    const {data, error, isFetching} = query

    return (
        <MenuWrapper>
            <Access permission={Permission.ManageVaultCreate}>
                <VaultNew/>
            </Access>
            <LinearProgressStateful color={"inherit"} loading={isFetching} line/>
            <VaultList vaults={data} error={error}/>
        </MenuWrapper>
    )
}
