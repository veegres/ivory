import {PermissionMap} from "../../../features/permission/api/PermissionType"
import {PermissionsList} from "../../../features/permission/component/PermissionsList"
import {AlertCentered} from "../../../shared/component/box/AlertCentered"
import {PageStartupBox} from "../../../shared/component/box/PageStartupBox"
import {PageStartupGreeting} from "../../../shared/component/box/PageStartupGreeting"

type Props = {
    username: string,
    permissions: PermissionMap,
}

export function PermissionsBody(props: Props) {
    const {username, permissions} = props
    return (
        <PageStartupBox header={"Permissions"} renderFooter={renderFooter()} position={"start"} padding={"50px 0px"}>
            <PageStartupGreeting username={username}/>
            <AlertCentered
                severity={"warning"}
                text={`
                You don't have any permissions yet. You can request them here. Once
                you request permissions, please wait for an authorized person to grant approval.
                `}
            />
        </PageStartupBox>
    )

    function renderFooter() {
        return <PermissionsList username={username} permissions={permissions}/>
    }
}