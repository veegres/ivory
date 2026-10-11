import {Box} from "@mui/material"

import {AlertAlign} from "../../../shared/component/box/AlertAlign"
import {CopyField} from "../../../shared/component/box/CopyField"
import {SxPropsMap} from "../../../shared/helper/HelperType"
import {buildRouteUrl} from "../../../shared/helper/HelperUrl"
import {DateTimeFormatter} from "../../../shared/helper/HelperUtils"
import {UserRegistration} from "../api/UserType"

const SX: SxPropsMap = {
    box: {display: "flex", flexDirection: "column", gap: 0.5},
}

type Props = {
    registration: UserRegistration,
    reset?: boolean,
}

export function UserRegistrationLink(props: Props) {
    const {registration, reset} = props
    const url = buildRouteUrl("/user", registration.token)

    return (
        <Box sx={SX.box}>
            <CopyField
                value={url}
                tooltip={"Copy the link to your clipboard"}
                message={"Link copied to clipboard!"}
            />
            <AlertAlign align={"justify"} text={renderWarning()}/>
        </Box>
    )

    function renderWarning() {
        return (<>
            Anyone holding this link can set the password
            {reset ? " of the existing user " : " for "}<b>{registration.username}</b>, so hand it to that
            person only. It works <b>once</b> and expires
            on <b>{DateTimeFormatter.utc(registration.expiresAt)}</b>. Ivory shows it once as well,
            so copy it now - to hand out another one, issue a new one.
        </>)
    }
}
