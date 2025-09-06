import {CredentialsRow} from "./CredentialsRow";
import {useState} from "react";
import {Password, PasswordType} from "../../../../type/password";
import {CancelIconButton, SaveIconButton} from "../../../view/button/IconButtons";
import {useRouterPasswordCreate} from "../../../../router/password";

export function CredentialsNew() {
    const initCredential: Password = { username: "", password: "", type: PasswordType.POSTGRES }
    const [credential, setCredential] = useState(initCredential)
    const [empty, setEmpty] = useState(false)
    const [clean, setClean] = useState(false)
    const createCredentials = useRouterPasswordCreate(handleCancel)

    return (
        <CredentialsRow
            renderButtons={renderButtons()}
            disabled={false}
            credential={credential}
            error={clean}
            onChangeCredential={(credential) => { setCredential(credential); setClean(true) }}
            onEmpty={(v) => setEmpty(v)}
        />
    )

    function renderButtons() {
        return (
            <>
                <CancelIconButton size={36} onClick={handleCancel} disabled={!clean || createCredentials.isPending}/>
                <SaveIconButton size={36} loading={createCredentials.isPending} onClick={handleCreate} disabled={empty}/>
            </>
        )
    }

    function handleCancel() {
        setCredential({ ...initCredential, type: credential.type })
        setClean(false)
    }

    function handleCreate() {
        createCredentials.mutate(credential)
    }
}
