import {Box} from "@mui/material"

import {PermissionMap} from "../../../features/permission/api/PermissionType"
import {PermissionList} from "../../../features/permission/component/PermissionList"
import {AlertCentered} from "../../../shared/component/box/AlertCentered"
import {PageStartupBox} from "../../../shared/component/box/PageStartupBox"
import {PageStartupGreeting} from "../../../shared/component/box/PageStartupGreeting"
import {SxPropsMap} from "../../../shared/helper/HelperType"

const SX: SxPropsMap = {
    footer: {maxHeight: "550px", width: "100%", overflow: "auto"},
}


type Props = {
    username: string,
    permissions: PermissionMap,
}

export function PermissionsBody(props: Props) {
    const {username, permissions} = props
    return (
        <PageStartupBox header={"Permissions"} renderFooter={renderFooter()} position={"start"} padding={"15px 0px"}>
            <PageStartupGreeting username={username}/>
            <AlertCentered
                severity={"warning"}
                text={`
                You don't have any permissions yet. Request them here, then wait for an authorised
                 person to approve your request.
                `}
            />
        </PageStartupBox>
    )

    function renderFooter() {
        return (
            <Box sx={SX.footer}>
                <PermissionList username={username} permissions={permissions}/>
            </Box>
        )
    }
}