import {useRouterPassword} from "../../../../api/password/hook"
import {LinearProgressStateful} from "../../../view/progress/LinearProgressStateful"
import {MenuWrapper} from "../menu/MenuWrapper"
import {CredentialsList} from "./CredentialsList"
import {CredentialsNew} from "./CredentialsNew"

export function Credentials() {
    const query = useRouterPassword()
    const {data, error, isFetching} = query

    return (
        <MenuWrapper>
            <CredentialsNew/>
            <LinearProgressStateful color={"inherit"} loading={isFetching} line/>
            <CredentialsList credentials={data} error={error}/>
        </MenuWrapper>
    )
}
