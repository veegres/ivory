import {LogoutButton} from "../../../features/auth/component/LogoutButton"
import {AlertCentered} from "../../../shared/component/box/AlertCentered"
import {PageStartupBox} from "../../../shared/component/box/PageStartupBox"
import {PageStartupGreeting} from "../../../shared/component/box/PageStartupGreeting"

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
                <PageStartupGreeting username={username}/>
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