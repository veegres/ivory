import { LinkOff, LockReset } from "@mui/icons-material";
import { Box } from "@mui/material";
import { useState } from "react";

import {DeleteIconButton, IconButton} from "../../../shared/component/button/IconButtons"
import {SxPropsMap} from "../../../shared/helper/HelperType"
import {Feature} from "../../Feature"
import {useRouterInfo} from "../../management/api/ManagementHook"
import {useHasAccess} from "../../management/component/ManageAccess"
import {
    useRouterUserDelete,
    useRouterUserPasswordReset,
    useRouterUserPasswordResetRevoke,
    useRouterUserUpdate,
} from "../api/UserHook"
import {User, UserAuthType, UserRegistration, UserRegistrationStatus} from "../api/UserType"
import {UserAuthTypes} from "./UserAuthTypes"
import {UserRegistrationLink} from "./UserRegistrationLink"

const SX: SxPropsMap = {
    box: {
        display: "flex", flexDirection: "column", gap: 0.5, padding: "3px",
        borderTop: 1, borderBottom: 1, borderColor: "divider",
    },
    footer: {display: "flex", justifyContent: "space-between", alignItems: "center", gap: 0.5, flexWrap: "wrap"},
    actions: {display: "flex", alignItems: "center", gap: 0.5, justifyContent: "space-between"},
}

type Props = {
    user: User,
}

export function UserUpdate(props: Props) {
    const {user} = props

    const [registration, setRegistration] = useState<UserRegistration>()
    const info = useRouterInfo()
    const deleteAccess = useHasAccess(Feature.ManageUserDelete)
    const updateAccess = useHasAccess(Feature.ManageUserUpdate)
    const resetAccess = useHasAccess(Feature.ManageUserPasswordReset)
    const deleteUser = useRouterUserDelete()
    const updateUser = useRouterUserUpdate()
    const issueReset = useRouterUserPasswordReset(setRegistration)
    const revokeReset = useRouterUserPasswordResetRevoke()

    return (
        <Box sx={SX.box}>
            <Box sx={SX.footer}>
                {renderAuthTypes()}
                <Box sx={SX.actions}>
                    {renderIssue()}
                    {renderRevoke()}
                    {renderDelete()}
                </Box>
            </Box>
            {registration && (
                <UserRegistrationLink registration={registration} reset/>
            )}
        </Box>
    )


    function renderAuthTypes() {
        const reason = getUpdateReason()
        return (
            <UserAuthTypes
                value={user.authTypes}
                size={"small"}
                disabled={!!reason || updateUser.isPending}
                reason={reason}
                onChange={(authTypes) => updateUser.mutate({username: user.username, body: {authTypes}})}
            />
        )
    }

    function renderIssue() {
        const reason = getIssueReason()
        return (
            <IconButton
                icon={<LockReset/>}
                tooltip={reason ?? "Reset the password - issues a link to set a new one"}
                disabled={!!reason}
                loading={issueReset.isPending}
                onClick={() => issueReset.mutate(user.username)}
            />
        )
    }

    function renderRevoke() {
        const reason = getRevokeReason()
        return (
            <IconButton
                icon={<LinkOff/>}
                tooltip={reason ?? "Make the outstanding link useless straight away"}
                disabled={!!reason}
                loading={revokeReset.isPending}
                onClick={() => revokeReset.mutate(user.username)}
            />
        )
    }

    function renderDelete() {
        const reason = getDeleteReason()
        return (
            <DeleteIconButton
                tooltip={reason ?? "Delete this user and their permissions"}
                disabled={!!reason}
                loading={deleteUser.isPending}
                onClick={() => deleteUser.mutate(user.username)}
            />
        )
    }


    function getUpdateReason() {
        if (updateAccess !== "allowed") return "You are not permitted to change how a user signs in"
        if (user.superuser && !isSuperuser()) return "Only a superuser can change a superuser"
        return undefined
    }

    function getIssueReason() {
        if (resetAccess !== "allowed") return "You are not permitted to reset passwords"
        if (!user.authTypes.includes(UserAuthType.BASIC)) return "This user does not sign in with a password"
        if (user.superuser && !isSuperuser()) return "Only a superuser can reset a superuser's password"
        return undefined
    }

    function getRevokeReason() {
        if (resetAccess !== "allowed") return "You are not permitted to revoke links"
        if (!isLinkOutstanding()) return "There is no link outstanding"
        if (user.superuser && !isSuperuser()) return "Only a superuser can revoke a superuser's link"
        return undefined
    }

    function getDeleteReason() {
        if (deleteAccess !== "allowed") return "You are not permitted to delete users"
        if (isYourself()) return "You cannot delete yourself"
        if (user.superuser && !isSuperuser()) return "Only a superuser can delete a superuser"
        return undefined
    }

    function isLinkOutstanding() {
        const status = user.registration?.status
        return status === UserRegistrationStatus.PENDING || status === UserRegistrationStatus.EXPIRED
    }

    function isYourself() {
        return info.data?.auth.user?.username === user.username
    }

    function isSuperuser() {
        return info.data?.auth.user?.superuser === true
    }
}