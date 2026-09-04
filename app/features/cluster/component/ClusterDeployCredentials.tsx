import {AddTwoTone, BlockTwoTone, LockTwoTone} from "@mui/icons-material"
import {Box, TextField, ToggleButton, ToggleButtonGroup, Tooltip} from "@mui/material"
import {ReactElement, ReactNode} from "react"

import {FieldLabel} from "../../../shared/component/input/FieldLabel"
import {FieldRow} from "../../../shared/component/input/FieldRow"
import {SxPropsMap} from "../../../shared/helper/HelperType"
import {VaultType} from "../../vault/api/VaultType"
import {ClusterOptionsVault} from "./ClusterOptionsVault"

const SX: SxPropsMap = {
    label: {gridColumn: "1", justifyContent: "flex-start"},
    mode: {gridColumn: "2", justifySelf: {xs: "end", sm: "start"}},
    button: {padding: "3px 6px"},
    icon: {fontSize: "20px"},
    content: {gridColumn: {xs: "1 / -1", sm: "3"}, minWidth: 0},
    none: {fontSize: 12, lineHeight: 1, color: "text.secondary"},
}

export type CredentialMode = "new" | "vault" | "none"

export type Credential = {
    username: string,
    password: string,
}

const ModeOptions: { [key in CredentialMode]: {tooltip: string, icon: ReactElement} } = {
    vault: {tooltip: "Use an existing vault entry", icon: <LockTwoTone sx={SX.icon}/>},
    new: {tooltip: "Enter a new username and password", icon: <AddTwoTone sx={SX.icon}/>},
    none: {tooltip: "Not used by this cluster", icon: <BlockTwoTone sx={SX.icon}/>},
}

type Props = {
    type: VaultType,
    label: string,
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

// ClusterDeployCredentials is one answer to one credential question: which
// account, and where it comes from. It renders its three cells straight into
// the caller's grid, so the name, the source and the fields of every credential
// line up in the same three columns whatever source each one is on.
export function ClusterDeployCredentials(props: Props) {
    const {type, label, mode, credential, vaultId, locked = false, optional = false, showErrors} = props
    const {onModeChange, onCredentialChange, onVaultChange} = props

    return (
        <>
            {renderLabel()}
            {renderMode()}
            {renderContent()}
        </>
    )

    function renderLabel() {
        return <FieldLabel sx={SX.label}>{label}</FieldLabel>
    }

    function renderMode() {
        return (
            <ToggleButtonGroup sx={SX.mode} exclusive={true} value={mode} onChange={(_, v) => v && handleModeChange(v)}>
                {getModes().map(renderModeButton)}
            </ToggleButtonGroup>
        )
    }

    function renderModeButton(value: CredentialMode) {
        const {tooltip, icon} = ModeOptions[value]
        return (
            <ToggleButton key={value} sx={SX.button} value={value}>
                <Tooltip title={tooltip} placement={"top"}>{icon}</Tooltip>
            </ToggleButton>
        )
    }

    function renderContent() {
        return <Box sx={SX.content}>{getContent()}</Box>
    }

    function renderNew() {
        return (
            <FieldRow>
                <TextField
                    label={"Username"}
                    value={credential.username}
                    disabled={locked}
                    error={showErrors && !credential.username}
                    onChange={(e) => onCredentialChange({...credential, username: e.target.value})}
                />
                <TextField
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
                label={"Vault Entry"}
                selected={vaultId}
                username={locked ? credential.username : undefined}
                onUpdate={onVaultChange}
                error={showErrors && !vaultId}
            />
        )
    }

    function renderNone() {
        return <FieldLabel sx={SX.none}>Not used by this cluster</FieldLabel>
    }

    function handleModeChange(next: CredentialMode) {
        if (next !== "vault") onVaultChange(type, undefined)
        onModeChange(next)
    }

    function getContent(): ReactNode {
        switch (mode) {
            case "new":
                return renderNew()
            case "vault":
                return renderVault()
            case "none":
                return renderNone()
        }
    }

    function getModes(): CredentialMode[] {
        return optional ? ["vault", "new", "none"] : ["vault", "new"]
    }
}
