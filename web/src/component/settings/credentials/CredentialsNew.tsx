import {CredentialsRow} from "./CredentialsRow";
import {useState} from "react";
import {useMutation} from "@tanstack/react-query";
import {credentialApi} from "../../../app/api";
import {Credential, CredentialType} from "../../../app/types";
import {CancelIconButton, SaveIconButton} from "../../view/IconButtons";
import {useMutationOptions} from "../../../hook/QueryCustom";

export function CredentialsNew() {
    const initCredential: Credential = { username: "", password: "", type: CredentialType.POSTGRES }
    const [credential, setCredential] = useState(initCredential)
    const [empty, setEmpty] = useState(false)
    const [clean, setClean] = useState(false)

    const createOptions = useMutationOptions([["credentials"]], handleCancel)
    const createCredentials = useMutation(credentialApi.create, createOptions)

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
