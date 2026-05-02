import {Table, TableCell, TableRow} from "@mui/material"

import {Node} from "../../../../api/cluster/type"
import {SxPropsMap} from "../../../../app/type"
import {TableBody} from "../../../view/table/TableBody"

const SX: SxPropsMap = {
    title: {color: "text.disabled", fontWeight: "bold", verticalAlign: "top"},
    table: {td: {border: 0}, fontSize: "18px"},
    wrap: {wordBreak: "break-all"},
}

type Props = {
    node: Node,
}

export function NodeInfoTable(props: Props) {
    const {node} = props

    return (
        <Table size={"small"} sx={SX.table}>
            <TableBody isLoading={false} cellCount={2}>
                <TableRow>
                    <TableCell sx={SX.title}>Host</TableCell>
                    <TableCell sx={SX.wrap}>{node.connection.host}</TableCell>
                </TableRow>
                <TableRow>
                    <TableCell sx={SX.title}>SSH</TableCell>
                    <TableCell sx={SX.wrap}>{node.connection.sshPort?.toString() ?? "-"}</TableCell>
                </TableRow>
                <TableRow>
                    <TableCell sx={SX.title}>Keeper</TableCell>
                    <TableCell sx={SX.wrap}>{node.connection.keeperPort?.toString() ?? "-"}</TableCell>
                </TableRow>
                <TableRow>
                    <TableCell sx={SX.title}>Database</TableCell>
                    <TableCell sx={SX.wrap}>{node.connection.dbPort?.toString() ?? "-"}</TableCell>
                </TableRow>
                <TableRow>
                    <TableCell sx={SX.title}>State</TableCell>
                    <TableCell sx={SX.wrap}>{node.keeper.state}</TableCell>
                </TableRow>
            </TableBody>
        </Table>
    )
}
