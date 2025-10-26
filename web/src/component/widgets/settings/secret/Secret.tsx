import {Alert, Box, Button} from "@mui/material"
import {useState} from "react"

import {useRouterSecretChange} from "../../../../api/management/hook"
import {SxPropsMap} from "../../../../app/type"
import {KeyEnterInput} from "../../../view/input/KeyEnterInput"
import {MenuWrapper} from "../menu/MenuWrapper"

const SX: SxPropsMap = {
    alert: {margin: "0px 10px"},
    text: {display: "flex", flexDirection: "column", gap: 1, textAlign: "center"},
    form: {display: "flex", flexDirection: "column", gap: 2, margin: "10px 0px"},
    button: {marginTop: "5px"},
    bold: {fontWeight: "bold"},
}

export function Secret() {
    const [prevKey, setPrevKey] = useState("")
    const [newKey, setNewKey] = useState("")
    const changeReq = useRouterSecretChange()

    return (
        <MenuWrapper>
            <Alert sx={SX.alert} severity={"info"} variant={"outlined"} icon={false}>
                <Box sx={SX.text}>
                    <Box>You can change the secret word here!</Box>
                    <Box sx={SX.bold}>
                        It can help if your secret was compromised or if someone got access
                        to the Ivory and you want them to be logged out
                    </Box>
                    <Box>
                        By changing the secret <b>Ivory</b> will reencrypt all of your passwords and
                        all the Ivory login tokens will become invalid (it means that you and all who
                        have been logged in are going to be logged out)
                    </Box>
                </Box>
            </Alert>
            <Box sx={SX.form}>
                <KeyEnterInput
                    hidden
                    label={"Previous Secret"}
                    onChange={(e) => setPrevKey(e.target.value)}
                />
                <KeyEnterInput
                    hidden
                    label={"New Secret"}
                    onChange={(e) => setNewKey(e.target.value)}
                />
                <Button
                    sx={SX.button}
                    variant={"contained"}
                    loading={changeReq.isPending}
                    onClick={() => changeReq.mutate({previousKey: prevKey, newKey})}
                >
                    Change
                </Button>
            </Box>
        </MenuWrapper>
    )
}
