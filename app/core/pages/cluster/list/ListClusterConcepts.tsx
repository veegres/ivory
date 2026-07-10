import {Box} from "@mui/material"
import {Fragment} from "react"

import {SxPropsMap} from "../../../../shared/helper/HelperType"

const SX: SxPropsMap = {
    box: {display: "flex", flexDirection: "column", alignItems: "center", gap: 0.5},
    title: {color: "text.disabled", fontSize: "0.8rem", textTransform: "uppercase", letterSpacing: "0.2em"},
    concepts: {
        display: "grid", gridTemplateColumns: "auto 1fr", alignItems: "baseline",
        columnGap: 2, rowGap: 1, maxWidth: 420, width: "100%", fontSize: "0.8rem",
    },
    term: {color: "text.primary", whiteSpace: "nowrap"},
    description: {color: "text.secondary", textAlign: "left"},
}

const CONCEPTS: {term: string, description: string}[] = [
    {term: "Node", description: "A server with three ports — keeper, database, and platform (e.g. Linux) — though keeper and database can share one"},
    {term: "Keeper", description: "Manages your database cluster. Provides features like failover, reload, reinit, etc. Functionality depends on your keeper"},
    {term: "Database", description: "The actual database software. Allows you to execute some template and custom queries"},
    {term: "Platform", description: "Deploys and manages your containers, and shows some server information to help you with troubleshooting"},
]

export function ListClusterConcepts() {
    return (
        <Box sx={SX.box}>
            <Box sx={SX.title}>Concepts</Box>
            <Box sx={SX.concepts}>
                {CONCEPTS.map(({term, description}) => (
                    <Fragment key={term}>
                        <Box sx={SX.term}>{term}</Box>
                        <Box sx={SX.description}>{description}</Box>
                    </Fragment>
                ))}
            </Box>
        </Box>
    )
}
