import {Box} from "@mui/material";
import {ChangeEvent} from "react";

import {BasicConfig} from "../../../../api/config/type";
import {SxPropsMap} from "../../../../app/type";
import {KeyEnterInput} from "../../../view/input/KeyEnterInput";

const SX: SxPropsMap = {
    box: {display: "flex",  gap: 1},
}

type Props = {
    config: BasicConfig,
    onChange: (config: BasicConfig) => void,
}

export function ConfigAuthBasic(props: Props) {
    const {config, onChange} = props
    return (
        <Box sx={SX.box}>
            <KeyEnterInput label={"Username"} value={config.username} onChange={handleConfigChange("username")}/>
            <KeyEnterInput label={"Password"} value={config.password} onChange={handleConfigChange("password")} hidden/>
        </Box>
    )

    function handleConfigChange(key: string) {
        return (e: ChangeEvent<HTMLInputElement>) => {
            onChange({...config, [key]: e.target.value})
        }
    }
}