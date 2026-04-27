import {Cancel, CheckCircle, Delete, Edit} from "@mui/icons-material"
import {Box, Tooltip} from "@mui/material"
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

export function VaultListItem(props: Props) {
    const {uuid, vault: v} = props
    const [edit, setEdit] = useState(false)
    const [empty, setEmpty] = useState(false)
    const [vault, setVault] = useState(v)

    const deleteVault = useRouterVaultDelete(v.type)
    const updateVault = useRouterVaultUpdate(v.type, () => setEdit(false))

    useEffect(() => setVault(vault), [vault])

    return (
        <Tooltip title={`ID: ${uuid}`} placement={"top-start"}>
            <Box>
                <VaultRow
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
