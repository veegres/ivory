import {Box} from "@mui/material";
import {Database} from "../../../api/query/type";
import {SxPropsMap} from "../../../app/type";

const SX: SxPropsMap = {
    box: {display: "flex", columnGap: 2, justifyContent: "space-evenly", flexWrap: "wrap"},
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
                <Box>{db.name ? db.name : "postgres"}</Box>
                {!db.name && (<Box sx={SX.default}>(default)</Box>)}
            </Box>
        </Box>
    )
}
