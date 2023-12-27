import {Accordion, AccordionDetails, AccordionSummary, Box} from "@mui/material";
import {SimpleStickyTable} from "../../view/table/SimpleStickyTable";
import {SxPropsMap} from "../../../type/common";
import {QueryFields} from "../../../type/query";
import {useState} from "react";
import {KeyboardArrowDown} from "@mui/icons-material";

const SX: SxPropsMap = {
    box: {borderBottom: 1, borderColor: "divider", "&::before": {display: "none"}, "&:last-child": {borderBottom: 0}},
    summary: {minHeight: "auto", padding: "2px 10px", "& .MuiAccordionSummary-content": {margin: "5px 0", gap: 1, alignItems: "center"}},
    details: {padding: "0 5px 5px"},
    duration: {fontSize: "10px", fontFamily: "monospace",  color: "text.disabled"},
}

type Props = {
    query: QueryFields,
}

export function QueryBodyHistoryItem(props: Props) {
    const {query} = props
    const [open, setOpen] = useState(false)

    const columns = query.fields.map(field => (
        {name: field.name, description: field.dataType}
    ))

    const endDate = new Date(query.endTime)
    const durationSec = (query.endTime - query.startTime) / 1000
    return (
        <Accordion
            sx={SX.box}
            expanded={open}
            TransitionProps={{unmountOnExit: true}}
            onChange={() => setOpen(!open)}
            disableGutters={true}
        >
            <AccordionSummary sx={SX.summary} expandIcon={<KeyboardArrowDown sx={{fontSize: "20px"}}/>}>
                <Box>{endDate.toLocaleString()}</Box>
                <Box sx={SX.duration}>({durationSec}s)</Box>
            </AccordionSummary>
            <AccordionDetails sx={SX.details}>
                <SimpleStickyTable rows={query.rows} columns={columns}/>
            </AccordionDetails>
        </Accordion>
    )
}
