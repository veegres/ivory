import {Table, TableCell, TableRow} from "@mui/material";
import {TableBody} from "../../../view/table/TableBody";
import {InstanceWeb} from "../../../../api/instance/type";
import {SxPropsMap} from "../../../../app/type";

const SX: SxPropsMap = {
    title: {color: "text.disabled", fontWeight: "bold", verticalAlign: "top"},
    table: {td: {border: 0}, fontSize: "18px"},
    wrap: {wordBreak: "break-all"},
}

type Props = {
    instance: InstanceWeb,
}

export function InstanceInfoTable(props: Props) {
    const {instance} = props

    return (
        <Table size={"small"} sx={SX.table}>
            <TableBody isLoading={false} cellCount={2}>
                <TableRow>
                    <TableCell sx={SX.title}>State</TableCell>
                    <TableCell sx={SX.wrap}>{instance.state}</TableCell>
                </TableRow>
                <TableRow>
                    <TableCell sx={SX.title}>Sidecar</TableCell>
                    <TableCell sx={SX.wrap}>{instance.sidecar.host}:{instance.sidecar.port.toString()}</TableCell>
                </TableRow>
                <TableRow>
                    <TableCell sx={SX.title}>Database</TableCell>
                    <TableCell sx={SX.wrap}>{instance.database.host}:{instance.database.port.toString()}</TableCell>
                </TableRow>
            </TableBody>
        </Table>
    )
}
