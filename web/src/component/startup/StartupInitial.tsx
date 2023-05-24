import {useState} from "react";
import {StartupBlock} from "./StartupBlock";
import {SecretInput} from "../shared/secret/SecretInput";
import {Typography} from "@mui/material";
import {SecretButton} from "../shared/secret/SecretButton";
import {useMutationOptions} from "../../hook/QueryCustom";
import {useMutation} from "@tanstack/react-query";
import {secretApi} from "../../app/api";

export function StartupInitial() {
    const [key, setKey] = useState("")
    const [ref, setRef] = useState("")
    const setReqOptions = useMutationOptions([["info"]])
    const setReq = useMutation(secretApi.set, setReqOptions)

    return (
        <StartupBlock header={"Welcome"} renderFooter={<SecretButton keyWord={key} refWord={ref}/>}>
            <Typography variant={"caption"}>
                This is <b>Ivory</b> — the tool that will help you to manage and troubleshoot your postgres clusters.
                Ivory needs some information to make your sensitive data safe.
                <ul style={{ padding: "0 10px "}}>
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
            <SecretInput
                label={"Reference word"}
                onChange={(e) => setRef(e.target.value)}
            />
            <SecretInput
                label={"Secret word"}
                onChange={(e) => setKey(e.target.value)}
                onEnterPress={() => setReq.mutate({ref, key})}
                hidden
            />
        </StartupBlock>
    )
}
