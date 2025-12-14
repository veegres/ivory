import {Box, Button} from "@mui/material"

import {
    useRouterPermissionApprove,
    useRouterPermissionReject,
    useRouterPermissionRequest,
} from "../../../../api/permission/hook"
import {PermissionStatus} from "../../../../api/permission/type"
import {SxPropsMap} from "../../../../app/type"

const SX: SxPropsMap = {
    box: {display: "flex", alignItems: "center"},
    button: {padding: "2px 5px"},
}

type Props = {
    username: string,
    permissions: [string, PermissionStatus][],
    approve?: boolean,
    reject?: boolean,
    request?: boolean,
    count?: boolean,
}

export function PermissionsButtons(props: Props) {
    const {username, permissions, approve = false, reject = false, request = false, count = false} = props
    const requestRouter = useRouterPermissionRequest()
    const approveRouter = useRouterPermissionApprove()
    const rejectRouter = useRouterPermissionReject()

    return (
        <Box sx={SX.box}>
            {request && renderRequest()}
            {approve && renderApprove()}
            {reject && renderReject()}
        </Box>
    )

    function renderApprove() {
        const predicate = (s: PermissionStatus) => s !== PermissionStatus.GRANTED
        const p = getFilteredPermissions(predicate)
        if (count && p.length <= 1) return
        return <Button
            sx={SX.button}
            color={"success"}
            size={"small"}
            loading={approveRouter.isPending}
            disabled={p.length === 0}
            onClick={() => approveRouter.mutate({username, permissions: p})}
            endIcon={count && p.length}
        >
            Grant
        </Button>
    }
    function renderReject() {
        const predicate = (s: PermissionStatus) => s !== PermissionStatus.NOT_PERMITTED
        const p = getFilteredPermissions(predicate)
        if (count && p.length <= 1) return
        return <Button
            sx={SX.button}
            color={"warning"}
            size={"small"}
            loading={rejectRouter.isPending}
            disabled={p.length === 0}
            onClick={() => rejectRouter.mutate({username, permissions: p})}
            endIcon={count && p.length}
        >
            Deny
        </Button>
    }
    function renderRequest() {
        const predicate = (s: PermissionStatus) => s === PermissionStatus.NOT_PERMITTED
        const p = getFilteredPermissions(predicate)
        if (count && p.length <= 1) return
        return (
            <Button
                sx={SX.button}
                color={"secondary"}
                size={"small"}
                loading={requestRouter.isPending}
                disabled={p.length === 0}
                onClick={() => requestRouter.mutate(p)}
                endIcon={count && p.length}
            >
                Request
            </Button>
        )
    }

    function getFilteredPermissions(predicate: (status: PermissionStatus) => boolean) {
        return permissions.filter(([_, status]) => predicate(status)).map(v => v[0])
    }
}