import React, {useState} from "react";
import {Secret} from "./Secret";
import {SecretTextField} from "./SecretTextField";
import {Typography} from "@mui/material";

export function InitialSecret() {
    const [key, setKey] = useState("")
    const [ref, setRef] = useState("")

    return (
        <Secret keyWord={key} refWord={ref} clean={false} header={"Welcome"}>
            <Typography variant={"caption"}>
                This is <b>Ivory</b> — the tool that will help you manage your postgres clusters. Please, provide
                some information to make your sensitive data safe while using Ivory.
                <ul style={{ padding: "0 10px "}}>
                    <li>
                        <b>Reference word</b> — this word will help Ivory to detect if your secret word is
                        correct if the system will be suddenly rebooted.
                    </li>
                    <li>
                        <b>Secret word</b> — this word will be used for encrypt and decrypt sensitive
                        data like passwords. Please, don't spread this word among a lot of people
                        it will prevent leaking sensitive data. Ivory keeps this word only in memory,
                        it means every time when Ivory will be rebooted you need to pass this word, so
                        please do not forget it.
                    </li>
                </ul>
            </Typography>
            <SecretTextField label={"Reference word"} onChange={(e) => setRef(e.target.value)}/>
            <SecretTextField label={"Secret word"} onChange={(e) => setKey(e.target.value)} hidden/>
        </Secret>
    )
}
