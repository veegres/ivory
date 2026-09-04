import {LinkOff, LockReset, Stars} from "@mui/icons-material"
import {Box, Tooltip} from "@mui/material"
import {useState} from "react"

import {InfoColorBox} from "../../../shared/component/box/InfoColorBox"
import {InfoColorBoxRow} from "../../../shared/component/box/InfoColorBoxRow"
import {DeleteIconButton, IconButton} from "../../../shared/component/button/IconButtons"
import {SxPropsMap} from "../../../shared/helper/HelperType"
import {DateTimeFormatter} from "../../../shared/helper/HelperUtils"
import {Feature} from "../../Feature"
import {useRouterInfo} from "../../management/api/ManagementHook"
import {useHasAccess} from "../../management/component/ManageAccess"
import {
    useRouterUserDelete,
    useRouterUserPasswordReset,
    useRouterUserPasswordResetRevoke,
    useRouterUserUpdate
} from "../api/UserHook"
import {User, UserAuthType, UserRegistration, UserRegistrationStatus} from "../api/UserType"
import {UserAuthTypes} from "./UserAuthTypes"
import {UserRegistrationLink} from "./UserRegistrationLink"

const SX: SxPropsMap = {
    row: {
        display: "flex", flexDirection: "column", gap: 0.5, padding: "5px 8px",
        border: 1, borderColor: "divider", borderRadius: 2,
    },
    title: {display: "flex", justifyContent: "space-between", alignItems: "center", gap: 1, overflow: "hidden"},
    name: {
        display: "flex", alignItems: "center", gap: 0.5, fontWeight: 600, textFamily: "monospace",
        whiteSpace: "nowrap", overflow: "hidden", textOverflow: "ellipsis", lineHeight: 1,
    },
    icon: {fontSize: "19px"},
    footer: {display: "flex", justifyContent: "space-between", alignItems: "end", gap: 0.5, flexWrap: "wrap"},
    actions: {display: "flex", alignItems: "center", gap: 0.5, justifyContent: "end", flexGrow: 1},
}

const STATUS: {[key in UserRegistrationStatus]: {label: string, color: string}} = {
    [UserRegistrationStatus.ACTIVE]: {label: "registered", color: "success"},
    [UserRegistrationStatus.PENDING]: {label: "link outstanding", color: "info"},
    [UserRegistrationStatus.EXPIRED]: {label: "link expired", color: "warning"},
    [UserRegistrationStatus.MISSING]: {label: "no password", color: "inherit"},
}

type Props = {
    user: User,
}

export function UserListItem(props: Props) {
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
        <Box sx={SX.row}>
            <Box sx={SX.title}>
                <Box sx={SX.name}>
                    <Tooltip title={user.superuser ? "Superuser" : "User"} placement={"top"} arrow disableInteractive>
                        <Stars sx={SX.icon} color={user.superuser ? "inherit" : "disabled"}/>
                    </Tooltip>
                    {user.username}
                </Box>
                <Box sx={SX.actions}>
                    <InfoColorBoxRow>{renderTags()}</InfoColorBoxRow>
                </Box>
            </Box>
            <Box sx={SX.footer}>
                {renderAuthTypes()}
                <Box sx={SX.actions}>
                    {renderIssue()}
                    {renderRevoke()}
                    {renderDelete()}
                </Box>
            </Box>
            {registration && <UserRegistrationLink registration={registration} reset/>}
        </Box>
    )


    function renderTags() {
        const status = user.registration?.status
        if (!status) return <InfoColorBox label={"not configured"} title={"This user doesn't use BASIC auth"}/>
        const {label, color} = STATUS[status]
        return <InfoColorBox label={label} color={color} title={getStatusTitle()}/>
    }

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

    function getStatusTitle() {
        const expiresAt = user.registration?.expiresAt
        if (!expiresAt) return undefined
        return `Expires on ${DateTimeFormatter.utc(expiresAt)}`
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
