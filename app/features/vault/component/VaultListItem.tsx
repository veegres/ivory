import {Cancel, CheckCircle, ContentCopy, Delete, Edit} from "@mui/icons-material"
import {useEffect, useState} from "react"

import {SimpleButton} from "../../../shared/component/button/SimpleButton"
import {useSnackbar} from "../../../shared/provider/SnackbarProvider"
import {Feature} from "../../Feature"
import {ManageAccess} from "../../management/component/ManageAccess"
import {useRouterVaultDelete, useRouterVaultUpdate} from "../api/VaultHook"
import {Vault, VaultType} from "../api/VaultType"
import {VaultRowPassword} from "./VaultRowPassword"
import {VaultRowSshKey} from "./VaultRowSshKey"

type Props = {
    uuid: string,
    vault: Vault,
}

export function VaultListItem(props: Props) {
    const {uuid, vault: v} = props
    const [edit, setEdit] = useState(false)
    const [empty, setEmpty] = useState(false)
    const [vault, setVault] = useState(v)
    const snackbar = useSnackbar()

    const deleteVault = useRouterVaultDelete(v.type)
    const updateVault = useRouterVaultUpdate(v.type, () => setEdit(false))

    useEffect(handleEffectVaultSync, [v])

    const isSshKey = v.type === VaultType.SSH_KEY
    const Row = isSshKey ? VaultRowSshKey : VaultRowPassword

    return (
        <Row
            uuid={uuid}
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
                <ManageAccess feature={Feature.ManageVaultUpdate}>
                    {isSshKey ? (
                        <SimpleButton onClick={handleCopy}><ContentCopy/></SimpleButton>
                    ) : (
                        <SimpleButton onClick={() => setEdit(true)} disabled={deleteVault.isPending}><Edit/></SimpleButton>
                    )}
                </ManageAccess>
                <ManageAccess feature={Feature.ManageVaultDelete}>
                    <SimpleButton loading={deleteVault.isPending} onClick={handleDelete}><Delete/></SimpleButton>
                </ManageAccess>
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

    function handleCopy() {
        if (v.metadata) {
            navigator.clipboard.writeText(v.metadata).then(() => {
                snackbar("Public key copied to clipboard!", "info")
            })
        }
    }

    function handleEffectVaultSync() {
        setVault(v)
    }
}
