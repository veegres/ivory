import {Box, SxProps, Theme} from "@mui/material"

import {Code} from "../../../../shared/component/box/Code"
import {DynamicInputs} from "../../../../shared/component/input/DynamicInputs"
import {ColorsMap, SxPropsMap} from "../../../../shared/helper/type"

const SX: SxPropsMap = {
    box: {display: "flex", flexDirection: "column", gap: 0.5},
    label: {display: "flex", fontSize: "8px", textTransform: "uppercase", lineHeight: 1.4, gap: 0.3},
}

type Props = {
    inputs: string[],
    onChange: (values: string[]) => void,
    colors?: ColorsMap,
    editable: boolean,
    minLength?: number,
    InputProps?: SxProps<Theme>,
}

export function ListNodeInput(props: Props) {
    const {inputs, InputProps, minLength, editable, colors, onChange} = props
    return (
        <DynamicInputs
            InputProps={InputProps}
            colors={colors}
            minLength={minLength}
            inputs={inputs}
            onChange={onChange}
            editable={editable}
            placeholder={"Node "}
            helper={renderPost()}
        />
    )

    function renderPost() {
        return (
            <Box sx={SX.label}>
                <Code>Host</Code>:<Code>Keeper Port</Code>:<Code>DB Port</Code>:<Code>SSH Port</Code>
            </Box>
        )
    }
}