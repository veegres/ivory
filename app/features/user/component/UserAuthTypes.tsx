import {Box, ToggleButton, Tooltip} from "@mui/material"

import {SxPropsMap} from "../../../shared/helper/HelperType"
import {UserAuthType} from "../api/UserType"

const SX: SxPropsMap = {
    box: {display: "flex", alignItems: "center", gap: 0.5},
    toggle: {padding: "3px 10px", whiteSpace: "nowrap"},
}

const OPTIONS: {[key in UserAuthType]: {label: string, description: string}} = {
    [UserAuthType.BASIC]: {label: "Basic", description: "Sign in with an Ivory password of their own"},
    [UserAuthType.LDAP]: {label: "LDAP", description: "Sign in through your LDAP directory"},
    [UserAuthType.OIDC]: {label: "OIDC", description: "Sign in through your SSO provider"},
}

type Props = {
    value: UserAuthType[],
    disabled?: boolean,
    reason?: string,
    onChange: (authTypes: UserAuthType[]) => void,
    size?: "small" | "medium" | "large",
}

export function UserAuthTypes(props: Props) {
    const {value, disabled, reason, onChange, size} = props

    return (
        <Box sx={SX.box}>
            {Object.values(UserAuthType).map(renderToggle)}
        </Box>
    )

    function renderToggle(authType: UserAuthType) {
        const option = OPTIONS[authType]
        const selected = value.includes(authType)
        const last = selected && value.length === 1
        const tooltip = last ? "A user needs at least one way to sign in" : (reason ?? option.description)
        return (
            <Tooltip key={authType} title={tooltip} placement={"top"} arrow disableInteractive>
                <Box component={"span"}>
                    <ToggleButton
                        sx={SX.toggle}
                        size={size}
                        value={authType}
                        selected={selected}
                        disabled={disabled || last}
                        onClick={() => handleToggle(authType)}
                    >
                        {option.label}
                    </ToggleButton>
                </Box>
            </Tooltip>
        )
    }

    function handleToggle(authType: UserAuthType) {
        if (value.includes(authType)) onChange(value.filter((it) => it !== authType))
        else onChange([...value, authType])
    }
}
