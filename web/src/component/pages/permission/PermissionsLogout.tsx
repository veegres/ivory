import {Typography} from "@mui/material"

import {AlertCentered} from "../../view/box/AlertCentered"
import {PageStartupBox} from "../../view/box/PageStartupBox"
import {LogoutButton} from "../../widgets/actions/LogoutButton"

type Props = {
    username?: string,
    error?: string,
}

export function PermissionsLogout(props: Props) {
    const {username, error} = props
    return (
        <PageStartupBox header={"Permissions"} renderFooter={renderFooter()} position={"start"}>
            {!username ? renderUsernameProblem() : renderPermissionProblem()}
        </PageStartupBox>
    )

    function renderPermissionProblem() {
        return (
            <>
                <Typography variant={"h6"}>Glad to see you, {username}!</Typography>
                <AlertCentered
                    severity={"error"}
                    text={`Something went wrong, there are no permissions. Please, try to logout and login again. (${error ?? "unknown error"})`}
                />
            </>
        )
    }

    function renderUsernameProblem() {
        return (
            <AlertCentered
                severity={"error"}
                text={`Something went wrong, there is no user information. Please, try to logout and login again. (${error ?? "unknown error"})`}
            />
        )
    }

    function renderFooter() {
        return <LogoutButton/>
    }
}