import { Delete, LinkOff, LockReset } from "@mui/icons-material";
import { Box } from "@mui/material";
import { useState } from "react";

import {SimpleButton} from "../../../shared/component/button/SimpleButton"
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
        display: "flex", flexDirection: "column", gap: 0.5, paddingY: 1,
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
                supported={info.data?.auth.supported}
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
            <SimpleButton
                size={"small"}
                tooltip={reason ?? "Reset the password - issues a link to set a new one"}
                disabled={!!reason}
                loading={issueReset.isPending}
                onClick={() => issueReset.mutate(user.username)}
            >
                <LockReset fontSize={"small"}/>
            </SimpleButton>
        )
    }

    function renderRevoke() {
        const reason = getRevokeReason()
        return (
            <SimpleButton
                size={"small"}
                tooltip={reason ?? "Make the outstanding link useless straight away"}
                disabled={!!reason}
                loading={revokeReset.isPending}
                onClick={() => revokeReset.mutate(user.username)}
            >
                <LinkOff fontSize={"small"}/>
            </SimpleButton>
        )
    }

    function renderDelete() {
        const reason = getDeleteReason()
        return (
            <SimpleButton
                size={"small"}
                tooltip={reason ?? "Delete this user and their permissions"}
                disabled={!!reason}
                loading={deleteUser.isPending}
                onClick={() => deleteUser.mutate(user.username)}
            >
                <Delete fontSize={"small"}/>
            </SimpleButton>
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
        if (!isLinkIssued()) return "There is no link"
        if (user.superuser && !isSuperuser()) return "Only a superuser can revoke a superuser's link"
        return undefined
    }

    function getDeleteReason() {
        if (deleteAccess !== "allowed") return "You are not permitted to delete users"
        if (isYourself()) return "You cannot delete yourself"
        if (user.superuser && !isSuperuser()) return "Only a superuser can delete a superuser"
        return undefined
    }

    function isLinkIssued() {
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