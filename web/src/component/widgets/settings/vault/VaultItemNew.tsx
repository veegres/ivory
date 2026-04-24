import {Cancel, CheckCircle} from "@mui/icons-material"
import {Box} from "@mui/material"
import {useEffect, useState} from "react"

import {useRouterVaultCreate} from "../../../../api/vault/hook"
import {Vault, VaultType} from "../../../../api/vault/type"
import {SxPropsMap} from "../../../../app/type"
import {VaultOptions} from "../../../../app/utils"
import {SimpleButton} from "../../../view/button/SimpleButton"
import {VaultRow} from "./VaultRow"

const SX: SxPropsMap = {
    box: {paddingRight: "10px", display: "flex", flexDirection: "column", gap: 1},
    head: {display: "flex", alignItems: "center", padding: "8px 5px", color: "text.secondary", fontSize: 17},
}

type Props = {
    type: VaultType,
}

export function VaultItemNew(props: Props) {
    const {type} = props
    const initVault: Vault = {username: "", secret: "", type}
    const options = VaultOptions[type]
    const [vault, setVault] = useState(initVault)
    const [empty, setEmpty] = useState(false)
    const [clean, setClean] = useState(false)
    const createVault = useRouterVaultCreate(type, handleCancel)

    useEffect(() => setVault(v => ({...v, type})), [type])

    return (
        <Box sx={SX.box}>
            <Box sx={SX.head}>{options.icon}&nbsp;New {options.label}</Box>
            <VaultRow
                renderButtons={renderButtons()}
                disabled={false}
                vault={vault}
                onChangeVault={(vault) => {setVault(vault); setClean(true)}}
                onEmpty={(v) => setEmpty(v)}
            />
        </Box>
    )

    function renderButtons() {
        return (
            <>
                <SimpleButton disabled={!clean || createVault.isPending} onClick={handleCancel}><Cancel/></SimpleButton>
                <SimpleButton loading={createVault.isPending} disabled={empty} onClick={handleCreate}><CheckCircle/></SimpleButton>
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
}
