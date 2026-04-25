import {Box, Button} from "@mui/material"
import {cloneElement} from "react"

import {Feature} from "../../../../api/feature"
import {Status} from "../../../../api/permission/type"
import {SxPropsMap} from "../../../../app/type"
import {PermissionOptions} from "../../../../app/utils"
import {PermissionsButtons} from "./PermissionsButtons"

const SX: SxPropsMap = {
    item: {
        display: "flex", justifyContent: "space-between", alignItems: "center",
        borderBottom: 1, borderColor: "divider", padding: "4px 8px", height: "35px",
        "&:first-of-type": {borderTop: 1, borderColor: "divider"},
    },
    button: {padding: "2px 5px"},
    wrap: {display: "flex", alignItems: "center", gap: 1},
}

type Props = {
    name: Feature,
    status: Status,
    username: string,
    view?: "admin" | "user",
}

export function PermissionsListItem(props: Props) {
    const {username, name, status, view = "user"} = props
    const options = PermissionOptions[status]

    return (
        <Box sx={SX.item}>
            <Box sx={SX.wrap}>
                {cloneElement(options.icon, {sx: {color: options.color, fontSize: "20px"}})}
                <Box>{name}</Box>
            </Box>
            {renderButton()}
        </Box>
    )

    function renderButton() {
        if (view === "admin") {
            return <PermissionsButtons
                username={username}
                permissions={[[name, status]]}
                approve={status === Status.NOT_PERMITTED || status === Status.PENDING}
                reject={status === Status.GRANTED || status === Status.PENDING}
            />
        }
        if (status === Status.NOT_PERMITTED) {
            return <PermissionsButtons username={username} permissions={[[name, status]]} request={true}/>
        }
        return renderStatus(status)
    }

    function renderStatus(status: Status) {
        return <Button sx={SX.button} disabled={true} size={"small"}>{Status[status]}</Button>
    }


}