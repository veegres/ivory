import {useState} from "react"

import {useRouterCertList} from "../../../../api/cert/hook"
import {Permission} from "../../../../api/permission/type"
import {LinearProgressStateful} from "../../../view/progress/LinearProgressStateful"
import {Access} from "../../access/Access"
import {MenuWrapper} from "../menu/MenuWrapper"
import {CertsList} from "./CertsList"
import {CertsNew} from "./CertsNew"
import {CertsTab, TABS} from "./CertsTab"

export function Certs() {
    const [tab, setTab] = useState(0)
    const type = TABS[tab].type
    const query = useRouterCertList(type)
    const {data, error, isFetching} = query

    return (
        <MenuWrapper>
            <CertsTab tab={tab} setTab={setTab}/>
            <Access permission={Permission.ManageCertCreate}><CertsNew type={type}/></Access>
            <LinearProgressStateful color={"inherit"} loading={isFetching} line/>
            <CertsList certs={data} error={error}/>
        </MenuWrapper>
    )
}
