import {CancelIconButton, DeleteIconButton, EditIconButton, SaveIconButton} from "../view/IconButtons";
import React, {useEffect, useState} from "react";
import {Credential} from "../../app/types";
import {useMutation} from "@tanstack/react-query";
import {credentialApi} from "../../app/api";
import {CredentialsRow} from "./CredentialsRow";
import {useMutationOptions} from "../../app/hooks";

type Props = {
    uuid: string,
    credential: Credential,
}

export function CredentialsItem(props: Props) {
    const { uuid } = props
    const [edit, setEdit] = useState(false)
    const [empty, setEmpty] = useState(false)
    const [credential, setCredential] = useState(props.credential)


    useEffect(() => { setCredential(credential) }, [credential])

    const deleteOptions = useMutationOptions(["credentials"])
    const deleteCredentials = useMutation(credentialApi.delete, deleteOptions)
    const updateOptions = useMutationOptions(["credentials"], () => setEdit(false))
    const updateCredentials = useMutation(credentialApi.update, updateOptions)

    return (
        <CredentialsRow
            renderButtons={edit ? renderWriteButtons() : renderReadButtons()}
            disabled={!edit}
            error={edit}
            credential={props.credential}
            onChangeCredential={(credential) => setCredential(credential)}
            onEmpty={(v) => setEmpty(v)}
        />
    )

    function renderReadButtons() {
        return (
            <>
                <EditIconButton loading={false} onClick={() => setEdit(true)} disabled={deleteCredentials.isLoading}/>
                <DeleteIconButton loading={deleteCredentials.isLoading} onClick={handleDelete}/>
            </>
        )
    }

    function renderWriteButtons() {
        return (
            <>
                <CancelIconButton loading={false} onClick={() => setEdit(false)} disabled={updateCredentials.isLoading}/>
                <SaveIconButton loading={updateCredentials.isLoading} onClick={handleUpdate} disabled={empty}/>
            </>
        )
    }

    function handleDelete() {
        deleteCredentials.mutate(uuid)
    }

    function handleUpdate() {
        updateCredentials.mutate({ uuid, credential })
    }
}
