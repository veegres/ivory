import {useState} from "react"

import {useRouterVaultCreate} from "../../../../api/vault/hook"
import {Vault, VaultType} from "../../../../api/vault/type"
import {CancelIconButton, SaveIconButton} from "../../../view/button/IconButtons"
import {VaultRow} from "./VaultRow"

export function VaultNew() {
    const initVault: Vault = {username: "", secret: "", type: VaultType.DATABASE_PASSWORD}
    const [vault, setVault] = useState(initVault)
    const [empty, setEmpty] = useState(false)
    const [clean, setClean] = useState(false)
    const createVault = useRouterVaultCreate(handleCancel)

    return (
        <VaultRow
            renderButtons={renderButtons()}
            disabled={false}
            vault={vault}
            error={clean}
            onChangeVault={(vault) => { setVault(vault); setClean(true) }}
            onEmpty={(v) => setEmpty(v)}
        />
    )

    function renderButtons() {
        return (
            <>
                <CancelIconButton size={36} onClick={handleCancel} disabled={!clean || createVault.isPending}/>
                <SaveIconButton size={36} loading={createVault.isPending} onClick={handleCreate} disabled={empty}/>
            </>
        )
    }

    function handleCancel() {
        setVault({...initVault, type: vault.type})
        setClean(false)
    }

    function handleCreate() {
        createVault.mutate(vault)
    }
}
