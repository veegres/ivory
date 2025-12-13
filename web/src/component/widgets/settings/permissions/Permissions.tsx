import {KeyboardArrowDown} from "@mui/icons-material"
import {Accordion, AccordionDetails, AccordionSummary, Box, Paper, Skeleton, Tab, Tabs} from "@mui/material"
import {useState} from "react"

import {useRouterInfo} from "../../../../api/management/hook"
import {useRouterPermissions} from "../../../../api/permission/hook"
import {SxPropsMap} from "../../../../app/type"
import {ErrorSmart} from "../../../view/box/ErrorSmart"
import {NoBox} from "../../../view/box/NoBox"
import {MenuWrapper} from "../menu/MenuWrapper"
import {PermissionsList} from "./PermissionsList"

const SX: SxPropsMap = {
    box: {overflowY: "scroll", padding: "0px 2px 0px 0px", display: "flex", flexDirection: "column", gap: 2},
    summary: {
        backgroundColor: "rgba(0, 0, 0, .03)", minHeight: "auto", padding: "5px 10px", gap: 2, fontSize: "15px",
        "& .MuiAccordionSummary-content": {margin: "5px 0", gap: 1, alignItems: "center"}
    },
    details: {backgroundColor: "rgba(255, 255, 255, .05)", padding: "0px"},
}

export function Permissions() {
    const [tab, setTab] = useState<"self" | "users">("self")
    const permissions = useRouterPermissions()
    const info = useRouterInfo()

    return (
        <MenuWrapper>
            <Tabs variant={"fullWidth"} value={tab} onChange={(_, value) => setTab(value as "self" | "users")}>
                <Tab label={"Self"} value={"self"} />
                <Tab label={"Users"} value={"users"} />
            </Tabs>
            <Box>
                {tab === "self" && renderSelfPermissions()}
                {tab === "users" && renderUserPermissions()}
            </Box>
        </MenuWrapper>
    )

    function renderSelfPermissions() {
        if (!info.data?.auth.user) return <ErrorSmart error={"Can't get user info"}/>
        const {username, permissions} = info.data.auth.user
        if (permissions.length === 0) return <NoBox text={"You don't have any permissions"}/>
        return (
            <Paper variant={"outlined"}>
                <PermissionsList permissions={permissions} username={username}/>
            </Paper>
        )
    }

    function renderUserPermissions() {
        if (permissions.isLoading) return renderLoading()
        if (permissions.isError) return <ErrorSmart error={permissions.error}/>
        if (!permissions.data || permissions.data.length === 0) return <NoBox text={"There is no users"}/>

        return permissions.data.map(({username, permissions}) => (
            <Accordion
                key={username}
                elevation={0}
                variant={"outlined"}
                slotProps={{transition: {unmountOnExit: true}}}
                disableGutters={true}
            >
                <AccordionSummary sx={SX.summary} expandIcon={<KeyboardArrowDown sx={{fontSize: "20px"}}/>}>
                    {username}
                </AccordionSummary>
                <AccordionDetails sx={SX.details}>
                    <PermissionsList permissions={permissions} username={username} view={"admin"}/>
                </AccordionDetails>
            </Accordion>
        ))
    }

    function renderLoading() {
        return (
            <>
                <Skeleton width={"100%"} height={"40px"}/>
                <Skeleton width={"100%"} height={"40px"}/>
                <Skeleton width={"100%"} height={"40px"}/>
            </>
        )
    }
}