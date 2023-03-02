import {ConsoleBlock} from "../../view/ConsoleBlock";
import {TransitionGroup} from "react-transition-group";
import {Collapse} from "@mui/material";
import {OverviewBloatJobItem} from "./OverviewBloatJobItem";
import {CompactTable, StylePropsMap, SxPropsMap} from "../../../app/types";

const SX: SxPropsMap = {
    jobsEmpty: {textAlign: "center"},
}


const style: StylePropsMap = {
    transition: {display: "flex", flexDirection: "column", gap: "10px"}
}

type Props = {
    list: CompactTable[],
}

export function OverviewBloatJob(props: Props) {
    const {list} = props
    if (list.length === 0) return <ConsoleBlock sx={SX.jobsEmpty}>There is no jobs yet</ConsoleBlock>

    return (
        <TransitionGroup style={style.transition}>
            {list.map((value) => (
                <Collapse key={value.uuid}>
                    <OverviewBloatJobItem key={value.uuid} compactTable={value}/>
                </Collapse>
            ))}
        </TransitionGroup>
    )
}
