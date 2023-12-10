import {Database, SxPropsMap} from "../../../type/common";
import {Box} from "@mui/material";

const SX: SxPropsMap = {
    box: {display: "flex", gap: 2},
    item: {display: "flex", gap: 1},
    label: {fontWeight: "bold", color: "text.secondary"},
    default: {color: "text.disabled"},
}

type Props = {
    db: Database,
}

export function DatabaseBox(props: Props) {
    const {db} = props

    return (
        <Box sx={SX.box}>
            <Box sx={SX.item}>
                <Box sx={SX.label}>Host:</Box><Box>{db.host}</Box>
            </Box>
            <Box sx={SX.item}>
                <Box sx={SX.label}>Port:</Box><Box>{db.port}</Box>
            </Box>
            <Box sx={SX.item}>
                <Box sx={SX.label}>Database:</Box>
                <Box>{db.database ? db.database : "postgres"}</Box>
                {!db.database && (<Box sx={SX.default}>(default)</Box>)}
            </Box>
        </Box>
    )
}
