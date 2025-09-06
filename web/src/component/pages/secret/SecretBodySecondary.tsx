import {useState} from "react";
import {PageStartupBox} from "../../view/box/PageStartupBox";
import {KeyEnterInput} from "../../view/input/KeyEnterInput";
import {Typography} from "@mui/material";
import {EraseButton} from "../../widgets/actions/EraseButton";
import {SecretButton} from "../../widgets/actions/SecretButton";
import {useRouterSecretSet} from "../../../api/secret/hook";


export function SecretBodySecondary() {
    const [key, setKey] = useState("")
    const secret = useRouterSecretSet()

    return (
        <PageStartupBox header={"Welcome Back"} renderFooter={renderButtons()}>
            <Typography variant={"caption"}>
                Looks like Ivory was rebooted. Please, provide the <b>Secret word</b> to be able
                to continue working with sensitive data.
                You can start from scratch by cleaning the secret word and sensitive data, just
                press <b>ERASE</b> button. You will need to provide new secret word. Also it can
                help if you suddenly forget the secret word.
            </Typography>
            <KeyEnterInput
                label={"Secret word"}
                onChange={(e) => setKey(e.target.value)}
                onEnterPress={() => secret.mutate({ref: "", key})}
                hidden
            />
        </PageStartupBox>
    )

    function renderButtons() {
        return (
            <>
                <EraseButton safe={false}/>
                <SecretButton keyWord={key} refWord={""}/>
            </>
        )
    }
}
