import {CredentialsRow} from "./CredentialsRow";
import {useState} from "react";
import {useMutation} from "@tanstack/react-query";
import {passwordApi} from "../../../app/api";
import {Password, PasswordType} from "../../../type/password";
import {CancelIconButton, SaveIconButton} from "../../view/button/IconButtons";
import {useMutationOptions} from "../../../hook/QueryCustom";

export function CredentialsNew() {
    const initCredential: Password = { username: "", password: "", type: PasswordType.POSTGRES }
    const [credential, setCredential] = useState(initCredential)
    const [empty, setEmpty] = useState(false)
    const [clean, setClean] = useState(false)

    const createOptions = useMutationOptions([["credentials"]], handleCancel)
    const createCredentials = useMutation(passwordApi.create, createOptions)

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
                <CancelIconButton size={36} onClick={handleCancel} disabled={!clean || createCredentials.isLoading}/>
                <SaveIconButton size={36} loading={createCredentials.isLoading} onClick={handleCreate} disabled={empty}/>
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
