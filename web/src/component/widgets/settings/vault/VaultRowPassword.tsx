import {Box} from "@mui/material"
import {ReactElement, useEffect, useRef, useState} from "react"

import {Vault} from "../../../../api/vault/type"
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

export function VaultRowPassword(props: Props) {
    const {renderButtons, vault, disabled, onChangeVault, onEmpty} = props
    const [username, setUsername] = useState(vault.username)
    const [secret, setSecret] = useState(vault.secret)

    const prevDisabledRef = useRef(false)
    useEffect(handleEffectVaultUpdate, [vault])
    useEffect(handleEffectDisabledUpdate, [disabled, vault])
    useEffect(handleEffectEmptyUpdate, [username, secret, onEmpty])

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
                    onChangeVault({username: v, secret, type: vault.type, metadata: vault.metadata})
                }}
            />
            <VaultInput
                sx={{flexGrow: 1}}
                label={"Password"}
                type={"password"}
                value={secret}
                disabled={disabled}
                onChange={(v) => {
                    setSecret(v)
                    onChangeVault({username, secret: v, type: vault.type, metadata: vault.metadata})
                }}
            />
            {renderButtons}
        </Box>
    )

    function handleEffectVaultUpdate() {
        setUsername(vault.username)
        setSecret(vault.secret)
    }

    function handleEffectDisabledUpdate() {
        if (!prevDisabledRef.current && disabled) {
            setUsername(vault.username)
            setSecret(vault.secret)
        }
        if (prevDisabledRef.current && !disabled) {
            setSecret("")
        }
        prevDisabledRef.current = disabled
    }

    function handleEffectEmptyUpdate() {
        if (!username) onEmpty(true)
        else if (!secret) onEmpty(true)
        else onEmpty(false)
    }
}
