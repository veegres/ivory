import {Table, TableCell, TableRow} from "@mui/material";
import {TableBody} from "../../view/TableBody";
import {SxPropsMap} from "../../../type/common";
import {ActiveInstance, InstanceInfo} from "../../../type/instance";

const SX: SxPropsMap = {
    title: {color: "text.secondary", fontWeight: "bold"},
    table: {td: {border: 0}, fontSize: "18px"},
}

type Props = {
    loading: boolean,
    activeInstance: ActiveInstance,
    instance: InstanceInfo,
}

export function InstanceTable(props: Props) {
    const {loading, activeInstance, instance} = props

    return (
        <Table size={"small"} sx={SX.table}>
            <TableBody isLoading={loading} cellCount={2}>
                <TableRow>
                    <TableCell sx={SX.title}>State</TableCell>
                    <TableCell>{instance.state}</TableCell>
                </TableRow>
                <TableRow>
                    <TableCell sx={SX.title}>Sidecar</TableCell>
                    <TableCell>{activeInstance.sidecar.host}:{activeInstance.sidecar.port.toString()}</TableCell>
                </TableRow>
                <TableRow>
                    <TableCell sx={SX.title}>Database</TableCell>
                    <TableCell>{activeInstance.database.host}:{activeInstance.database.port.toString()}</TableCell>
                </TableRow>
            </TableBody>
        </Table>
    )
}
