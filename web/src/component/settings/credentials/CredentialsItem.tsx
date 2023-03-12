import {CancelIconButton, DeleteIconButton, EditIconButton, SaveIconButton} from "../../view/IconButtons";
import {useEffect, useState} from "react";
import {Password} from "../../../type/password";
import {useMutation} from "@tanstack/react-query";
import {passwordApi} from "../../../app/api";
import {CredentialsRow} from "./CredentialsRow";
import {useMutationOptions} from "../../../hook/QueryCustom";

type Props = {
    uuid: string,
    credential: Password,
}

export function CredentialsItem(props: Props) {
    const {uuid} = props
    const [edit, setEdit] = useState(false)
    const [empty, setEmpty] = useState(false)
    const [credential, setCredential] = useState(props.credential)


    useEffect(() => {
        setCredential(credential)
    }, [credential])

    const deleteOptions = useMutationOptions([["credentials"]])
    const deleteCredentials = useMutation(passwordApi.delete, deleteOptions)
    const updateOptions = useMutationOptions([["credentials"]], () => setEdit(false))
    const updateCredentials = useMutation(passwordApi.update, updateOptions)

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
                <EditIconButton size={36} onClick={() => setEdit(true)} disabled={deleteCredentials.isLoading}/>
                <DeleteIconButton size={36} loading={deleteCredentials.isLoading} onClick={handleDelete}/>
            </>
        )
    }

    function renderWriteButtons() {
        return (
            <>
                <CancelIconButton size={36} onClick={() => setEdit(false)} disabled={updateCredentials.isLoading}/>
                <SaveIconButton size={36} loading={updateCredentials.isLoading} onClick={handleUpdate} disabled={empty}/>
            </>
        )
    }

    function handleDelete() {
        deleteCredentials.mutate(uuid)
    }

    function handleUpdate() {
        updateCredentials.mutate({uuid, credential})
    }
}
