import {useEffect, useState} from "react"

import {useRouterVaultDelete, useRouterVaultUpdate} from "../../../../api/vault/hook"
import {Vault} from "../../../../api/vault/type"
import {Permission} from "../../../../api/permission/type"
import {CancelIconButton, DeleteIconButton, EditIconButton, SaveIconButton} from "../../../view/button/IconButtons"
import {Access} from "../../access/Access"
import {VaultRow} from "./VaultRow"

type Props = {
    uuid: string,
    vault: Vault,
}

export function VaultItem(props: Props) {
    const {uuid} = props
    const [edit, setEdit] = useState(false)
    const [empty, setEmpty] = useState(false)
    const [vault, setVault] = useState(props.vault)

    const deleteVault = useRouterVaultDelete()
    const updateVault = useRouterVaultUpdate(() => setEdit(false))

    useEffect(() => {setVault(vault)}, [vault])

    return (
        <VaultRow
            renderButtons={edit ? renderWriteButtons() : renderReadButtons()}
            disabled={!edit}
            error={edit}
            vault={props.vault}
            onChangeVault={(vault) => setVault(vault)}
            onEmpty={(v) => setEmpty(v)}
        />
    )

    function renderReadButtons() {
        return (
            <>
                <Access permission={Permission.ManageVaultUpdate}>
                    <EditIconButton size={36} onClick={() => setEdit(true)} disabled={deleteVault.isPending}/>
                </Access>
                <Access permission={Permission.ManageVaultDelete}>
                    <DeleteIconButton size={36} loading={deleteVault.isPending} onClick={handleDelete}/>
                </Access>
            </>
        )
    }

    function renderWriteButtons() {
        return (
            <>
                <CancelIconButton size={36} onClick={() => setEdit(false)} disabled={updateVault.isPending}/>
                <SaveIconButton size={36} loading={updateVault.isPending} onClick={handleUpdate} disabled={empty}/>
            </>
        )
    }

    function handleDelete() {
        deleteVault.mutate(uuid)
    }

    function handleUpdate() {
        updateVault.mutate({uuid, vault})
    }
}
