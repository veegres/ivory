import {useState} from "react"

import {useRouterCertList} from "../../../../api/cert/hook"
import {CertTabs, CertType} from "../../../../api/cert/type"
import {Feature} from "../../../../api/feature"
import {TabsButton} from "../../../view/button/TabsButton"
import {LinearProgressStateful} from "../../../view/progress/LinearProgressStateful"
import {Access} from "../../access/Access"
import {MenuWrapper} from "../menu/MenuWrapper"
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
        <MenuWrapper>
            <TabsButton tab={tab} tabs={TABS} setTab={setTab}/>
            <Access feature={Feature.ManageCertCreate}><CertsNew type={type}/></Access>
            <LinearProgressStateful color={"inherit"} loading={isFetching} line/>
            <CertsList certs={data} error={error}/>
        </MenuWrapper>
    )
}
