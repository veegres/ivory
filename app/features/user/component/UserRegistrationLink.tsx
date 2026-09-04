import {Box} from "@mui/material"

import {AlertCentered} from "../../../shared/component/box/AlertCentered"
import {CopyIconButton} from "../../../shared/component/button/IconButtons"
import {SxPropsMap} from "../../../shared/helper/HelperType"
import {DateTimeFormatter} from "../../../shared/helper/HelperUtils"
import {useSnackbar} from "../../../shared/provider/SnackbarProvider"
import {UserRegistration} from "../api/UserType"
import {buildUserRegistrationUrl} from "../api/UserUrl"

const SX: SxPropsMap = {
    box: {display: "flex", flexDirection: "column", gap: 1},
    field: {
        display: "flex", alignItems: "center", gap: 1, height: "40px", padding: "0px 4px 0px 12px",
        border: 1, borderColor: "divider", borderRadius: 1,
    },
    link: {
        flexGrow: 1, minWidth: 0, whiteSpace: "nowrap", overflow: "hidden", textOverflow: "ellipsis",
        fontFamily: "monospace", fontSize: "12px", color: "text.secondary",
    },
}

type Props = {
    registration: UserRegistration,
    reset?: boolean,
}

export function UserRegistrationLink(props: Props) {
    const {registration, reset} = props
    const url = buildUserRegistrationUrl(registration.token)
    const snackbar = useSnackbar()

    return (
        <Box sx={SX.box}>
            <Box sx={SX.field}>
                <Box sx={SX.link} title={url}>{url}</Box>
                <CopyIconButton tooltip={"Copy the link to your clipboard"} onClick={handleCopy}/>
            </Box>
            <AlertCentered severity={"warning"} text={renderWarning()}/>
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

    function handleCopy() {
        navigator.clipboard.writeText(url).then(() => {
            snackbar("Link copied to clipboard!", "info")
        })
    }
}
