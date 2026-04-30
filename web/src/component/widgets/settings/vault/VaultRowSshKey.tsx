import {Box} from "@mui/material"
import {ReactElement, useEffect, useState} from "react"

import {Vault, VaultType} from "../../../../api/vault/type"
import {SxPropsMap} from "../../../../app/type"
import {VaultInput} from "./VaultInput"

const SX: SxPropsMap = {
    row: {display: "flex", alignItems: "center", gap: 1, margin: "6px 0px"},
}

type Props = {
    renderButtons: ReactElement,
    disabled: boolean,
    vault: Vault,
    onChangeVault: (value: Vault) => void
    onEmpty: (value: boolean) => void
}

export function VaultRowSshKey(props: Props) {
    const {renderButtons, vault, disabled, onChangeVault, onEmpty} = props
    const [username, setUsername] = useState(vault.username)

    useEffect(handleEffectVaultUpdate, [vault])
    useEffect(handleEffectEmptyUpdate, [username, onEmpty])

    return (
        <Box sx={SX.row}>
            <VaultInput
                sx={{width: "150px"}}
                label={"Username"}
                type={"username"}
                value={username}
                disabled={disabled}
                onChange={(v) => {
                    setUsername(v)
                    onChangeVault({username: v, secret: vault.secret, type: VaultType.SSH_KEY, metadata: vault.metadata})
                }}
            />
            <VaultInput
                sx={{flexGrow: 1}}
                label={disabled ? "Public Key" : "Key"}
                type={"text"}
                value={disabled ? (vault.metadata || "") : (username ? "Key will be generated automatically" : "Enter username to generate key")}
                disabled={true}
                onChange={() => {}}
            />
            {renderButtons}
        </Box>
    )

    function handleEffectVaultUpdate() {
        setUsername(vault.username)
    }

    function handleEffectEmptyUpdate() {
        onEmpty(!username)
    }
}
