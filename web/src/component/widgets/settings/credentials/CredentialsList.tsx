import {Collapse} from "@mui/material"
import {TransitionGroup} from "react-transition-group"

import {Password, PasswordMap} from "../../../../api/password/type"
import {ErrorSmart} from "../../../view/box/ErrorSmart"
import {NoBox} from "../../../view/box/NoBox"
import {CredentialsItem} from "./CredentialsItem"

type Props = {
    credentials?: PasswordMap,
    error: any,
}

export function CredentialsList(props: Props) {
    const {credentials, error} = props
    if (error) return <ErrorSmart error={error}/>

    const list = Object.entries<Password>(credentials ?? {})
    if (list.length === 0) return <NoBox text={"There are no credentials yet"}/>

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
