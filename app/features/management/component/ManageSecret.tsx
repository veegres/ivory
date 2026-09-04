import {Box, Button, Tooltip} from "@mui/material"
import {useState} from "react"

import {AlertCentered} from "../../../shared/component/box/AlertCentered"
import {KeyEnterInput} from "../../../shared/component/input/KeyEnterInput"
import {SxPropsMap} from "../../../shared/helper/HelperType"
import {useRouterSecretChange} from "../api/ManagementHook"

const SX: SxPropsMap = {
    box: {display: "flex", flexDirection: "column", gap: 2},
    inputs: {display: "flex", flexDirection: "column", gap: 1.5},
    button: {padding: "5px"},
}

export function ManageSecret() {
    const [prevKey, setPrevKey] = useState("")
    const [newKey, setNewKey] = useState("")
    const changeReq = useRouterSecretChange()

    return (
        <Box sx={SX.box}>
            <AlertCentered severity={"warning"} text={renderDescription()}/>
            <Box sx={SX.inputs}>
                <KeyEnterInput
                    label={"Previous secret"}
                    value={prevKey}
                    hidden
                    required={false}
                    onChange={(e) => setPrevKey(e.target.value)}
                />
                <KeyEnterInput
                    label={"New secret"}
                    value={newKey}
                    hidden
                    required={false}
                    onChange={(e) => setNewKey(e.target.value)}
                    onEnterPress={handleChange}
                />
            </Box>
            {renderButton()}
        </Box>
    )

    function renderButton() {
        return (
            <Tooltip title={"Re-encrypt everything with the new secret"} placement={"top"} arrow disableInteractive>
                <Box component={"span"}>
                    <Button
                        sx={SX.button}
                        fullWidth={true}
                        loading={changeReq.isPending}
                        onClick={handleChange}
                    >
                        Change secret
                    </Button>
                </Box>
            </Tooltip>
        )
    }

    function renderDescription() {
        return (<>
            Changing the secret re-encrypts every stored vault and invalidates every login token,
            so you and everybody else signed in are logged out. Leave <i>Previous secret</i> empty
            if you skipped setting one during the initial setup, and <i>New secret</i> empty to
            fall back to the default.
        </>)
    }

    function handleChange() {
        changeReq.mutate({previousKey: prevKey, newKey})
    }
}
