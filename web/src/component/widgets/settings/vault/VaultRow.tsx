import {Box} from "@mui/material"
import {ReactElement, useEffect, useRef, useState} from "react"

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

export function VaultRow(props: Props) {
    const {renderButtons, vault, disabled, onChangeVault, onEmpty} = props
    const [username, setUsername] = useState(vault.username)
    const [secret, setSecret] = useState(vault.secret)

    const prevDisabledRef = useRef(false)
    useEffect(handleEffectProps, [vault])
    useEffect(handleEffectDisable, [disabled, vault])
    useEffect(handleEffectEmpty, [username, secret, onEmpty])

    return (
        <Box sx={SX.row}>
            <VaultInput
                sx={{width: "150px"}}
                label={"Username"}
                type={"username"}
                value={username}
                disabled={disabled}
                onChange={handleUsername}
            />
            <VaultInput
                sx={{flexGrow: 1}}
                label={vault.type === VaultType.SSH_KEY ? "Key" : "Password"}
                type={"password"}
                value={secret}
                disabled={disabled}
                onChange={handleSecret}
            />
            {renderButtons}
        </Box>
    )

    function handleUsername(v: string) {
        setUsername(v)
        onChangeVault({username: v, secret, type: vault.type})
    }

    function handleSecret(v: string) {
        setSecret(v)
        onChangeVault({username, secret: v, type: vault.type})
    }

    function handleEffectProps() {
        setUsername(vault.username)
        setSecret(vault.secret)
    }

    function handleEffectDisable() {
        if (!prevDisabledRef.current && disabled) {
            setUsername(vault.username)
            setSecret(vault.secret)
        }
        if (prevDisabledRef.current && !disabled) {
            setSecret("")
        }
        prevDisabledRef.current = disabled
    }

    function handleEffectEmpty() {
        if (!username) onEmpty(true)
        else if (!secret) onEmpty(true)
        else onEmpty(false)
    }
}
