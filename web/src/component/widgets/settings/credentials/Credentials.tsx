import {useRouterPassword} from "../../../../api/password/hook"
import {Permission} from "../../../../api/permission/type"
import {LinearProgressStateful} from "../../../view/progress/LinearProgressStateful"
import {Access} from "../../access/Access"
import {MenuWrapper} from "../menu/MenuWrapper"
import {CredentialsList} from "./CredentialsList"
import {CredentialsNew} from "./CredentialsNew"

export function Credentials() {
    const query = useRouterPassword()
    const {data, error, isFetching} = query

    return (
        <MenuWrapper>
            <Access permission={Permission.ManagePasswordCreate}>
                <CredentialsNew/>
            </Access>
            <LinearProgressStateful color={"inherit"} loading={isFetching} line/>
            <CredentialsList credentials={data} error={error}/>
        </MenuWrapper>
    )
}
