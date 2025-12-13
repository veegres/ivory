import {Typography} from "@mui/material"

import {AlertCentered} from "../../view/box/AlertCentered"
import {PageStartupBox} from "../../view/box/PageStartupBox"
import {LogoutButton} from "../../widgets/actions/LogoutButton"

type Props = {
    username?: string,
}

export function PermissionsLogout(props: Props) {
    const {username} = props
    return (
        <PageStartupBox header={"Permissions"} renderFooter={renderFooter()} position={"start"}>
            <Typography variant={"h6"}>Glad to see you, {username}!</Typography>
            <AlertCentered
                severity={"error"}
                text={"Something goes wrong, there are no permissions for you. Please, try to logout and login again."}
            />
        </PageStartupBox>
    )

    function renderFooter() {
        return <LogoutButton/>
    }
}