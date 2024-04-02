import {CredentialsNew} from "./CredentialsNew";
import {MenuWrapper} from "../menu/MenuWrapper";
import {CredentialsList} from "./CredentialsList";
import {LinearProgressStateful} from "../../view/progress/LinearProgressStateful";
import {useRouterPassword} from "../../../router/password";

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
