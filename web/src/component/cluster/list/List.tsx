import {Box, IconButton, Table, TableCell, TableHead, TableRow, Tooltip} from "@mui/material";
import {useQuery} from "@tanstack/react-query";
import {clusterApi} from "../../../app/api";
import {ListRow} from "./ListRow";
import {useState} from "react";
import {ErrorAlert} from "../../view/ErrorAlert";
import {TableBody} from "../../view/TableBody";
import {TableCellLoader} from "../../view/TableCellLoader";
import {Add} from "@mui/icons-material";
import {PageBlock} from "../../view/PageBlock";
import {ListRowNew} from "./ListRowNew";
import {InfoAlert} from "../../view/InfoAlert";
import {SxPropsMap} from "../../../app/types";
import {ListTags} from "./ListTags";

const SX: SxPropsMap = {
    table: {"tr:last-child td": {border: 0}},
    tags: {position: "relative", height: 0, top: "-35px"},
    nameCell: {width: "220px"},
    buttonCell: {width: "1%"},
    padding: {padding: "10px 20px"},
}

export function List() {
    const [editNode, setEditNode] = useState("")
    const [showNewElement, setShowNewElement] = useState(false)
    const query = useQuery(["cluster/list"], () => clusterApi.list())
    const {data: clusterMap, isLoading, isFetching, isError, error} = query
    const rows = Object.entries(clusterMap ?? {})

    return (
        <PageBlock withMarginTop={"35px"}>
            <Box sx={SX.tags}><ListTags/></Box>
            {renderContent()}
            {!isError && !showNewElement && !rows.length && (
                <Box sx={SX.padding}>
                    <InfoAlert text={"Please, add a cluster to get started"}/>
                </Box>
            )}
        </PageBlock>
    )

    function renderContent() {
        if (isError) return <ErrorAlert error={error}/>

        return (
            <Table size={"small"} sx={SX.table}>
                <TableHead>
                    <TableRow sx={{height: "20px"}}>
                        <TableCell sx={SX.nameCell}>Cluster Name</TableCell>
                        <TableCell>Instances</TableCell>
                        <TableCellLoader sx={SX.buttonCell} isFetching={isFetching && !isLoading}>
                            <Tooltip title={"Add new cluster"} disableInteractive>
                                <IconButton disabled={showNewElement} onClick={() => setShowNewElement(true)}>
                                    <Add/>
                                </IconButton>
                            </Tooltip>
                        </TableCellLoader>
                    </TableRow>
                </TableHead>
                <TableBody isLoading={isLoading} cellCount={3}>
                    {renderRows()}
                    <ListRowNew show={showNewElement} close={() => setShowNewElement(false)}/>
                </TableBody>
            </Table>
        )
    }

    function renderRows() {
        return rows.map(([name, cluster]) => {
            const editable = name === editNode
            const toggle = () => setEditNode(editable ? "" : name)

            return <ListRow key={name} cluster={cluster} name={name} editable={editable} toggle={toggle}/>
        })
    }
}
