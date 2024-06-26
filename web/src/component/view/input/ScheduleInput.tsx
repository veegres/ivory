import {DateTimeField} from "@mui/x-date-pickers";
import {DateTimeFormatter} from "../../../app/utils";
import dayjs, {Dayjs} from "dayjs";
import {Box, Tooltip} from "@mui/material";
import {SxPropsMap} from "../../../type/general";

const SX: SxPropsMap = {
    box: {display: "flex", alignItems: "center", gap: 2},
    time: {flexGrow: 1},
    zone: {cursor: "pointer"},
}

type Props = {
    value: Dayjs | null,
    onChange: (v: Dayjs | null) => void,
}

export function ScheduleInput(props: Props) {
    const {value, onChange} = props
    return (
        <Box sx={SX.box}>
            <DateTimeField
                sx={SX.time}
                label={"Schedule"}
                size={"small"}
                disablePast={true}
                format={DateTimeFormatter.format}
                value={value ?? null}
                onChange={onChange}
            />
            <Tooltip
                placement={"top"}
                arrow={true}
                title={`This is your local timezone. You need to provide 
                 time in your timezone, it is going to be converted into UTC.`}
            >
                <Box sx={SX.zone}>{dayjs().format("Z")}</Box>
            </Tooltip>
        </Box>
    )
}
