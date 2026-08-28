import {Box} from "@mui/material"

import {CodeToken} from "../../../shared/component/box/CodeToken"
import {Note} from "../../../shared/component/box/Note"
import {SubContentBox} from "../../../shared/component/box/SubContentBox"
import {DeployVarOptions} from "../../../shared/helper/HelperUtils"
import {DeployVar} from "../../node/api/NodeType"

// NOTE: not typed as SxPropsMap - Code takes a plain SystemStyleObject, and
// the annotation is what makes the two disagree
const SX = {
    // NOTE: auto-fit rather than a fixed column count - it takes as many pairs
    // per row as the box can hold, three where there is room and two where
    // there is not, and stretches them to use the width either way
    grid: {
        display: "grid", gridTemplateColumns: "repeat(auto-fit, minmax(170px, 1fr))",
        columnGap: 2, rowGap: 0,
    },
    pair: {display: "flex", alignItems: "center", gap: 1.5, minWidth: 0},
    // NOTE: the width is reserved by the cell, not by the chip - stretching the
    // chip itself made a list of tokens read as a row of buttons
    name: {minWidth: "104px"},
    // NOTE: nowrap - "from the vault" breaking in two is what made those rows
    // twice the height of the others
    example: {whiteSpace: "nowrap"},
}

// DeploymentVariables explains what a variable is for before listing any: a
// reader opening a template for the first time has no reason to guess that
// these are placeholders the deploy screen fills in - so that sentence sits in
// the heading, the way a CodeField's hint sits beside its label. Each variable
// sits next to what it turns into, which says more than a description would
// and costs a column instead of a line.
export function DeploymentVariables() {
    return (
        <SubContentBox
            label={"Variables"}
            hint={"Variables can be used in a command or its post script - Ivory replaces them with each node's values when you deploy"}
            island={true}
            collapsible={false}
        >
            <Box sx={SX.grid}>{Object.values(DeployVar).map(renderVariable)}</Box>
        </SubContentBox>
    )

    function renderVariable(name: DeployVar) {
        return (
            <Box key={name} sx={SX.pair}>
                <Box sx={SX.name}><CodeToken tooltip={getTitle(name)}>{name}</CodeToken></Box>
                <Box sx={SX.example}><Note>{DeployVarOptions[name].example}</Note></Box>
            </Box>
        )
    }

    // NOTE: a secret's tooltip has to say that what a preview shows is not
    // what runs - every other variable a reader sees filled in really is the
    // value that lands in the command
    function getTitle(name: DeployVar) {
        const {label, secret} = DeployVarOptions[name]
        if (!secret) return label
        return `${label} - a preview only shows the mask; the server substitutes the real one, so it never reaches the browser`
    }
}
