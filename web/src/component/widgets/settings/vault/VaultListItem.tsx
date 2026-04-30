import {Cancel, CheckCircle, ContentCopy, Delete, Edit} from "@mui/icons-material"
import {Box, Tooltip} from "@mui/material"
import {useEffect, useState} from "react"

import {Feature} from "../../../../api/feature"
import {useRouterVaultDelete, useRouterVaultUpdate} from "../../../../api/vault/hook"
import {Vault, VaultType} from "../../../../api/vault/type"
import {useSnackbar} from "../../../../provider/SnackbarProvider"
import {SimpleButton} from "../../../view/button/SimpleButton"
import {Access} from "../../access/Access"
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
        <Tooltip title={`ID: ${uuid}`} placement={"top-start"}>
            <Box>
                <Row
                    renderButtons={edit ? renderWriteButtons() : renderReadButtons()}
                    disabled={!edit}
                    vault={v}
                    onChangeVault={(vault) => setVault(vault)}
                    onEmpty={(v) => setEmpty(v)}
                />
            </Box>
        </Tooltip>
    )

    function renderReadButtons() {
        return (
            <>
                <Access feature={Feature.ManageVaultUpdate}>
                    {isSshKey ? (
                        <SimpleButton onClick={handleCopy}><ContentCopy/></SimpleButton>
                    ) : (
                        <SimpleButton onClick={() => setEdit(true)} disabled={deleteVault.isPending}><Edit/></SimpleButton>
                    )}
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
