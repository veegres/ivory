import {Alert, Box, Button} from "@mui/material"
import {useState} from "react"

import {useRouterSecretChange} from "../../../../api/management/hook"
import {SxPropsMap} from "../../../../app/type"
import {KeyEnterInput} from "../../../view/input/KeyEnterInput"
import {MenuWrapperScroll} from "../menu/MenuWrapperScroll"

const SX: SxPropsMap = {
    text: {display: "flex", flexDirection: "column", gap: 1, textAlign: "center"},
    form: {display: "flex", flexDirection: "column", gap: 2, margin: "20px 0px"},
    bold: {fontWeight: "bold", fontSize: "14px"},
    desc: {textAlign: "justify", fontSize: "12px", padding: "0px 25px"},
}

export function Secret() {
    const [prevKey, setPrevKey] = useState("")
    const [newKey, setNewKey] = useState("")
    const changeReq = useRouterSecretChange()

    return (
        <MenuWrapperScroll>
            <Alert severity={"info"} variant={"outlined"} icon={false}>
                <Box sx={SX.text}>
                    <Box>You can change your secret here.</Box>
                    <Box sx={SX.bold}>
                        This is useful if your secret has been compromised or if someone gained access
                        to Ivory and you want to force all users to log out.
                    </Box>
                    <Box sx={SX.desc}>
                        When you change the secret, Ivory will re-encrypt all your stored passwords, and
                        all existing Ivory login tokens will be invalidated. This means that you and
                        everyone else currently logged in will be logged out.
                        If you skipped setting a secret during the initial setup, you can set it here—just
                        leave the <i>Previous secret</i> field empty. You can also revert to the default behavior
                        by leaving the <i>New secret</i> field empty.
                    </Box>
                </Box>
            </Alert>
            <Box sx={SX.form}>
                <KeyEnterInput
                    hidden={true}
                    required={false}
                    label={"Previous Secret"}
                    onChange={(e) => setPrevKey(e.target.value)}
                />
                <KeyEnterInput
                    hidden={true}
                    required={false}
                    label={"New Secret"}
                    onChange={(e) => setNewKey(e.target.value)}
                />
                <Button
                    variant={"contained"}
                    loading={changeReq.isPending}
                    onClick={() => changeReq.mutate({previousKey: prevKey, newKey})}
                >
                    Change
                </Button>
            </Box>
        </MenuWrapperScroll>
    )
}
