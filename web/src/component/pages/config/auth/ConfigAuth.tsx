import {ChangeEvent} from "react";
import {AlertColor, Box, FormControlLabel, Radio, RadioGroup} from "@mui/material";
import {SxPropsMap} from "../../../../api/management/type";
import {ConfigBox} from "../../../view/box/ConfigBox";
import {KeyEnterInput} from "../../../view/input/KeyEnterInput";
import {AuthConfig, AuthType} from "../../../../api/config/type";

const SX: SxPropsMap = {
    radio: {display: "flex", gap: 3},
    basic: {display: "flex", gap: 1},
}

const Description = {
    [AuthType.NONE]: <>
        No authentications means that anyone will have direct access to <b>Ivory</b> and
        to all the information provided like certs, password, etc. Anyone who has a link to
        Ivory will be able to perform any action with database like switchover, reinit, etc.
    </>,
    [AuthType.BASIC]: <>
        Basic authentication helps to protect <b>Ivory</b> data from direct access
        with general username and password. It provides sign in form to get access
        to <b>Ivory</b>. It will generate unique token for each sign in. Each token
        has expiration time. You need to provide general username and password
        that will be used below.
    </>
}

const Recommendation = {
    [AuthType.NONE]: <>
        This type of authentication is recommended for <b>local usage</b>, when
        you use <b>Ivory</b> in your personal laptop / computer.
    </>,
    [AuthType.BASIC]: <>
        This type of authentication is recommended for <b>teams</b> or <b>small
        group of people</b>.
    </>
}

const Severity: { [key in AuthType]: AlertColor } = {
    [AuthType.NONE]: "warning",
    [AuthType.BASIC]: "info",
}

type Props = {
    onChange: (auth: AuthConfig) => void,
    auth: AuthConfig,
}

export function ConfigAuth(props: Props) {
    const {onChange, auth} = props

    return (
        <ConfigBox
            label={"Authentication"}
            renderAction={renderAction()}
            renderBody={renderBasic()}
            showBody={auth.type === AuthType.BASIC}
            description={Description[auth.type]}
            recommendation={Recommendation[auth.type]}
            severity={Severity[auth.type]}
        />
    )

    function renderAction() {
        return (
            <RadioGroup sx={SX.radio} row value={auth.type} onChange={handleAuthTypeChange}>
                <FormControlLabel
                    value={AuthType.NONE}
                    control={<Radio size={"small"}/>}
                    label={AuthType[AuthType.NONE].toLowerCase()}
                />
                <FormControlLabel
                    value={AuthType.BASIC}
                    control={<Radio size={"small"}/>}
                    label={AuthType[AuthType.BASIC].toLowerCase()}
                />
            </RadioGroup>
        )
    }

    function renderBasic() {
        return (
            <Box sx={SX.basic}>
                <KeyEnterInput label={"Username"} value={auth.body?.username ?? ""} onChange={handleUsernameChange}/>
                <KeyEnterInput label={"Password"} value={auth.body?.password ?? ""} onChange={handlePasswordChange} hidden />
            </Box>
        )
    }

    function handleUsernameChange(e: ChangeEvent<HTMLInputElement>) {
        onChange({...auth, body: {...auth.body, username: e.target.value}})
    }

    function handlePasswordChange(e: ChangeEvent<HTMLInputElement>) {
        onChange({...auth, body: {...auth.body, password: e.target.value}})
    }

    function handleAuthTypeChange(e: ChangeEvent<HTMLInputElement>) {
        const type = parseInt(e.target.value)
        onChange({type: type as AuthType, body: undefined})
    }
}
