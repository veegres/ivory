import {Typography} from "@mui/material";
import {useState} from "react";

import {useRouterSecretSet} from "../../../api/secret/hook";
import {PageStartupBox} from "../../view/box/PageStartupBox";
import {KeyEnterInput} from "../../view/input/KeyEnterInput";
import {EraseButton} from "../../widgets/actions/EraseButton";
import {SecretButton} from "../../widgets/actions/SecretButton";


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
                <EraseButton safe={false}/>
                <SecretButton keyWord={key}/>
            </>
        )
    }
}
