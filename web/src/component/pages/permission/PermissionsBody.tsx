import {Alert} from "@mui/material"

import {PermissionStatus} from "../../../api/permission/type"
import {PageStartupBox} from "../../view/box/PageStartupBox"
import {Menu} from "../../widgets/settings/menu/Menu"
import {PermissionsList} from "../../widgets/settings/permissions/PermissionsList"

type Props = {
    username: string,
    permissions: { [permission: string]: PermissionStatus },
}

export function PermissionsBody(props: Props) {
    const {username, permissions} = props
    return (
        <PageStartupBox header={"Permissions"} renderFooter={renderFooter()} position={"start"}>
            <Menu/>
            <Alert variant={"outlined"} color={"warning"} icon={false}>
                You don't have any permissions yet. You can request them here. Once
                you request permissions, please wait for an authorized person to grant approval.
            </Alert>
        </PageStartupBox>
    )

    function renderFooter() {
        return <PermissionsList username={username} permissions={permissions}/>
    }
}