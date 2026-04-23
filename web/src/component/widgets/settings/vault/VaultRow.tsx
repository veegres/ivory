import {Box, FormHelperText, ToggleButton, Tooltip} from "@mui/material"
import {cloneElement, ReactElement, useEffect, useRef, useState} from "react"

import {Vault, VaultType} from "../../../../api/vault/type"
import {SxPropsMap} from "../../../../app/type"
import {VaultOptions} from "../../../../app/utils"
import {VaultInput} from "./VaultInput"

const SX: SxPropsMap = {
    row: {display: "flex", alignItems: "center", gap: "15px", margin: "5px 10px 0px"},
    buttons: {display: "flex"},
    toggle: {height: "36px"},
}

type Props = {
    renderButtons: ReactElement,
    disabled: boolean,
    vault: Vault,
    error: boolean,
    onChangeVault: (value: Vault) => void
    onEmpty: (value: boolean) => void
}

export function VaultRow(props: Props) {
    const {renderButtons, vault, disabled, onChangeVault, onEmpty, error} = props
    const [username, setUsername] = useState(vault.username)
    const [secret, setSecret] = useState(vault.secret)
    const [type, toggleType] = useState(vault.type)

    const prevDisabledRef = useRef(false)
    useEffect(handleEffectProps, [vault])
    useEffect(handleEffectDisable, [disabled, vault])
    useEffect(handleEffectEmpty, [username, secret, onEmpty])

    const isSsh = type === VaultType.SSH_PASSWORD || type === VaultType.SSH_KEY

    return (
        <Box sx={SX.row}>
            <Box display={"inline-flex"} flexDirection={"column"}>
                {renderToggle()}
                {/* we need this to have the same indent as in input */}
                <FormHelperText>{" "}</FormHelperText>
            </Box>
            <VaultInput
                label={"username"}
                type={"username"}
                value={username}
                disabled={disabled}
                error={error}
                onChange={handleUsername}
            />
            <VaultInput
                label={isSsh && type === VaultType.SSH_KEY ? "key" : "password"}
                type={disabled ? "string" : (type === VaultType.SSH_KEY ? "string" : "password")}
                value={secret}
                disabled={disabled}
                error={error}
                onChange={handleSecret}
            />
            <Box display={"inline-flex"} flexDirection={"column"}>
                <Box sx={SX.buttons}>{renderButtons}</Box>
                {/* we need this to have the same indent as in input */}
                <FormHelperText>{" "}</FormHelperText>
            </Box>
        </Box>
    )

    function renderToggle() {
        const option = VaultOptions[type]

        return (
            <Tooltip title={option.name} placement={"top"} disableInteractive>
                <Box component={"span"}>
                    <ToggleButton sx={SX.toggle} size={"small"} value={type} disabled={disabled} onClick={handleToggleType}>
                        {cloneElement(option.icon, {sx: {color: option.color}})}
                    </ToggleButton>
                </Box>
            </Tooltip>
        )
    }

    function handleUsername(v: string) {
        setUsername(v)
        onChangeVault({username: v, secret, type})
    }

    function handleSecret(v: string) {
        setSecret(v)
        onChangeVault({username, secret: v, type})
    }

    function handleToggleType() {
        let tmpType: VaultType
        switch (type) {
            case VaultType.DATABASE_PASSWORD:
                tmpType = VaultType.KEEPER_PASSWORD
                break
            case VaultType.KEEPER_PASSWORD:
                tmpType = VaultType.SSH_PASSWORD
                break
            case VaultType.SSH_PASSWORD:
                tmpType = VaultType.SSH_KEY
                break
            case VaultType.SSH_KEY:
                tmpType = VaultType.DATABASE_PASSWORD
                break
            default:
                tmpType = VaultType.DATABASE_PASSWORD
        }
        toggleType(tmpType)
        onChangeVault({username, secret, type: tmpType})
    }

    function handleEffectProps() {
        setUsername(vault.username)
        setSecret(vault.secret)
        toggleType(vault.type)
    }

    function handleEffectDisable() {
        if (!prevDisabledRef.current && disabled) {
            setUsername(vault.username)
            setSecret(vault.secret)
            toggleType(vault.type)
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
