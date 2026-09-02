import {Box, TextField} from "@mui/material"

import {CodeToken} from "../../../shared/component/box/CodeToken"
import {DeployVarOptions} from "../../../shared/helper/HelperUtils"
import {DeployVar} from "../../node/api/NodeType"

// NOTE: not typed as SxPropsMap - Code takes a plain SystemStyleObject, and
// the annotation is what makes the two disagree
const SX = {
    grid: {display: "grid", gridTemplateColumns: "repeat(auto-fit, minmax(200px, 1fr))", columnGap: 1, rowGap: 0.5},
    pair: {display: "grid", alignItems: "center", gap: 0.5, gridTemplateColumns: "repeat(auto-fit, minmax(0, 1fr))", gridAutoFlow: "column"},
    input: {
        flexGrow: 1, minWidth: 0,
        "& .MuiOutlinedInput-notchedOutline": {borderColor: "divider"},
        "& .MuiOutlinedInput-root.Mui-disabled .MuiOutlinedInput-notchedOutline": {borderColor: "divider"},
    },
    inputText: {fontSize: "13px", fontFamily: "monospace", paddingLeft: "8px", paddingRight: "8px"},
}

export type DefaultField = {
    variable: DeployVar,
    value: string,
    numeric?: boolean,
    disabled?: boolean,
    hint?: string,
    onChange: (value: string) => void,
}

type Props = {
    fields: (DefaultField | undefined)[],
    editable?: boolean,
}

export function DeploymentDefaultsGrid(props: Props) {
    const {fields, editable = true} = props
    return <Box sx={SX.grid}>{fields.map(renderField)}</Box>

    function renderField(field: DefaultField | undefined, index: number) {
        if (!field) return <Box key={index}/>
        const {variable, value, numeric, disabled, onChange} = field
        return (
            <Box key={variable} sx={SX.pair}>
                <CodeToken tooltip={getTitle(field)}>{variable}</CodeToken>
                <TextField
                    sx={SX.input}
                    size={"small"}
                    type={numeric ? "number" : "text"}
                    disabled={!editable || disabled}
                    placeholder={DeployVarOptions[variable].example}
                    value={value}
                    slotProps={{htmlInput: {sx: SX.inputText}}}
                    onChange={(e) => onChange(e.target.value)}
                />
            </Box>
        )
    }

    function getTitle(field: DefaultField) {
        const {label} = DeployVarOptions[field.variable]
        return field.hint ? `${label} - ${field.hint}` : label
    }
}
