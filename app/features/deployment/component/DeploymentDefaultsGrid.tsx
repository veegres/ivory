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
    inputText: {fontSize: "13px", fontFamily: "monospace", padding: "4px 8px"},
}

export type DefaultField = {
    variable: DeployVar,
    value: string,
    numeric?: boolean,
    // NOTE: disabled independently of the grid's own editable - a variable
    // that is never part of a template (the cluster name, a password, the
    // host) stays disabled even while the rest of the grid is being edited
    disabled?: boolean,
    // NOTE: shown as the disabled field's placeholder and appended to its
    // tooltip, so a reader sees why it is empty without having to ask
    hint?: string,
    onChange: (value: string) => void,
}

type Props = {
    fields: DefaultField[],
    editable?: boolean,
}

export function DeploymentDefaultsGrid(props: Props) {
    const {fields, editable = true} = props
    return <Box sx={SX.grid}>{fields.map(renderField)}</Box>

    function renderField(field: DefaultField) {
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
