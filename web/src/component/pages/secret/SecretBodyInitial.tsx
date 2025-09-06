import {useState} from "react";
import {PageStartupBox} from "../../view/box/PageStartupBox";
import {KeyEnterInput} from "../../view/input/KeyEnterInput";
import {Typography} from "@mui/material";
import {SecretButton} from "../../widgets/actions/SecretButton";
import {useRouterSecretSet} from "../../../api/secret/hook";

export function SecretBodyInitial() {
    const [key, setKey] = useState("")
    const [ref, setRef] = useState("")
    const secret = useRouterSecretSet()

    return (
        <PageStartupBox header={"Welcome"} renderFooter={<SecretButton keyWord={key} refWord={ref}/>}>
            <Typography variant={"caption"}>
                <b>This is Ivory</b> — the tool that will help you to manage and troubleshoot your postgres clusters.
                Ivory needs some information to make your sensitive data safe.
                <ul style={{paddingLeft: "20px"}}>
                    <li>
                        <b>Secret word</b> — this word will be used to encrypt and decrypt sensitive
                        data like passwords, tokens, etc. Please, don't spread this word among a lot of people
                        it will prevent leaking sensitive data. Ivory keeps this word only in memory,
                        it means every time when Ivory will be rebooted you need to pass this word, so
                        please do not forget it.
                    </li>
                    <li>
                        <b>Reference word</b> — this word will help Ivory to detect either your secret word is
                        correct or not after Ivory reboot.
                    </li>
                </ul>
            </Typography>
            <KeyEnterInput
                label={"Reference word"}
                onChange={(e) => setRef(e.target.value)}
            />
            <KeyEnterInput
                label={"Secret word"}
                onChange={(e) => setKey(e.target.value)}
                onEnterPress={() => secret.mutate({ref, key})}
                hidden
            />
        </PageStartupBox>
    )
}
