import {Box, Button, Tooltip} from "@mui/material"
import {useState} from "react"

import {AlertCentered} from "../../../shared/component/box/AlertCentered"
import {KeyEnterInput} from "../../../shared/component/input/KeyEnterInput"
import {SxPropsMap} from "../../../shared/helper/HelperType"
import {useRouterUserPassword} from "../api/UserHook"

const SX: SxPropsMap = {
    box: {display: "flex", flexDirection: "column", gap: 2},
    inputs: {display: "flex", flexDirection: "column", gap: 1.5},
    button: {padding: "5px"},
}

export function UserAccount() {
    const [previousPassword, setPreviousPassword] = useState("")
    const [newPassword, setNewPassword] = useState("")
    const [repeat, setRepeat] = useState("")
    const updatePassword = useRouterUserPassword(handleUpdated)

    return (
        <Box sx={SX.box}>
            <AlertCentered text={renderDescription()}/>
            <Box sx={SX.inputs}>
                <KeyEnterInput
                    label={"Current password"}
                    value={previousPassword}
                    hidden
                    onChange={(e) => setPreviousPassword(e.target.value)}
                />
                <KeyEnterInput
                    label={"New password"}
                    value={newPassword}
                    hidden
                    onChange={(e) => setNewPassword(e.target.value)}
                />
                <KeyEnterInput
                    label={getRepeatLabel()}
                    value={repeat}
                    hidden
                    error={isMismatched()}
                    onChange={(e) => setRepeat(e.target.value)}
                    onEnterPress={handleUpdate}
                />
            </Box>
            {renderButton()}
        </Box>
    )

    function renderButton() {
        return (
            <Tooltip title={"Replace your password with the new one"} placement={"top"} arrow disableInteractive>
                <Box component={"span"}>
                    <Button
                        sx={SX.button}
                        fullWidth={true}
                        loading={updatePassword.isPending}
                        disabled={!isComplete()}
                        onClick={handleUpdate}
                    >
                        Change password
                    </Button>
                </Box>
            </Tooltip>
        )
    }

    function renderDescription() {
        return (
            "This changes the password of your Ivory user. It is the only thing you can change about " +
            "yourself - a username is never updated, and nobody but you can set your password."
        )
    }

    function getRepeatLabel() {
        return isMismatched() ? "Repeat new password - they do not match" : "Repeat new password"
    }

    function handleUpdate() {
        if (!isComplete()) return
        updatePassword.mutate({previousPassword, newPassword})
    }

    function handleUpdated() {
        setPreviousPassword("")
        setNewPassword("")
        setRepeat("")
    }

    function isComplete() {
        return !!previousPassword && !!newPassword && newPassword === repeat
    }

    function isMismatched() {
        return repeat !== "" && newPassword !== repeat
    }
}
