import {TextField, ToggleButton, ToggleButtonGroup} from "@mui/material"

import {TitleBox} from "../../../shared/component/box/TitleBox"
import {FieldRow} from "../../../shared/component/input/FieldRow"
import {SxPropsMap} from "../../../shared/helper/HelperType"
import {VaultType} from "../../vault/api/VaultType"
import {ClusterOptionsVault} from "./ClusterOptionsVault"

const SX: SxPropsMap = {
    toggleButton: {padding: "0px 10px"},
}

// CredentialMode is how one credential is answered: an existing vault entry, a
// username and password the deploy writes to a new vault before it starts, or
// none at all. None exists because whether a deployment has credentials is the
// user's answer rather than the engine's - a template names the account its
// commands create, and a deployment that creates none is answered with nothing.
export type CredentialMode = "new" | "vault" | "none"

export type Credential = {
    username: string,
    password: string,
}

type Props = {
    title: string,
    type: VaultType,
    mode: CredentialMode,
    credential: Credential,
    vaultId?: string,
    // NOTE: NONE is offered on every credential but ssh, which is how Ivory
    // reaches the host at all and so is never an answer the deploy can skip
    optional?: boolean,
    showErrors: boolean,
    onModeChange: (mode: CredentialMode) => void,
    onCredentialChange: (credential: Credential) => void,
    onVaultChange: (type: VaultType, vaultId?: string) => void,
}

// ClusterDeployCredentials is one credential answer on the deploy screen: ssh,
// keeper and database each get their own. An engine that is its own keeper is
// asked twice, and pointing both at one vault entry is the user's answer.
export function ClusterDeployCredentials(props: Props) {
    const {title, type, mode, credential, vaultId, optional = false, showErrors} = props
    const {onModeChange, onCredentialChange, onVaultChange} = props

    return (
        <TitleBox label={title} renderActions={renderActions()} island={true} dense={true} collapsible={false}>
            {renderContent()}
        </TitleBox>
    )

    function renderContent() {
        switch (mode) {
            case "new":
                return renderNew()
            case "vault":
                return renderVault()
        }
    }

    function renderActions() {
        return (
            <ToggleButtonGroup size={"small"} exclusive={true} value={mode} onChange={(_, v) => v && handleModeChange(v)}>
                <ToggleButton sx={SX.toggleButton} value={"new"}>NEW</ToggleButton>
                <ToggleButton sx={SX.toggleButton} value={"vault"}>VAULT</ToggleButton>
                {optional && <ToggleButton sx={SX.toggleButton} value={"none"}>NONE</ToggleButton>}
            </ToggleButtonGroup>
        )
    }

    function renderNew() {
        return (
            <FieldRow>
                <TextField
                    fullWidth
                    size={"small"}
                    label={"Username"}
                    value={credential.username}
                    error={showErrors && !credential.username}
                    onChange={(e) => onCredentialChange({...credential, username: e.target.value})}
                />
                <TextField
                    fullWidth
                    size={"small"}
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
