import {Tab, Tabs} from "@mui/material";
import {CertTabs, CertType} from "../../../../api/cert/type";

export const TABS: CertTabs = {
    0: {label: "Client CA", type: CertType.CLIENT_CA},
    1: {label: "Client Cert", type: CertType.CLIENT_CERT},
    2: {label: "Client Key", type: CertType.CLIENT_KEY}
}

type Props = {
    tab: number,
    setTab: (index: number) => void,
}

export function CertsTab(props: Props) {
    const { tab, setTab } = props
    return (
        <Tabs value={tab} variant={"fullWidth"} centered onChange={(_, value) => setTab(value)}>
            {Object.entries(TABS).map(([key, value]) => (
                <Tab key={key} label={value.label}/>
            ))}
        </Tabs>
    )
}
