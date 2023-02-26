import {Table, TableBody, TableCell, TableHead, TableRow} from "@mui/material";

type Props = {
    id: string,
}

export function QueryItemPlay(props: Props) {
    const {id} = props
    return (
        <Table size={"small"}>
            <TableHead>
                <TableRow>
                    <TableCell>head 1</TableCell>
                    <TableCell>head 2</TableCell>
                    <TableCell>head 3</TableCell>
                </TableRow>
            </TableHead>
            <TableBody>
                <TableRow>
                    <TableCell>{id}</TableCell>
                    <TableCell>cell 2</TableCell>
                    <TableCell>cell 3</TableCell>
                </TableRow>
                <TableRow>
                    <TableCell>{id}</TableCell>
                    <TableCell>cell 2</TableCell>
                    <TableCell>cell 3</TableCell>
                </TableRow>
                <TableRow>
                    <TableCell>{id}</TableCell>
                    <TableCell>cell 2</TableCell>
                    <TableCell>cell 3</TableCell>
                </TableRow>
            </TableBody>
        </Table>
    )
}
