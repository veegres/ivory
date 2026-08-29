import {TextField, ToggleButton, ToggleButtonGroup} from "@mui/material"

import {OptionsVault} from "../../../core/widgets/options/OptionsVault"
import {TitledBox} from "../../../shared/component/box/TitledBox"
import {FieldRow} from "../../../shared/component/input/FieldRow"
import {SxPropsMap} from "../../../shared/helper/HelperType"
import {VaultType} from "../../vault/api/VaultType"

const SX: SxPropsMap = {
    toggleButton: {padding: "0px 10px"},
}

// CredentialMode is how one credential is answered: an existing vault entry, or
// a username and password the deploy writes to a new vault before it starts.
export type CredentialMode = "new" | "vault"

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
    // NOTE: prefilled, not editable, and the only vault entries offered are
    // the ones that carry it
    lockedUser?: string,
    showErrors: boolean,
    onModeChange: (mode: CredentialMode) => void,
    onCredentialChange: (credential: Credential) => void,
    onVaultChange: (type: VaultType, vaultId?: string) => void,
}

// ClusterDeployCredentials is one credential answer on the deploy screen: ssh,
// keeper and database each get their own. An engine that is its own keeper is
// asked twice, and pointing both at one vault entry is the user's answer.
export function ClusterDeployCredentials(props: Props) {
    const {title, type, mode, credential, vaultId, lockedUser, showErrors} = props
    const {onModeChange, onCredentialChange, onVaultChange} = props

    return (
        <TitledBox title={title} renderActions={renderActions()} island={true}>
            {mode === "new" ? renderNew() : renderVault()}
        </TitledBox>
    )

    function renderActions() {
        return (
            <ToggleButtonGroup size={"small"} exclusive={true} value={mode} onChange={(_, v) => v && handleModeChange(v)}>
                <ToggleButton sx={SX.toggleButton} value={"new"}>NEW</ToggleButton>
                <ToggleButton sx={SX.toggleButton} value={"vault"}>VAULT</ToggleButton>
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
                    disabled={!!lockedUser}
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
            <OptionsVault
                type={type}
                selected={vaultId}
                onUpdate={onVaultChange}
                username={lockedUser || undefined}
                error={showErrors && !vaultId}
            />
        )
    }

    // NOTE: switching to NEW clears the vault id - leaving one behind would
    // send both a vault and a password for the same credential, which the
    // server refuses as two answers to one question
    function handleModeChange(next: CredentialMode) {
        if (next === "new") onVaultChange(type, undefined)
        onModeChange(next)
    }
}
