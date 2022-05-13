import React, {useState} from "react";
import {Secret} from "./Secret";
import {SecretTextField} from "./SecretTextField";

export function InitialSecret() {
    const [key, setKey] = useState("")
    const [ref, setRef] = useState("")

    return (
        <Secret keyWord={key} refWord={ref} clean={false}>
            <SecretTextField label={"Reference word"} onChange={(e) => setRef(e.target.value)}/>
            <SecretTextField label={"Secret word"} onChange={(e) => setKey(e.target.value)}/>
        </Secret>
    )
}
