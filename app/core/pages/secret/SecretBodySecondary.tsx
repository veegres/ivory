import {Typography} from "@mui/material"
import {useState} from "react"

import {useRouterSecretSet} from "../../../features/secret/api/SecretHook"
import {SecretButton} from "../../../features/secret/component/SecretButton"
import {PageStartupBox} from "../../../shared/component/box/PageStartupBox"
import {KeyEnterInput} from "../../../shared/component/input/KeyEnterInput"


export function SecretBodySecondary() {
    const [key, setKey] = useState("")
    const secret = useRouterSecretSet()

    return (
        <PageStartupBox header={"Welcome Back"} renderFooter={<SecretButton keyWord={key}/>}>
            <Typography variant={"caption"}>
                Oops! <b>Ivory</b> was just rebooted. Please enter your <b>Secret word</b> to continue working with sensitive data.
                There is nothing to press here if you have forgotten it: wiping everything is something only a
                signed-in administrator can ask for, from <i>Settings</i>. Without the word, the way to start
                fresh is to reinstall <b>Ivory</b>.
            </Typography>
            <KeyEnterInput
                label={"Secret word"}
                onChange={(e) => setKey(e.target.value)}
                onEnterPress={() => secret.mutate({key})}
                hidden
            />
        </PageStartupBox>
    )
}
