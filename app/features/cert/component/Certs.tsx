import {useState} from "react"

import {TabsButton} from "../../../shared/component/button/TabsButton"
import {LinearProgressStateful} from "../../../shared/component/progress/LinearProgressStateful"
import {LastElementScrolling} from "../../../shared/component/scrolling/LastElementScrolling"
import {Feature} from "../../Feature"
import {ManageAccess} from "../../management/component/ManageAccess"
import {useRouterCertList} from "../api/CertHook"
import {CertTabs, CertType} from "../api/CertType"
import {CertsList} from "./CertsList"
import {CertsNew} from "./CertsNew"

export const TABS: CertTabs = {
    0: {label: "Client CA", type: CertType.CLIENT_CA},
    1: {label: "Client Cert", type: CertType.CLIENT_CERT},
    2: {label: "Client Key", type: CertType.CLIENT_KEY},
}

export function Certs() {
    const [tab, setTab] = useState(0)
    const type = TABS[tab].type
    const query = useRouterCertList(type)
    const {data, error, isFetching} = query

    return (
        <LastElementScrolling>
            <TabsButton tab={tab} tabs={TABS} setTab={setTab}/>
            <ManageAccess feature={Feature.ManageCertCreate}><CertsNew type={type}/></ManageAccess>
            <LinearProgressStateful color={"inherit"} loading={isFetching} line/>
            <CertsList certs={data} error={error}/>
        </LastElementScrolling>
    )
}
