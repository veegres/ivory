import {Alert, Box, Button} from "@mui/material"
import {useState} from "react"

import {useRouterSecretChange} from "../../../../api/management/hook"
import {SxPropsMap} from "../../../../app/type"
import {KeyEnterInput} from "../../../view/input/KeyEnterInput"
import {MenuWrapperScroll} from "../menu/MenuWrapperScroll"

const SX: SxPropsMap = {
    text: {display: "flex", flexDirection: "column", gap: 1, textAlign: "center"},
    form: {display: "flex", flexDirection: "column", gap: 2, margin: "20px 0px"},
    bold: {fontWeight: "bold"},
}

export function Secret() {
    const [prevKey, setPrevKey] = useState("")
    const [newKey, setNewKey] = useState("")
    const changeReq = useRouterSecretChange()

    return (
        <MenuWrapperScroll>
            <Alert severity={"info"} variant={"outlined"} icon={false}>
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
