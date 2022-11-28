import {CredentialsRow} from "./CredentialsRow";
import React, {useState} from "react";
import {useMutation, useQueryClient} from "@tanstack/react-query";
import {credentialApi} from "../../app/api";
import {Credential, CredentialMap, CredentialType} from "../../app/types";
import {CancelIconButton, SaveIconButton} from "../view/IconButtons";
import {useMutationOptions} from "../../app/hooks";

export function CredentialsNew() {
    const initCredential: Credential = { username: "", password: "", type: CredentialType.POSTGRES }
    const [credential, setCredential] = useState(initCredential)
    const [empty, setEmpty] = useState(false)
    const [clean, setClean] = useState(false)
    const { onError } = useMutationOptions()

    const queryClient = useQueryClient();
    const createCredentials = useMutation(credentialApi.create, {
        onSuccess: ({key, credential}) => {
            const map = queryClient.getQueryData<CredentialMap>(["credentials"]) ?? {}
            map[key] = credential
            queryClient.setQueryData<CredentialMap>(["credentials"], map)
            handleCancel()
        },
        onError,
    })

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
                <CancelIconButton loading={false} onClick={handleCancel} disabled={!clean || createCredentials.isLoading}/>
                <SaveIconButton loading={createCredentials.isLoading} onClick={handleCreate} disabled={empty}/>
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
