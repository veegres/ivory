import {Collapse} from "@mui/material";
import {TransitionGroup} from "react-transition-group";

import {Password, PasswordMap} from "../../../../api/password/type";
import {ErrorSmart} from "../../../view/box/ErrorSmart";
import {InfoAlert} from "../../../view/box/InfoAlert";
import {CredentialsItem} from "./CredentialsItem";

type Props = {
    credentials?: PasswordMap,
    error: any,
}

export function CredentialsList(props: Props) {
    const {credentials, error} = props
    if (error) return <ErrorSmart error={error}/>

    const list = Object.entries<Password>(credentials ?? {})
    if (list.length === 0) return <InfoAlert text={"There are no credentials yet"}/>

    return (
        <TransitionGroup appear={false}>
            {list.map(([key, credential]) => (
                <Collapse key={key}>
                    <CredentialsItem uuid={key} credential={credential}/>
                </Collapse>
            ))}
        </TransitionGroup>
    )
}
