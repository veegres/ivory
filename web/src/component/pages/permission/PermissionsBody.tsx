import {Typography} from "@mui/material"

import {PermissionMap} from "../../../api/permission/type"
import {AlertCentered} from "../../view/box/AlertCentered"
import {PageStartupBox} from "../../view/box/PageStartupBox"
import {Menu} from "../../widgets/settings/menu/Menu"
import {PermissionsList} from "../../widgets/settings/permissions/PermissionsList"

type Props = {
    username: string,
    permissions: PermissionMap,
}

export function PermissionsBody(props: Props) {
    const {username, permissions} = props
    return (
        <PageStartupBox header={"Permissions"} renderFooter={renderFooter()} position={"start"} padding={"50px 0px"}>
            <Menu/>
            <Typography variant={"h6"}>Glad to see you, {username}!</Typography>
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