import {KeyboardArrowDown} from "@mui/icons-material"
import {Accordion, AccordionDetails, AccordionSummary, Skeleton} from "@mui/material"

import {useRouterPermissions} from "../../../../api/permission/hook"
import {SxPropsMap} from "../../../../app/type"
import {ErrorSmart} from "../../../view/box/ErrorSmart"
import {NoBox} from "../../../view/box/NoBox"
import {PermissionsList} from "./PermissionsList"

const SX: SxPropsMap = {
    box: {overflowY: "scroll", padding: "0px 2px 0px 0px", display: "flex", flexDirection: "column", gap: 2},
    summary: {
        backgroundColor: "rgba(0, 0, 0, .03)", minHeight: "auto", padding: "5px 10px", gap: 2, fontSize: "15px",
        "& .MuiAccordionSummary-content": {margin: "5px 0", gap: 1, alignItems: "center"}
    },
    details: {backgroundColor: "rgba(255, 255, 255, .05)", padding: "0px"},
}

export function PermissionsUserList() {
    const permissions = useRouterPermissions()

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