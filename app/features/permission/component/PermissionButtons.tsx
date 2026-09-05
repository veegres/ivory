import {Box, Button} from "@mui/material"

import {SxPropsMap} from "../../../shared/helper/HelperType"
import {Feature} from "../../Feature"
import {ManageAccess} from "../../management/component/ManageAccess"
import {
    useRouterPermissionApprove,
    useRouterPermissionReject,
    useRouterPermissionRequest,
} from "../api/PermissionHook"
import {Status} from "../api/PermissionType"

const SX: SxPropsMap = {
    box: {display: "flex", alignItems: "center"},
}

type Props = {
    username: string,
    permissions: [Feature, Status][],
    approve?: boolean,
    reject?: boolean,
    request?: boolean,
    count?: boolean,
}

export function PermissionButtons(props: Props) {
    const {username, permissions, approve = false, reject = false, request = false, count = false} = props
    const requestRouter = useRouterPermissionRequest()
    const approveRouter = useRouterPermissionApprove()
    const rejectRouter = useRouterPermissionReject()

    return (
        <Box sx={SX.box}>
            {request && renderRequest()}
            <ManageAccess feature={Feature.ManagePermissionUpdate}>{approve && renderApprove()}</ManageAccess>
            <ManageAccess feature={Feature.ManagePermissionUpdate}>{reject && renderReject()}</ManageAccess>
        </Box>
    )

    function renderApprove() {
        const predicate = (s: Status) => s !== Status.GRANTED
        const p = getFilteredPermissions(predicate)
        if (count && p.length <= 1) return
        return <Button
            color={"success"}
            loading={approveRouter.isPending}
            disabled={p.length === 0}
            onClick={() => approveRouter.mutate({username, permissions: p})}
            endIcon={count && p.length}
        >
            Grant
        </Button>
    }
    function renderReject() {
        const predicate = (s: Status) => s !== Status.NOT_PERMITTED
        const p = getFilteredPermissions(predicate)
        if (count && p.length <= 1) return
        return <Button
            color={"warning"}
            loading={rejectRouter.isPending}
            disabled={p.length === 0}
            onClick={() => rejectRouter.mutate({username, permissions: p})}
            endIcon={count && p.length}
        >
            Revoke
        </Button>
    }
    function renderRequest() {
        const predicate = (s: Status) => s === Status.NOT_PERMITTED
        const p = getFilteredPermissions(predicate)
        if (count && p.length <= 1) return
        return (
            <Button
                color={"secondary"}
                loading={requestRouter.isPending}
                disabled={p.length === 0}
                onClick={() => requestRouter.mutate(p)}
                endIcon={count && p.length}
            >
                Request
            </Button>
        )
    }

    function getFilteredPermissions(predicate: (status: Status) => boolean) {
        return permissions.filter(([_, status]) => predicate(status)).map(v => v[0])
    }
}