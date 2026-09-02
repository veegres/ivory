import {Box, TextField, ToggleButton, ToggleButtonGroup} from "@mui/material"

import {Hint} from "../../../shared/component/box/Hint"
import {FieldRow} from "../../../shared/component/input/FieldRow"
import {SxPropsMap} from "../../../shared/helper/HelperType"
import {VaultType} from "../../vault/api/VaultType"
import {ClusterOptionsVault} from "./ClusterOptionsVault"

const SX: SxPropsMap = {
    toggleButton: {lineHeight: 1},
}

export type CredentialMode = "new" | "vault" | "none"

export type Credential = {
    username: string,
    password: string,
}

type Props = {
    type: VaultType,
    mode: CredentialMode,
    credential: Credential,
    vaultId?: string,
    locked?: boolean,
    optional?: boolean,
    showErrors: boolean,
    onModeChange: (mode: CredentialMode) => void,
    onCredentialChange: (credential: Credential) => void,
    onVaultChange: (type: VaultType, vaultId?: string) => void,
}

export function ClusterDeployCredentials(props: Props) {
    const {type, mode, credential, vaultId, locked = false, optional = false, showErrors} = props
    const {onModeChange, onCredentialChange, onVaultChange} = props

    return (
        <Box sx={{display: "flex", gap: 0.5}}>
            <Box sx={{flexGrow: 1}}>{renderContent()}</Box>
            {renderActions()}
        </Box>
    )

    function renderContent() {
        switch (mode) {
            case "new":
                return renderNew()
            case "vault":
                return renderVault()
            case "none":
                return renderNone()
        }
    }

    function renderActions() {
        return (
            <ToggleButtonGroup exclusive={true} value={mode} onChange={(_, v) => v && handleModeChange(v)}>
                <ToggleButton sx={SX.toggleButton} value={"vault"}>VAULT</ToggleButton>
                <ToggleButton sx={SX.toggleButton} value={"new"}>NEW</ToggleButton>
                {optional && (
                    <ToggleButton sx={SX.toggleButton} value={"none"} disabled={!!credential.username}>NONE</ToggleButton>
                )}
            </ToggleButtonGroup>
        )
    }

    function renderNone() {
        return <Hint>Cluster will not use credentials</Hint>
    }

    function renderNew() {
        return (
            <FieldRow>
                <TextField
                    fullWidth
                    label={"Username"}
                    value={credential.username}
                    disabled={locked}
                    error={showErrors && !credential.username}
                    onChange={(e) => onCredentialChange({...credential, username: e.target.value})}
                />
                <TextField
                    fullWidth
                    type={"password"}
                    label={"Password"}
                    value={credential.password}
                    error={showErrors && !credential.password}
                    onChange={(e) => onCredentialChange({...credential, password: e.target.value})}
                />
            </FieldRow>
        )
    }

    function renderVault() {
        return (
            <ClusterOptionsVault
                type={type}
                selected={vaultId}
                username={locked ? credential.username : undefined}
                onUpdate={onVaultChange}
                error={showErrors && !vaultId}
            />
        )
    }

    function handleModeChange(next: CredentialMode) {
        if (next !== "vault") onVaultChange(type, undefined)
        onModeChange(next)
    }
}
