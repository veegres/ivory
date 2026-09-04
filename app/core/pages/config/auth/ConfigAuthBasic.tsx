import {Box, ToggleButton, ToggleButtonGroup} from "@mui/material"

import {AlertCentered} from "../../../../shared/component/box/AlertCentered"
import {SxPropsMap} from "../../../../shared/helper/HelperType"

const SX: SxPropsMap = {
    box: {display: "flex", flexDirection: "column", gap: 1},
}

type Props = {
    enabled: boolean,
    onChange: (enabled: boolean) => void,
}

export function ConfigAuthBasic(props: Props) {
    const {enabled, onChange} = props
    return (
        <Box sx={SX.box}>
            <AlertCentered text={renderDescription()}/>
            <ToggleButtonGroup value={enabled} exclusive fullWidth>
                <ToggleButton value={true} onClick={() => onChange(true)}>Enabled</ToggleButton>
                <ToggleButton value={false} onClick={() => onChange(false)}>Disabled</ToggleButton>
            </ToggleButtonGroup>
        </Box>
    )

    function renderDescription() {
        return (
            "Basic authentication signs in the Ivory users. There are no credentials to type here beyond " +
            "the superuser above: everybody after them is invited with a link and sets their own password. " +
            "With this off, that superuser signs in through LDAP or SSO instead."
        )
    }
}
