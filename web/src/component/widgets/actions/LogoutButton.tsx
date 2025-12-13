import {Button} from "@mui/material"

import {useRouterLogout} from "../../../api/auth/hook"

export function LogoutButton() {
    const logoutRouter = useRouterLogout()

    return (
        <Button onClick={() => logoutRouter.mutate()} loading={logoutRouter.isPending}>Sign out</Button>
    )
}