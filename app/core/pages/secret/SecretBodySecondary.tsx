import {Typography} from "@mui/material"
import {useState} from "react"

import {ManageEraseButton} from "../../../features/management/component/ManageEraseButton"
import {useRouterSecretSet} from "../../../features/secret/api/hook"
import {SecretButton} from "../../../features/secret/component/SecretButton"
import {PageStartupBox} from "../../../shared/component/box/PageStartupBox"
import {KeyEnterInput} from "../../../shared/component/input/KeyEnterInput"


export function SecretBodySecondary() {
    const [key, setKey] = useState("")
    const secret = useRouterSecretSet()

    return (
        <PageStartupBox header={"Welcome Back"} renderFooter={renderButtons()}>
            <Typography variant={"caption"}>
                Oops! <b>Ivory</b> was just rebooted. Please enter your <b>Secret word</b> to continue working with sensitive data.
                If you’ve forgotten it or want to start fresh, simply press <b>ERASE</b> to remove all sensitive data and set a new <b>Secret word</b>.
            </Typography>
            <KeyEnterInput
                label={"Secret word"}
                onChange={(e) => setKey(e.target.value)}
                onEnterPress={() => secret.mutate({key})}
                hidden
            />
        </PageStartupBox>
    )

    function renderButtons() {
        return (
            <>
                <ManageEraseButton safe={false}/>
                <SecretButton keyWord={key}/>
            </>
        )
    }
}
