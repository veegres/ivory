import {KeyboardArrowDown} from "@mui/icons-material"
import {Accordion, AccordionDetails, AccordionSummary, Box} from "@mui/material"
import {useMemo, useState} from "react"

import {QueryFields} from "../../../api/postgres"
import {SxPropsMap} from "../../../app/type"
import {VirtualizedTable} from "../../view/table/VirtualizedTable"
import {QueryResponseInfo} from "./QueryResponseInfo"

const SX: SxPropsMap = {
    box: {borderBottom: 1, borderColor: "divider", "&::before": {display: "none"}, "&:last-child": {borderBottom: 0}},
    summary: {
        minHeight: "auto", padding: "2px 10px", gap: 2,
        "& .MuiAccordionSummary-content": {margin: "5px 0", gap: 1, alignItems: "center"}
    },
    details: {padding: "0 5px 5px"},
    number: {color: "text.secondary", width: "20px", textAlign: "center"},
}

type Props = {
    index: number,
    query: QueryFields,
}

export function QueryLogItem(props: Props) {
    const {query, index} = props
    const {fields, options, endTime, startTime} = query
    const [open, setOpen] = useState(false)
    const [show, setShow] = useState(false)

    const columns = useMemo(handleMemoColumns, [fields])

    const endDate = new Date(endTime)
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
                {!open && !show ? (
                    <Box>{endDate.toLocaleString()}</Box>
                ) : (
                    <QueryResponseInfo url={query.url} time={{start: startTime, end: endTime}} options={options}/>
                )}
            </AccordionSummary>
            <AccordionDetails sx={SX.details}>
                <VirtualizedTable rows={query.rows} columns={columns}/>
            </AccordionDetails>
        </Accordion>
    )

    function handleMemoColumns() {
        return fields.map(field => ({name: field.name, description: field.dataType}))
    }
}
