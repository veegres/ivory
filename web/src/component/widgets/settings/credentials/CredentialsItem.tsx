import {useEffect, useState} from "react"

import {useRouterPasswordDelete, useRouterPasswordUpdate} from "../../../../api/password/hook"
import {Password} from "../../../../api/password/type"
import {Permission} from "../../../../api/permission/type"
import {CancelIconButton, DeleteIconButton, EditIconButton, SaveIconButton} from "../../../view/button/IconButtons"
import {Access} from "../../access/Access"
import {CredentialsRow} from "./CredentialsRow"

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
                <Access permission={Permission.ManagePasswordUpdate}>
                    <EditIconButton size={36} onClick={() => setEdit(true)} disabled={deleteCredentials.isPending}/>
                </Access>
                <Access permission={Permission.ManagePasswordDelete}>
                    <DeleteIconButton size={36} loading={deleteCredentials.isPending} onClick={handleDelete}/>
                </Access>
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
