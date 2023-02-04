import React, {useState} from "react";
import {Box, Tab, Tabs} from "@mui/material";
import {CertsNew} from "./CertsNew";
import {CertsList} from "./CertsList";

const SX = {
    tabs: { minHeight: "25px", margin: "0 0 10px 0" },
    tab: { minHeight: "25px", padding: "8px 12px", textTransform: "none" },
}

export function CertsContent() {
    const [tab, setTab] = useState(0)

    return (
        <Box>
            <Tabs sx={SX.tabs} value={tab} variant={"fullWidth"} centered>
                <Tab sx={SX.tab} label={"Client RootCA"} onClick={() => setTab(0)}/>
                <Tab sx={SX.tab} label={"Client CA"} onClick={() => setTab(1)}/>
                <Tab sx={SX.tab} label={"Client Key"} onClick={() => setTab(2)}/>
            </Tabs>
            <CertsNew />
            <CertsList />
        </Box>
    )
}
