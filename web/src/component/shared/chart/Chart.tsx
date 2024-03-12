import {Database, SxPropsMap} from "../../../type/common";
import {Box} from "@mui/material";
import {ChartRow} from "./ChartRow";
import {ChartCommon} from "./ChartCommon";
import {ChartDatabase} from "./ChartDatabase";

const SX: SxPropsMap = {
    box: {display: "flex", flexDirection: "column", justifyContent: "center", gap: 2},
    info: {display: "flex"},
}

type Props = {
    credentialId: string,
    db: Database,
}

export function Chart(props: Props) {
    const {credentialId, db} = props

    return (
        <Box sx={SX.box}>
            <ChartRow>
                <ChartCommon credentialId={credentialId} db={db}/>
            </ChartRow>
            <ChartRow label={db.name && `${db.name}`}>
                <ChartDatabase credentialId={credentialId} db={db}/>
            </ChartRow>
        </Box>
    )
}
