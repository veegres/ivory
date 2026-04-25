import {Cancel, CheckCircle, Delete, Edit} from "@mui/icons-material"
import {useEffect, useState} from "react"

import {Feature} from "../../../../api/feature"
import {useRouterVaultDelete, useRouterVaultUpdate} from "../../../../api/vault/hook"
import {Vault} from "../../../../api/vault/type"
import {SimpleButton} from "../../../view/button/SimpleButton"
import {Access} from "../../access/Access"
import {VaultRow} from "./VaultRow"

type Props = {
    uuid: string,
    vault: Vault,
}

export function VaultItem(props: Props) {
    const {uuid, vault: v} = props
    const [edit, setEdit] = useState(false)
    const [empty, setEmpty] = useState(false)
    const [vault, setVault] = useState(v)

    const deleteVault = useRouterVaultDelete()
    const updateVault = useRouterVaultUpdate(v.type, () => setEdit(false))

    useEffect(() => setVault(vault), [vault])

    return (
        <VaultRow
            renderButtons={edit ? renderWriteButtons() : renderReadButtons()}
            disabled={!edit}
            vault={v}
            onChangeVault={(vault) => setVault(vault)}
            onEmpty={(v) => setEmpty(v)}
        />
    )

    function renderReadButtons() {
        return (
            <>
                <Access feature={Feature.ManageVaultUpdate}>
                    <SimpleButton onClick={() => setEdit(true)} disabled={deleteVault.isPending}><Edit/></SimpleButton>
                </Access>
                <Access feature={Feature.ManageVaultDelete}>
                    <SimpleButton  loading={deleteVault.isPending} onClick={handleDelete}><Delete/></SimpleButton>
                </Access>
            </>
        )
    }

    function renderWriteButtons() {
        return (
            <>
                <SimpleButton onClick={() => setEdit(false)} disabled={updateVault.isPending}><Cancel/></SimpleButton>
                <SimpleButton loading={updateVault.isPending} onClick={handleUpdate} disabled={empty}><CheckCircle/></SimpleButton>
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
