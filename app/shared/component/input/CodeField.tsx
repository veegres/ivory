import {Box} from "@mui/material"

import {SxPropsMap} from "../../helper/HelperType"
import {Caption} from "../box/Caption"
import {Note} from "../box/Note"
import {CodeEditor} from "./CodeEditor"

const SX: SxPropsMap = {
    box: {display: "flex", flexDirection: "column", gap: 0.5},
    // NOTE: a small inset so the label sits over the field it names rather
    // than against the edge of the box around it
    // NOTE: a small row gap - the hint only takes a row of its own when it
    // will not fit beside the label, and there it belongs to the label rather
    // than reading as a line of its own
    head: {
        display: "flex", alignItems: "baseline", flexWrap: "wrap",
        columnGap: 1, rowGap: 0.25, padding: "0px 5px",
    },
}

type Props = {
    label: string,
    value: string,
    editable: boolean,
    hint?: string,
    placeholder?: string,
    minHeight?: string,
    onUpdate?: (value: string) => void,
}

// CodeField is a labelled block of code - the deploy dialog's only way to show
// one, so a command and a post script are the same field at different sizes
// rather than two things that happen to look similar.
export function CodeField(props: Props) {
    const {label, value, editable, hint, placeholder, minHeight = "80px", onUpdate} = props

    return (
        <Box sx={SX.box}>
            <Box sx={SX.head}>
                <Caption>{label}</Caption>
                {hint && <Note>{hint}</Note>}
            </Box>
            <CodeEditor
                value={value}
                editable={editable}
                minHeight={minHeight}
                placeholder={placeholder}
                onUpdate={onUpdate}
            />
        </Box>
    )
}
