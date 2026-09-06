import {Box, Button, Tooltip} from "@mui/material"
import {useState} from "react"

import {useRouterLogin} from "../../../features/auth/api/AuthHook"
import {useRouterUserRegistrationPassword, useRouterUserRegistrationVerify} from "../../../features/user/api/UserHook"
import {UserAuthType} from "../../../features/user/api/UserType"
import {AlertCentered} from "../../../shared/component/box/AlertCentered"
import {ErrorSmart} from "../../../shared/component/box/ErrorSmart"
import {PageStartupBox} from "../../../shared/component/box/PageStartupBox"
import {PageStartupGreeting} from "../../../shared/component/box/PageStartupGreeting"
import {KeyEnterInput} from "../../../shared/component/input/KeyEnterInput"
import {LogoProgress} from "../../../shared/component/progress/LogoProgress"
import {redirectToHome} from "../../../shared/helper/HelperUrl"
import {DateTimeFormatter} from "../../../shared/helper/HelperUtils"

type Props = {
    token: string,
}

export function UserRegistrationBody(props: Props) {
    const {token} = props
    const [password, setPassword] = useState("")
    const [repeat, setRepeat] = useState("")

    const registration = useRouterUserRegistrationVerify(token)
    const login = useRouterLogin(redirectToHome)
    const setPasswordRequest = useRouterUserRegistrationPassword(handlePasswordSet)

    if (registration.isPending) return <LogoProgress/>

    return (
        <PageStartupBox header={"Registration"} renderFooter={renderFooter()} position={"start"}>
            {renderContent()}
        </PageStartupBox>
    )

    function renderContent() {
        if (registration.isError) return <ErrorSmart error={registration.error}/>
        return (<>
            <PageStartupGreeting username={registration.data?.username}/>
            <AlertCentered severity={"info"} text={renderDescription()}/>
            <AlertCentered severity={"warning"} text={renderWarning()}/>
            <KeyEnterInput
                label={"Password"}
                value={password}
                hidden={true}
                onChange={(e) => setPassword(e.target.value)}
            />
            <KeyEnterInput
                label={getRepeatLabel()}
                value={repeat}
                hidden={true}
                error={isMismatched()}
                onChange={(e) => setRepeat(e.target.value)}
                onEnterPress={handleSet}
            />
        </>)
    }

    function renderDescription() {
        return (<>
            Somebody from your team registered you in <b>Ivory</b> as <b>{registration.data?.username}</b>.
            Set the password you will sign in with - you are taken straight into Ivory afterwards.
            The link works once and expires
            on <b>{registration.data ? DateTimeFormatter.utc(registration.data.expiresAt) : ""}</b>.
        </>)
    }

    function renderWarning() {
        return (
            "If you opened this link by accident, or you were not expecting it, close this page " +
            "and tell the person who sent it. The link stays usable until it is used or revoked."
        )
    }

    function renderFooter() {
        if (registration.isError) {
            return (
                <Tooltip title={"Open Ivory"} placement={"top"} arrow disableInteractive>
                    <Button variant={"contained"} onClick={redirectToHome} fullWidth>BACK</Button>
                </Tooltip>
            )
        }
        return (
            <Tooltip title={"Save the password and sign in to Ivory"} placement={"top"} arrow disableInteractive>
                <Box component={"span"}>
                    <Button
                        variant={"contained"}
                        loading={setPasswordRequest.isPending || login.isPending}
                        disabled={!password || password !== repeat}
                        onClick={handleSet}
                        fullWidth
                    >
                        SET PASSWORD
                    </Button>
                </Box>
            </Tooltip>
        )
    }

    function getRepeatLabel() {
        return isMismatched() ? "Repeat password - they do not match" : "Repeat password"
    }

    function handleSet() {
        if (!password || password !== repeat) return
        setPasswordRequest.mutate({token, password})
    }

    function handlePasswordSet() {
        login.mutate({type: UserAuthType.BASIC, subject: {username: registration.data?.username, password}})
    }

    function isMismatched() {
        return repeat !== "" && password !== repeat
    }
}
