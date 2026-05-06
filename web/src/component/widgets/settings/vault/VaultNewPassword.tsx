import {Cancel, CheckCircle} from "@mui/icons-material"
import {useEffect, useState} from "react"

import {useRouterVaultCreate} from "../../../../api/vault/hook"
import {Vault, VaultType} from "../../../../api/vault/type"
import {VaultOptions} from "../../../../app/utils"
import {SimpleButton} from "../../../view/button/SimpleButton"
import {VaultNewWrapper} from "./VaultNewWrapper"
import {VaultRowPassword} from "./VaultRowPassword"

type Props = {
    type: VaultType,
}

export function VaultNewPassword(props: Props) {
    const {type} = props
    const initVault: Vault = {username: "", secret: "", type}
    const options = VaultOptions[type]
    const [vault, setVault] = useState(initVault)
    const [empty, setEmpty] = useState(false)
    const [clean, setClean] = useState(false)
    const createVault = useRouterVaultCreate(type, handleCancel)

    useEffect(handleEffectTypeChange, [type])

    return (
        <VaultNewWrapper
            icon={options.icon}
            label={options.label}
            description={renderDescription()}
        >
            <VaultRowPassword
                renderButtons={renderButtons()}
                disabled={false}
                vault={vault}
                onChangeVault={(vault) => {setVault(vault); setClean(true)}}
                onEmpty={(v) => setEmpty(v)}
            />
        </VaultNewWrapper>
    )

    function renderDescription() {
        return (
            "All passwords are encrypted by the secret key and safely stored inside Ivory. Ivory decrypt them " +
            "at the moment they're needed."
        )
    }

    function renderButtons() {
        return (
            <>
                <SimpleButton disabled={!clean || createVault.isPending} onClick={handleCancel}><Cancel/></SimpleButton>
                <SimpleButton loading={createVault.isPending} disabled={empty} onClick={handleCreate}>
                    <CheckCircle/>
                </SimpleButton>
            </>
        )
    }

    function handleCancel() {
        setVault(initVault)
        setClean(false)
    }

    function handleCreate() {
        createVault.mutate(vault)
    }

    function handleEffectTypeChange() {
        setVault(v => ({...v, type}))
    }
}
