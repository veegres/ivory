import {CancelIconButton, DeleteIconButton, EditIconButton, SaveIconButton} from "../../../view/button/IconButtons";
import {useEffect, useState} from "react";
import {Password} from "../../../../type/password";
import {CredentialsRow} from "./CredentialsRow";
import {useRouterPasswordDelete, useRouterPasswordUpdate} from "../../../../router/password";

type Props = {
    uuid: string,
    credential: Password,
}

export function CredentialsItem(props: Props) {
    const {uuid} = props
    const [edit, setEdit] = useState(false)
    const [empty, setEmpty] = useState(false)
    const [credential, setCredential] = useState(props.credential)

    const deleteCredentials = useRouterPasswordDelete()
    const updateCredentials = useRouterPasswordUpdate(() => setEdit(false))

    useEffect(() => {setCredential(credential)}, [credential])

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
                <EditIconButton size={36} onClick={() => setEdit(true)} disabled={deleteCredentials.isPending}/>
                <DeleteIconButton size={36} loading={deleteCredentials.isPending} onClick={handleDelete}/>
            </>
        )
    }

    function renderWriteButtons() {
        return (
            <>
                <CancelIconButton size={36} onClick={() => setEdit(false)} disabled={updateCredentials.isPending}/>
                <SaveIconButton size={36} loading={updateCredentials.isPending} onClick={handleUpdate} disabled={empty}/>
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
