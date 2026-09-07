import {Cancel, CheckCircle} from "@mui/icons-material"
import {useEffect, useState} from "react"

import {SimpleButton} from "../../../shared/component/button/SimpleButton"
import {useRouterVaultCreate} from "../api/VaultHook"
import {Vault, VaultType} from "../api/VaultType"
import {VaultNewWrapper} from "./VaultNewWrapper"
import {VaultRowPassword} from "./VaultRowPassword"

type Props = {
    type: VaultType,
}

export function VaultNewPassword(props: Props) {
    const {type} = props
    const initVault: Vault = {username: "", secret: "", type}
    const [vault, setVault] = useState(initVault)
    const [empty, setEmpty] = useState(false)
    const [clean, setClean] = useState(false)
    const createVault = useRouterVaultCreate(type, handleCancel)

    useEffect(handleEffectTypeChange, [type])

    return (
        <VaultNewWrapper description={renderDescription()}>
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
                <SimpleButton tooltip={"Cancel"} disabled={!clean || createVault.isPending} onClick={handleCancel}>
                    <Cancel fontSize={"small"}/>
                </SimpleButton>
                <SimpleButton tooltip={"Create"} loading={createVault.isPending} disabled={empty} onClick={handleCreate}>
                    <CheckCircle fontSize={"small"}/>
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
