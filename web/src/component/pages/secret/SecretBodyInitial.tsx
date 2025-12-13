import {Button, Typography} from "@mui/material"
import {useState} from "react"

import {useRouterSecretSet, useRouterSecretSkip} from "../../../api/secret/hook"
import {PageStartupBox} from "../../view/box/PageStartupBox"
import {KeyEnterInput} from "../../view/input/KeyEnterInput"
import {SecretButton} from "../../widgets/actions/SecretButton"

export function SecretBodyInitial() {
    const [key, setKey] = useState("")
    const secret = useRouterSecretSet()
    const skip = useRouterSecretSkip()

    return (
        <PageStartupBox header={"Welcome"} renderFooter={renderFooter()}>
            <Typography variant={"caption"}>
                Welcome to <b>Ivory</b> — your assistant for managing and troubleshooting PostgreSQL clusters.
                To keep your sensitive data safe, we need little information from you.
                <br/>
                <br/>
                <b>Secret word</b> — this word will be used to encrypt and decrypt sensitive information,
                like passwords and tokens. Keep it private and don’t share it with many people — this helps
                prevent data leaks. Ivory only keeps this word in memory, so each time you restart the tool,
                you’ll need to enter it again. Make sure you remember it! You may <i>skip</i> this step, in
                which case Ivory will apply a default value (this option is not recommended). You can modify
                it at any time.
            </Typography>
            <KeyEnterInput
                label={"Secret word"}
                onChange={(e) => setKey(e.target.value)}
                onEnterPress={() => secret.mutate({key})}
                hidden
            />
        </PageStartupBox>
    )

    function renderFooter() {
        return (
            <>
                <SecretButton keyWord={key}/>
                <Button color={"inherit"} variant={"contained"} size={"small"} onClick={() => skip.mutate()}>Skip</Button>
            </>
        )
    }
}
