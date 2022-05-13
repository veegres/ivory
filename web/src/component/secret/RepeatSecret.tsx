import React, {useState} from "react";
import {Secret} from "./Secret";
import {SecretTextField} from "./SecretTextField";

export function RepeatSecret() {
    const [key, setKey] = useState("")

    return (
        <Secret keyWord={key} refWord={""} clean={true}>
            <SecretTextField label={"Secret word"} onChange={(e) => setKey(e.target.value)} />
        </Secret>
    )
}
