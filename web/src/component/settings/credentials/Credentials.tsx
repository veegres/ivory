import {CredentialsNew} from "./CredentialsNew";
import {MenuWrapper} from "../menu/MenuWrapper";
import {CredentialsList} from "./CredentialsList";

export function Credentials() {
    return (
        <MenuWrapper>
            <CredentialsNew/>
            <CredentialsList/>
        </MenuWrapper>
    )
}
