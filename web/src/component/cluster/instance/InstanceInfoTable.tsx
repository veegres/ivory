import {Table, TableCell, TableRow} from "@mui/material";
import {TableBody} from "../../view/table/TableBody";
import {SxPropsMap} from "../../../type/common";
import {InstanceWeb} from "../../../type/Instance";

const SX: SxPropsMap = {
    title: {color: "text.secondary", fontWeight: "bold"},
    table: {td: {border: 0}, fontSize: "18px"},
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
                    <TableCell>{instance.state}</TableCell>
                </TableRow>
                <TableRow>
                    <TableCell sx={SX.title}>Sidecar</TableCell>
                    <TableCell>{instance.sidecar.host}:{instance.sidecar.port.toString()}</TableCell>
                </TableRow>
                <TableRow>
                    <TableCell sx={SX.title}>Database</TableCell>
                    <TableCell>{instance.database.host}:{instance.database.port.toString()}</TableCell>
                </TableRow>
            </TableBody>
        </Table>
    )
}
