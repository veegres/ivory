import {Box, FormHelperText, ToggleButton, Tooltip} from "@mui/material";
import {CredentialsInput} from "./CredentialsInput";
import {cloneElement, ReactElement, useEffect, useRef, useState} from "react";
import {CredentialOptions} from "../../../../app/utils";
import {Password, PasswordType} from "../../../../api/password/type";
import {SxPropsMap} from "../../../../api/management/type";

const SX: SxPropsMap = {
    row: {display: "flex", alignItems: "center", gap: "15px", margin: "5px 10px 0px"},
    buttons: {display: "flex"},
    toggle: {height: "36px"},
}

type Props = {
    renderButtons: ReactElement,
    disabled: boolean,
    credential: Password,
    error: boolean,
    onChangeCredential: (value: Password) => void
    onEmpty: (value: boolean) => void
}

export function CredentialsRow(props: Props) {
    const {renderButtons, credential, disabled, onChangeCredential, onEmpty, error} = props
    const [username, setUsername] = useState(credential.username)
    const [password, setPassword] = useState(credential.password)
    const [type, toggleType] = useState(credential.type)

    const prevDisabledRef = useRef(false)
    useEffect(handleEffectProps, [credential])
    useEffect(handleEffectDisable, [disabled, credential])
    useEffect(handleEffectEmpty, [username, password, onEmpty])

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
            <Tooltip title={option.name} placement={"top"} disableInteractive>
                <Box component={"span"}>
                    <ToggleButton sx={SX.toggle} size={"small"} value={type} disabled={disabled} onClick={handleToggleType}>
                        {cloneElement(option.icon, {sx: {color: option.color}})}
                    </ToggleButton>
                </Box>
            </Tooltip>
        )
    }

    function handleUsername(value: string) {
        setUsername(value)
        onChangeCredential({username: value, password, type})
    }

    function handlePassword(value: string) {
        setPassword(value)
        onChangeCredential({username, password: value, type})
    }

    function handleToggleType() {
        let tmpType: PasswordType
        switch (type) {
            case PasswordType.POSTGRES:
                tmpType = PasswordType.PATRONI
                break
            case PasswordType.PATRONI:
                tmpType = PasswordType.POSTGRES
                break
        }
        toggleType(tmpType)
        onChangeCredential({username, password, type: tmpType})
    }

    function handleEffectProps() {
        setUsername(credential.username)
        setPassword(credential.password)
        toggleType(credential.type)
    }

    function handleEffectDisable() {
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

    function handleEffectEmpty() {
        if (!username) onEmpty(true)
        else if (!password) onEmpty(true)
        else onEmpty(false)
    }
}
