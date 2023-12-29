import {Accordion, AccordionDetails, AccordionSummary, Box} from "@mui/material";
import {SimpleStickyTable} from "../../view/table/SimpleStickyTable";
import {SxPropsMap} from "../../../type/common";
import {QueryFields} from "../../../type/query";
import {useState} from "react";
import {KeyboardArrowDown} from "@mui/icons-material";

const SX: SxPropsMap = {
    box: {borderBottom: 1, borderColor: "divider", "&::before": {display: "none"}, "&:last-child": {borderBottom: 0}},
    summary: {
        minHeight: "auto", padding: "2px 10px", gap: 2,
        "& .MuiAccordionSummary-content": {margin: "5px 0", gap: 1, alignItems: "center"}
    },
    details: {padding: "0 5px 5px"},
    info: {display: "flex", justifyContent: "space-between", fontSize: "12px", color: "text.secondary", gap: 1},
}

type Props = {
    query: QueryFields,
}

export function QueryBodyHistoryItem(props: Props) {
    const {query} = props
    const [open, setOpen] = useState(false)
    const [show, setShow] = useState(false)

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
            onMouseOut={() => setShow(false)}
            onMouseOver={() => setShow(true)}
        >
            <AccordionSummary sx={SX.summary} expandIcon={<KeyboardArrowDown sx={{fontSize: "20px"}}/>}>
                <Box>{endDate.toLocaleString()}</Box>
                {(open || show) && (
                    <Box sx={SX.info}>
                        <Box>[ {query.url} ]</Box>
                        <Box>[ {durationSec}s ]</Box>
                    </Box>
                )}
            </AccordionSummary>
            <AccordionDetails sx={SX.details}>
                <SimpleStickyTable rows={query.rows} columns={columns}/>
            </AccordionDetails>
        </Accordion>
    )
}
