import React, {useState} from "react";
import {Box, Tab, Tabs} from "@mui/material";
import {CertsNew} from "./CertsNew";
import {CertsList} from "./CertsList";
import {CertTabs, CertType} from "../../../app/types";

const SX = {
    box: { display: "flex", flexDirection: "column", gap: 1 },
    tabs: { minHeight: "25px" },
    tab: { minHeight: "25px", padding: "8px 12px", textTransform: "none" },
}

const TABS: CertTabs = {
    0: { label: "Client CA", type: CertType.CLIENT_CA },
    1: { label: "Client Cert", type: CertType.CLIENT_CERT },
    2: { label: "Client Key", type: CertType.CLIENT_KEY }
}

export type CertTypeProps = {
    type: CertType,
}

export function CertsContent() {
    const [tab, setTab] = useState(0)

    return (
        <Box sx={SX.box}>
            <Tabs sx={SX.tabs} value={tab} variant={"fullWidth"} centered onChange={(_, value) => setTab(value)}>
                {Object.entries(TABS).map(([key, value]) => (<Tab key={key} sx={SX.tab} label={value.label} />))}
            </Tabs>
            <CertsNew type={TABS[tab].type} />
            <CertsList type={TABS[tab].type} />
        </Box>
    )
}
