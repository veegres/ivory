import {KeyboardArrowDown} from "@mui/icons-material"
import {Accordion, AccordionDetails, AccordionSummary, Box} from "@mui/material"

import {useRouterPermission} from "../../../../api/permission/hook"
import {SxPropsMap} from "../../../../app/type"
import {ErrorSmart} from "../../../view/box/ErrorSmart"
import {MenuWrapperScroll} from "../menu/MenuWrapperScroll"
import {PermissionsItem} from "./PermissionsItem"

const SX: SxPropsMap = {
    box: {overflowY: "scroll", padding: "0px 2px 0px 0px"},
    summary: {
        backgroundColor: "rgba(0, 0, 0, .03)", minHeight: "auto", padding: "5px 10px", gap: 2, fontSize: "15px",
        "& .MuiAccordionSummary-content": {margin: "5px 0", gap: 1, alignItems: "center"}
    },
    details: {backgroundColor: "rgba(255, 255, 255, .05)", padding: "0px"},
}

export function Permissions() {
    const permissions = useRouterPermission()
    if (permissions.isLoading) return <Box>Loading</Box>
    if (permissions.isError) return <ErrorSmart error={permissions.error}/>
    return (
        <MenuWrapperScroll sx={SX.box}>
            {permissions.data?.map(({username, permissions}) => (
                <Accordion
                    key={username}
                    elevation={0}
                    variant={"outlined"}
                    sx={SX.accordion}
                    slotProps={{transition: {unmountOnExit: true}}}
                    disableGutters={true}
                >
                    <AccordionSummary sx={SX.summary} expandIcon={<KeyboardArrowDown sx={{fontSize: "20px"}}/>}>
                        {username}
                    </AccordionSummary>
                    <AccordionDetails sx={SX.details}>
                        {Object.entries(permissions).map(([name, status]) => (
                            <PermissionsItem key={name} username={username} name={name} status={status} view={"admin"}/>
                        ))}
                    </AccordionDetails>
                </Accordion>
            ))}
        </MenuWrapperScroll>
    )
}