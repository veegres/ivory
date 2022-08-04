import {Box, FormHelperText, ToggleButton, Tooltip} from "@mui/material";
import {CredentialsInput} from "./CredentialsInput";
import React, {ReactElement, useEffect, useRef, useState} from "react";
import {CredentialOptions} from "../../app/utils";
import {Credential, CredentialType} from "../../app/types";

const SX = {
    row: { display: "flex", alignItems: "center", gap: "15px", padding: "5px 20px"},
    buttons: { display: "flex" },
    toggle: { height: "32px" },
}

type Props = {
    renderButtons: ReactElement,
    disabled: boolean,
    credential: Credential,
    error: boolean,
    onChangeCredential: (value: Credential) => void
    onEmpty: (value: boolean) => void
}

export function CredentialsRow(props: Props) {
    const { renderButtons, credential, disabled, onChangeCredential, onEmpty, error } = props
    const [username, setUsername] = useState(credential.username)
    const [password, setPassword] = useState(credential.password)
    const [type, toggleType] = useState(credential.type)

    const prevDisabledRef = useRef<boolean>()
    useEffect(handlePropsEffect, [credential])
    useEffect(handleDisableEffect, [disabled, credential.username, credential.password])
    useEffect(handleEmptyEffect, [username, password, onEmpty])

    return (
        <Box sx={SX.row}>
            <Box display={"inline-flex"} flexDirection={"column"}>
                {renderToggle()}
                {/* we need this to have the same indent as in input */}
                <FormHelperText>{" "}</FormHelperText>
            </Box>
            <CredentialsInput
                label={"username"}
                type={"username"}
                value={username}
                disabled={disabled}
                error={error}
                onChange={handleUsername}
            />
            <CredentialsInput
                label={"password"}
                type={disabled ? "string" : "password"}
                value={password}
                disabled={disabled}
                error={error}
                onChange={handlePassword}
            />
            <Box display={"inline-flex"} flexDirection={"column"}>
                <Box sx={SX.buttons}>{renderButtons}</Box>
                {/* we need this to have the same indent as in input */}
                <FormHelperText>{" "}</FormHelperText>
            </Box>
        </Box>
    )

    function renderToggle() {
        const option = CredentialOptions[type]

        return (
            <Tooltip title={option.name} placement={"top"}>
                <Box component={"span"}>
                    <ToggleButton value={type} size={"small"} sx={SX.toggle} disabled={disabled} onClick={handleToggleType}>
                        {option.icon}
                    </ToggleButton>
                </Box>
            </Tooltip>
        )
    }

    function handleUsername(value: string) {
        setUsername(value)
        onChangeCredential({ username: value, password, type })
    }

    function handlePassword(value: string) {
        setPassword(value)
        onChangeCredential({ username, password: value, type })
    }

    function handleToggleType() {
        switch (type) {
            case CredentialType.POSTGRES:
                toggleType(CredentialType.PATRONI)
                break
            case CredentialType.PATRONI:
                toggleType(CredentialType.POSTGRES)
                break
            default:
                toggleType(CredentialType.POSTGRES)
        }
    }

    function handlePropsEffect() {
        setUsername(credential.username)
        setPassword(credential.password)
        toggleType(credential.type)
    }

    function handleDisableEffect() {
        if (!prevDisabledRef.current && disabled) {
            setUsername(credential.username)
            setPassword(credential.password)
            toggleType(credential.type)
        }
        if (prevDisabledRef.current && !disabled) {
            setPassword("")
        }
        prevDisabledRef.current = disabled
    }

    function handleEmptyEffect() {
        if (!username) onEmpty(true)
        else if (!password) onEmpty(true)
        else onEmpty(false)
    }
}
