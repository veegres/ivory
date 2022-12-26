import React, {useState} from "react";
import {StartupBlock} from "./StartupBlock";
import {StartupTextField} from "./StartupTextField";
import {Typography} from "@mui/material";



export function StartupReboot() {
    const [key, setKey] = useState("")

    return (
        <StartupBlock keyWord={key} refWord={""} clean={true} header={"Welcome Back"}>
            <Typography variant={"caption"}>
                Looks like Ivory was rebooted. Please, provide the <b>Secret word</b> to be able
                to continue working with sensitive data.
                You can clean the secret word and sensitive data by press <b>Clean</b> button and
                start working with Ivory from scratch by providing new secret word. Also it will help
                if you suddenly forget the secret word.
            </Typography>
            <StartupTextField label={"Secret word"} onChange={(e) => setKey(e.target.value)} hidden/>
        </StartupBlock>
    )
}
