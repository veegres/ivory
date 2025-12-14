import {Box, Button} from "@mui/material"
import {cloneElement} from "react"

import {
    useRouterPermissionApprove,
    useRouterPermissionReject,
    useRouterPermissionRequest
} from "../../../../api/permission/hook"
import {PermissionStatus} from "../../../../api/permission/type"
import {SxPropsMap} from "../../../../app/type"
import {PermissionOptions} from "../../../../app/utils"

const SX: SxPropsMap = {
    item: {
        display: "flex", justifyContent: "space-between", alignItems: "center",
        borderBottom: 1, borderColor: "divider", padding: "3px 8px",
        "&:first-of-type": {borderTop: 1, borderColor: "divider"},
    },
    wrap: {display: "flex", alignItems: "center", gap: 1, height: "28px"},
    button: {padding: "2px 5px"},
}

type Props = {
    name: string,
    status: PermissionStatus,
    username: string,
    view?: "admin" | "user",
}

export function PermissionsListItem(props: Props) {
    const {username, name, status, view = "user"} = props
    const options = PermissionOptions[status]
    const request = useRouterPermissionRequest()
    const approve = useRouterPermissionApprove()
    const reject = useRouterPermissionReject()

    return (
        <Box sx={SX.item}>
            <Box sx={SX.wrap}>
                {cloneElement(options.icon, {sx: {color: options.color, fontSize: "20px"}})}
                <Box>{name}</Box>
            </Box>
            <Box sx={SX.wrap}>
                {renderButton()}
            </Box>
        </Box>
    )

    function renderButton() {
        if (view === "admin") {
            if (status === PermissionStatus.GRANTED) {
                return renderReject("Revoke")
            }
            if (status === PermissionStatus.NOT_PERMITTED) {
                return renderApprove("Grant")
            }
            if (status === PermissionStatus.PENDING) {
                return (
                    <>
                        {renderApprove("Approve")}
                        {renderReject("Reject")}
                    </>
                )
            }
        } else {
            if (status === PermissionStatus.NOT_PERMITTED) {
                return renderRequest()
            }
        }
        return renderStatus(status)
    }

    function renderStatus(status: PermissionStatus) {
        return <Button sx={SX.button} disabled={true} size={"small"}>{PermissionStatus[status]}</Button>
    }

    function renderApprove(label: string) {
        const r = {username, permission: name}
        return <Button sx={SX.button} color={"success"} size={"small"} loading={approve.isPending} onClick={() => approve.mutate(r)}>{label}</Button>
    }
    function renderReject(label: string) {
        const r = {username, permission: name}
        return <Button sx={SX.button} color={"warning"} size={"small"} loading={reject.isPending} onClick={() => reject.mutate(r)}>{label}</Button>
    }
    function renderRequest() {
        return (
            <Button sx={SX.button} size={"small"} loading={request.isPending} onClick={() => request.mutate(name)}>Request</Button>
        )
    }
}