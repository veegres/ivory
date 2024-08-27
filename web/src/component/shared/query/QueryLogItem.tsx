import {Accordion, AccordionDetails, AccordionSummary, Box} from "@mui/material";
import {VirtualizedTable} from "../../view/table/VirtualizedTable";
import {SxPropsMap} from "../../../type/general";
import {QueryFields} from "../../../type/query";
import {useMemo, useState} from "react";
import {KeyboardArrowDown} from "@mui/icons-material";

const SX: SxPropsMap = {
    box: {borderBottom: 1, borderColor: "divider", "&::before": {display: "none"}, "&:last-child": {borderBottom: 0}},
    summary: {
        minHeight: "auto", padding: "2px 10px", gap: 2,
        "& .MuiAccordionSummary-content": {margin: "5px 0", gap: 1, alignItems: "center"}
    },
    details: {padding: "0 5px 5px"},
    info: {display: "flex", justifyContent: "space-between", fontSize: "12px", color: "text.secondary", gap: 1},
    number: {color: "text.secondary", width: "20px", textAlign: "center"},
}

type Props = {
    index: number,
    query: QueryFields,
}

export function QueryLogItem(props: Props) {
    const {query, index} = props
    const [open, setOpen] = useState(false)
    const [show, setShow] = useState(false)

    const columns = useMemo(handleMemoColumns, [query.fields])

    const endDate = new Date(query.endTime)
    const durationSec = (query.endTime - query.startTime) / 1000
    return (
        <Accordion
            sx={SX.box}
            expanded={open}
            slotProps={{transition: {unmountOnExit: true}}}
            onChange={() => setOpen(!open)}
            disableGutters={true}
            onMouseOut={() => setShow(false)}
            onMouseOver={() => setShow(true)}
        >
            <AccordionSummary sx={SX.summary} expandIcon={<KeyboardArrowDown sx={{fontSize: "20px"}}/>}>
                <Box sx={SX.number}>{index+1}</Box>
                <Box>{endDate.toLocaleString()}</Box>
                {(open || show) && (
                    <Box sx={SX.info}>
                        <Box>[ {query.url} ]</Box>
                        <Box>[ {durationSec}s ]</Box>
                    </Box>
                )}
            </AccordionSummary>
            <AccordionDetails sx={SX.details}>
                <VirtualizedTable rows={query.rows} columns={columns}/>
            </AccordionDetails>
        </Accordion>
    )

    function handleMemoColumns() {
        return query.fields.map(field => ({name: field.name, description: field.dataType}))
    }
}
