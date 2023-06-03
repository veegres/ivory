import {StartupBlock} from "./StartupBlock";
import {LoadingButton} from "@mui/lab";
import {KeyEnterInput} from "../view/input/KeyEnterInput";
import {StartupConfigAuth} from "./StartupConfigAuth";
import {StartupConfigQuery} from "./StartupConfigQuery";
import {useState} from "react";
import {Authentication, AuthType} from "../../type/common";

export function StartupInitialConfig() {
    const [company, setCompany] = useState("")
    const [query, setQuery] = useState(false)
    const [auth, setAuth] = useState<Authentication>({type: AuthType.NONE, body: undefined})

    console.log(company, query, auth)

    return (
        <StartupBlock header={"Configuration"} renderFooter={renderFooter()}>
            <KeyEnterInput label={"Company"} onChange={(e) => setCompany(e.target.value)}/>
            <StartupConfigQuery onChange={setQuery}/>
            <StartupConfigAuth onChange={setAuth} auth={auth}/>
        </StartupBlock>
    )

    function renderFooter() {
        return (
            <LoadingButton>Finish</LoadingButton>
        )
    }
}
